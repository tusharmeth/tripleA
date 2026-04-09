# Ledger

## How to Run

**1. Clone the repository**
```bash
git clone https://github.com/tusharmeth/tripleA.git
cd java
cd ledger
```

**2. Start the application**
```bash
docker compose up --build
```

This starts both the PostgreSQL database and the Spring Boot app. The app will be ready when you see the Spring Boot banner in the logs.

**3. Open Swagger UI**

Navigate to: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

---

## Using the API

**Create accounts**

Use `POST /accounts` with a body like:
```json
{
  "account_id": "alice",
  "initial_balance": 1000.00
}
```

Create a second account:
```json
{
  "account_id": "bob",
  "initial_balance": 500.00
}
```

**Perform a transaction**

Use `POST /transactions` to transfer funds between accounts:
```json
{
  "source_account_id": "alice",
  "destination_account_id": "bob",
  "amount": "200.00"
}
```

**View accounts and transactions**

- `GET /accounts` — list all accounts
- `GET /accounts/{id}` — get a specific account
- `GET /transactions` — list all transactions
- `GET /transactions/{id}` — get a specific transaction

---

## Stop the application

```bash
docker compose down
```

To also remove the database volume:
```bash
docker compose down -v
```
