package accounts

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrAccountNotFound      = errors.New("account not found")
	ErrAccountExists        = errors.New("account already exists")
	ErrInsufficientFunds    = errors.New("insufficient funds")
	ErrSameAccount          = errors.New("source and destination accounts must differ")
)

type Store struct {
	mu       sync.RWMutex
	accounts map[string]*Account
}

func NewStore() *Store {
	return &Store{accounts: make(map[string]*Account)}
}

func (s *Store) Create(id string, initialBalance float64) (*Account, error) {
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

func (s *Store) GetByID(id string) (*Account, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	acc, exists := s.accounts[id]
	if !exists {
		return nil, ErrAccountNotFound
	}
	return acc, nil
}

// Transfer atomically debits source and credits destination.
// Lock order is alphabetical by ID to prevent deadlocks.
func (s *Store) Transfer(sourceID, destID string, amount float64) error {
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
