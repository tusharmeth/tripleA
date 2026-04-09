package transactions

import "time"

type Transaction struct {
	ID                  string    `json:"transaction_id"`
	SourceAccountID     string    `json:"source_account_id"`
	DestinationAccountID string   `json:"destination_account_id"`
	Amount              string    `json:"amount"`
	CreatedAt           time.Time `json:"created_at"`
}
