package accounts

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrAccountNotFound   = errors.New("account not found")
	ErrAccountExists     = errors.New("account already exists")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrSameAccount       = errors.New("source and destination accounts must differ")
)

// AccountStorer is the interface both in-memory and postgres stores satisfy.
type AccountStorer interface {
	Create(id string, initialBalance float64) (*Account, error)
	GetByID(id string) (*Account, error)
	Transfer(sourceID, destID string, amount float64) error
}

// ---- In-memory implementation ----

type MemoryStore struct {
	mu       sync.RWMutex
	accounts map[string]*Account
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{accounts: make(map[string]*Account)}
}

func (s *MemoryStore) Create(id string, initialBalance float64) (*Account, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.accounts[id]; exists {
		return nil, ErrAccountExists
	}

	now := time.Now().UTC()
	acc := &Account{
		ID:             id,
		InitialBalance: initialBalance,
		Balance:        initialBalance,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	s.accounts[id] = acc
	return acc, nil
}

func (s *MemoryStore) GetByID(id string) (*Account, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	acc, exists := s.accounts[id]
	if !exists {
		return nil, ErrAccountNotFound
	}
	return acc, nil
}

// Transfer atomically debits source and credits destination.
func (s *MemoryStore) Transfer(sourceID, destID string, amount float64) error {
	if sourceID == destID {
		return ErrSameAccount
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	src, ok := s.accounts[sourceID]
	if !ok {
		return ErrAccountNotFound
	}
	dst, ok := s.accounts[destID]
	if !ok {
		return ErrAccountNotFound
	}
	if src.Balance < amount {
		return ErrInsufficientFunds
	}

	now := time.Now().UTC()
	src.Balance -= amount
	src.UpdatedAt = now
	dst.Balance += amount
	dst.UpdatedAt = now
	return nil
}

// ---- Postgres implementation ----

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func (s *PostgresStore) Create(id string, initialBalance float64) (*Account, error) {
	now := time.Now().UTC()
	_, err := s.db.Exec(
		`INSERT INTO accounts (id, initial_balance, balance, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		id, initialBalance, initialBalance, now, now,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrAccountExists
		}
		return nil, fmt.Errorf("inserting account: %w", err)
	}
	return &Account{
		ID:             id,
		InitialBalance: initialBalance,
		Balance:        initialBalance,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

func (s *PostgresStore) GetByID(id string) (*Account, error) {
	row := s.db.QueryRow(
		`SELECT id, initial_balance, balance, created_at, updated_at
		 FROM accounts WHERE id = $1`, id,
	)
	var acc Account
	if err := row.Scan(&acc.ID, &acc.InitialBalance, &acc.Balance, &acc.CreatedAt, &acc.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAccountNotFound
		}
		return nil, fmt.Errorf("querying account: %w", err)
	}
	return &acc, nil
}

// Transfer debits source and credits destination within a single DB transaction.
func (s *PostgresStore) Transfer(sourceID, destID string, amount float64) error {
	if sourceID == destID {
		return ErrSameAccount
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	now := time.Now().UTC()

	// Lock rows in consistent order to avoid deadlocks.
	first, second := sourceID, destID
	if sourceID > destID {
		first, second = destID, sourceID
	}
	var dummy string
	if err := tx.QueryRow(`SELECT id FROM accounts WHERE id = $1 FOR UPDATE`, first).Scan(&dummy); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrAccountNotFound
		}
		return fmt.Errorf("locking account: %w", err)
	}
	if err := tx.QueryRow(`SELECT id FROM accounts WHERE id = $1 FOR UPDATE`, second).Scan(&dummy); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrAccountNotFound
		}
		return fmt.Errorf("locking account: %w", err)
	}

	var srcBalance float64
	if err := tx.QueryRow(`SELECT balance FROM accounts WHERE id = $1`, sourceID).Scan(&srcBalance); err != nil {
		return fmt.Errorf("reading source balance: %w", err)
	}
	if srcBalance < amount {
		return ErrInsufficientFunds
	}

	if _, err := tx.Exec(
		`UPDATE accounts SET balance = balance - $1, updated_at = $2 WHERE id = $3`,
		amount, now, sourceID,
	); err != nil {
		return fmt.Errorf("debiting source: %w", err)
	}
	if _, err := tx.Exec(
		`UPDATE accounts SET balance = balance + $1, updated_at = $2 WHERE id = $3`,
		amount, now, destID,
	); err != nil {
		return fmt.Errorf("crediting destination: %w", err)
	}

	return tx.Commit()
}

func isUniqueViolation(err error) bool {
	// lib/pq error code 23505 = unique_violation
	type pgErr interface{ Get(byte) string }
	if pe, ok := err.(pgErr); ok {
		return pe.Get('C') == "23505"
	}
	return false
}
