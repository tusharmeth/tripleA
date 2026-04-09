package transactions

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

var ErrTransactionNotFound = errors.New("transaction not found")

type Store struct {
	mu           sync.RWMutex
	transactions map[string]*Transaction
}

func NewStore() *Store {
	return &Store{transactions: make(map[string]*Transaction)}
}

func (s *Store) Create(sourceID, destID, amount string) (*Transaction, error) {
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

func (s *Store) GetByID(id string) (*Transaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tx, exists := s.transactions[id]
	if !exists {
		return nil, ErrTransactionNotFound
	}
	return tx, nil
}

func generateID() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
