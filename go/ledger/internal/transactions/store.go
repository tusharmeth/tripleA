package transactions

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/tusharmethwani/ledger/internal/accounts"
)

var ErrTransactionNotFound = errors.New("transaction not found")

// TransactionStorer is the interface both in-memory and postgres stores satisfy.
type TransactionStorer interface {
	// Create executes the transfer and records the transaction atomically.
	Create(sourceID, destID, amount string) (*Transaction, error)
	GetByID(id string) (*Transaction, error)
}

// ---- In-memory implementation ----

type MemoryStore struct {
	mu           sync.RWMutex
	transactions map[string]*Transaction
	accountStore accounts.AccountStorer
}

func NewMemoryStore(accountStore accounts.AccountStorer) *MemoryStore {
	return &MemoryStore{
		transactions: make(map[string]*Transaction),
		accountStore: accountStore,
	}
}

func (s *MemoryStore) Create(sourceID, destID, amount string) (*Transaction, error) {
	amountF, err := strconv.ParseFloat(amount, 64)
	if err != nil || amountF <= 0 {
		return nil, fmt.Errorf("invalid amount: %w", err)
	}

	if err := s.accountStore.Transfer(sourceID, destID, amountF); err != nil {
		return nil, err
	}

	id, err := generateID()
	if err != nil {
		return nil, err
	}

	tx := &Transaction{
		ID:                   id,
		SourceAccountID:      sourceID,
		DestinationAccountID: destID,
		Amount:               amount,
		CreatedAt:            time.Now().UTC(),
	}

	s.mu.Lock()
	s.transactions[id] = tx
	s.mu.Unlock()

	return tx, nil
}

func (s *MemoryStore) GetByID(id string) (*Transaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tx, exists := s.transactions[id]
	if !exists {
		return nil, ErrTransactionNotFound
	}
	return tx, nil
}

// ---- Postgres implementation ----

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

// Create executes the account transfer and inserts the transaction record within a
// single database transaction, guaranteeing atomicity.
func (s *PostgresStore) Create(sourceID, destID, amount string) (*Transaction, error) {
	if sourceID == destID {
		return nil, accounts.ErrSameAccount
	}

	amountF, err := strconv.ParseFloat(amount, 64)
	if err != nil || amountF <= 0 {
		return nil, fmt.Errorf("invalid amount")
	}

	dbTx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("beginning transaction: %w", err)
	}
	defer dbTx.Rollback() //nolint:errcheck

	now := time.Now().UTC()

	// Lock rows in consistent order to prevent deadlocks.
	first, second := sourceID, destID
	if sourceID > destID {
		first, second = destID, sourceID
	}
	var dummy string
	if err := dbTx.QueryRow(`SELECT id FROM accounts WHERE id = $1 FOR UPDATE`, first).Scan(&dummy); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, accounts.ErrAccountNotFound
		}
		return nil, fmt.Errorf("locking account: %w", err)
	}
	if err := dbTx.QueryRow(`SELECT id FROM accounts WHERE id = $1 FOR UPDATE`, second).Scan(&dummy); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, accounts.ErrAccountNotFound
		}
		return nil, fmt.Errorf("locking account: %w", err)
	}

	var srcBalance float64
	if err := dbTx.QueryRow(`SELECT balance FROM accounts WHERE id = $1`, sourceID).Scan(&srcBalance); err != nil {
		return nil, fmt.Errorf("reading source balance: %w", err)
	}
	if srcBalance < amountF {
		return nil, accounts.ErrInsufficientFunds
	}

	if _, err := dbTx.Exec(
		`UPDATE accounts SET balance = balance - $1, updated_at = $2 WHERE id = $3`,
		amountF, now, sourceID,
	); err != nil {
		return nil, fmt.Errorf("debiting source: %w", err)
	}
	if _, err := dbTx.Exec(
		`UPDATE accounts SET balance = balance + $1, updated_at = $2 WHERE id = $3`,
		amountF, now, destID,
	); err != nil {
		return nil, fmt.Errorf("crediting destination: %w", err)
	}

	id, err := generateID()
	if err != nil {
		return nil, err
	}

	if _, err := dbTx.Exec(
		`INSERT INTO transactions (id, source_account_id, destination_account_id, amount, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		id, sourceID, destID, amount, now,
	); err != nil {
		return nil, fmt.Errorf("inserting transaction: %w", err)
	}

	if err := dbTx.Commit(); err != nil {
		return nil, fmt.Errorf("committing transaction: %w", err)
	}

	return &Transaction{
		ID:                   id,
		SourceAccountID:      sourceID,
		DestinationAccountID: destID,
		Amount:               amount,
		CreatedAt:            now,
	}, nil
}

func (s *PostgresStore) GetByID(id string) (*Transaction, error) {
	row := s.db.QueryRow(
		`SELECT id, source_account_id, destination_account_id, amount, created_at
		 FROM transactions WHERE id = $1`, id,
	)
	var tx Transaction
	if err := row.Scan(&tx.ID, &tx.SourceAccountID, &tx.DestinationAccountID, &tx.Amount, &tx.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTransactionNotFound
		}
		return nil, fmt.Errorf("querying transaction: %w", err)
	}
	return &tx, nil
}

func generateID() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
