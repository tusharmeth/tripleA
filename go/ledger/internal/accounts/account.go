package accounts

import "time"

type Account struct {
	ID             string    `json:"account_id"`
	InitialBalance float64   `json:"initial_balance"`
	Balance        float64   `json:"balance"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
