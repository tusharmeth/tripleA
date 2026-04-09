package main

import (
	"log"
	"net/http"
	"os"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/tusharmethwani/ledger/docs"
	"github.com/tusharmethwani/ledger/internal/accounts"
	"github.com/tusharmethwani/ledger/internal/db"
	"github.com/tusharmethwani/ledger/internal/transactions"
)

// @title           Ledger API
// @version         1.0
// @description     API for managing accounts and transactions.

// @host      localhost:9090
// @BasePath  /

func main() {
	var (
		accountStore  accounts.AccountStorer
		txStore       transactions.TransactionStorer
	)

	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		log.Println("Connecting to postgres...")
		database, err := db.Connect(dsn)
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}
		if _, err := database.Exec(db.Schema); err != nil {
			log.Fatalf("failed to apply schema: %v", err)
		}
		log.Println("Postgres connected and schema applied.")
		accountStore = accounts.NewPostgresStore(database)
		txStore = transactions.NewPostgresStore(database)
	} else {
		log.Println("DATABASE_URL not set — using in-memory stores.")
		memAccountStore := accounts.NewMemoryStore()
		accountStore = memAccountStore
		txStore = transactions.NewMemoryStore(memAccountStore)
	}

	mux := http.NewServeMux()
	accounts.NewHandler(accountStore).RegisterRoutes(mux)
	transactions.NewHandler(txStore).RegisterRoutes(mux)

	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	log.Println("Ledger service listening on :9090")
	log.Println("Swagger UI available at http://localhost:9090/swagger/index.html")
	if err := http.ListenAndServe(":9090", mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
