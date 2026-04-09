package main

import (
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/tusharmethwani/ledger/docs"
	"github.com/tusharmethwani/ledger/internal/accounts"
	"github.com/tusharmethwani/ledger/internal/transactions"
)

// @title           Ledger API
// @version         1.0
// @description     API for managing accounts and transactions.

// @host      localhost:8080
// @BasePath  /

func main() {
	accountStore := accounts.NewStore()
	txStore := transactions.NewStore()

	mux := http.NewServeMux()
	accounts.NewHandler(accountStore).RegisterRoutes(mux)
	transactions.NewHandler(txStore, accountStore).RegisterRoutes(mux)

	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	log.Println("Ledger service listening on :8080")
	log.Println("Swagger UI available at http://localhost:8080/swagger/index.html")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
