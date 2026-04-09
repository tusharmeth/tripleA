package accounts

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Handler struct {
	store AccountStorer
}

func NewHandler(store AccountStorer) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /accounts", h.createAccount)
	mux.HandleFunc("GET /accounts/{id}", h.getAccount)
}

type createAccountRequest struct {
	AccountID      string  `json:"account_id"`
	InitialBalance float64 `json:"initial_balance"`
}

// createAccount godoc
// @Summary      Create an account
// @Description  Creates a new account with a given ID and initial balance
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        body  body      createAccountRequest  true  "Account to create"
// @Success      201   {object}  Account
// @Failure      400   {object}  map[string]string
// @Failure      409   {object}  map[string]string
// @Router       /accounts [post]
func (h *Handler) createAccount(w http.ResponseWriter, r *http.Request) {
	var req createAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.AccountID == "" {
		writeError(w, http.StatusBadRequest, "account_id is required")
		return
	}
	if req.InitialBalance < 0 {
		writeError(w, http.StatusBadRequest, "initial_balance must be non-negative")
		return
	}

	acc, err := h.store.Create(req.AccountID, req.InitialBalance)
	if err != nil {
		if errors.Is(err, ErrAccountExists) {
			writeError(w, http.StatusConflict, "account already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to create account")
		return
	}

	writeJSON(w, http.StatusCreated, acc)
}

// getAccount godoc
// @Summary      Get an account
// @Description  Retrieves an account by its ID
// @Tags         accounts
// @Produce      json
// @Param        id   path      string  true  "Account ID"
// @Success      200  {object}  Account
// @Failure      404  {object}  map[string]string
// @Router       /accounts/{id} [get]
func (h *Handler) getAccount(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	acc, err := h.store.GetByID(id)
	if err != nil {
		if errors.Is(err, ErrAccountNotFound) {
			writeError(w, http.StatusNotFound, "account not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get account")
		return
	}

	writeJSON(w, http.StatusOK, acc)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
