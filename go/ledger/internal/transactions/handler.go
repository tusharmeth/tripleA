package transactions

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/tusharmethwani/ledger/internal/accounts"
)

type Handler struct {
	store        *Store
	accountStore *accounts.Store
}

func NewHandler(store *Store, accountStore *accounts.Store) *Handler {
	return &Handler{store: store, accountStore: accountStore}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /transactions", h.createTransaction)
	mux.HandleFunc("GET /transactions/{id}", h.getTransaction)
}

type createTransactionRequest struct {
	SourceAccountID      any    `json:"source_account_id"`
	DestinationAccountID any    `json:"destination_account_id"`
	Amount               string `json:"amount"`
}

// createTransaction godoc
// @Summary      Create a transaction
// @Description  Transfers an amount from source account to destination account
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        body  body      createTransactionRequest  true  "Transaction to create"
// @Success      201   {object}  Transaction
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      422   {object}  map[string]string
// @Router       /transactions [post]
func (h *Handler) createTransaction(w http.ResponseWriter, r *http.Request) {
	var req createTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	sourceID := anyToString(req.SourceAccountID)
	destID := anyToString(req.DestinationAccountID)

	if sourceID == "" {
		writeError(w, http.StatusBadRequest, "source_account_id is required")
		return
	}
	if destID == "" {
		writeError(w, http.StatusBadRequest, "destination_account_id is required")
		return
	}
	if req.Amount == "" {
		writeError(w, http.StatusBadRequest, "amount is required")
		return
	}

	amount, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil || amount <= 0 {
		writeError(w, http.StatusBadRequest, "amount must be a positive number")
		return
	}

	if err := h.accountStore.Transfer(sourceID, destID, amount); err != nil {
		switch {
		case errors.Is(err, accounts.ErrAccountNotFound):
			writeError(w, http.StatusNotFound, "one or both accounts not found")
		case errors.Is(err, accounts.ErrInsufficientFunds):
			writeError(w, http.StatusUnprocessableEntity, "insufficient funds")
		case errors.Is(err, accounts.ErrSameAccount):
			writeError(w, http.StatusBadRequest, "source and destination accounts must differ")
		default:
			writeError(w, http.StatusInternalServerError, "failed to process transaction")
		}
		return
	}

	tx, err := h.store.Create(sourceID, destID, req.Amount)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "transaction processed but failed to record")
		return
	}

	writeJSON(w, http.StatusCreated, tx)
}

// getTransaction godoc
// @Summary      Get a transaction
// @Description  Retrieves a transaction by its ID
// @Tags         transactions
// @Produce      json
// @Param        id   path      string  true  "Transaction ID"
// @Success      200  {object}  Transaction
// @Failure      404  {object}  map[string]string
// @Router       /transactions/{id} [get]
func (h *Handler) getTransaction(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	tx, err := h.store.GetByID(id)
	if err != nil {
		if errors.Is(err, ErrTransactionNotFound) {
			writeError(w, http.StatusNotFound, "transaction not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get transaction")
		return
	}

	writeJSON(w, http.StatusOK, tx)
}

// anyToString handles account IDs sent as either a JSON number or a string.
func anyToString(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	default:
		return ""
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
