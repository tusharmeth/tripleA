package ledger.controller;

import ledger.dto.CreateTransactionRequest;
import ledger.exception.AccountNotFoundException;
import ledger.exception.InsufficientFundsException;
import ledger.exception.SameAccountException;
import ledger.exception.TransactionNotFoundException;
import ledger.model.Transaction;
import ledger.store.AccountStore;
import ledger.store.TransactionStore;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.Map;

@RestController
@Tag(name = "transactions")
@RequiredArgsConstructor
public class TransactionController {

    private final TransactionStore store;
    private final AccountStore accountStore;

    @Operation(summary = "Create a transaction", description = "Transfers an amount from source account to destination account")
    @PostMapping("/transactions")
    public ResponseEntity<?> createTransaction(@RequestBody CreateTransactionRequest req) {
        String sourceId = req.getSourceAccountId();
        String destId = req.getDestinationAccountId();

        if (sourceId == null || sourceId.isBlank()) {
            return ResponseEntity.badRequest().body(Map.of("error", "source_account_id is required"));
        }
        if (destId == null || destId.isBlank()) {
            return ResponseEntity.badRequest().body(Map.of("error", "destination_account_id is required"));
        }
        if (req.getAmount() == null || req.getAmount().isBlank()) {
            return ResponseEntity.badRequest().body(Map.of("error", "amount is required"));
        }

        double amount;
        try {
            amount = Double.parseDouble(req.getAmount());
        } catch (NumberFormatException e) {
            return ResponseEntity.badRequest().body(Map.of("error", "amount must be a positive number"));
        }
        if (amount <= 0) {
            return ResponseEntity.badRequest().body(Map.of("error", "amount must be a positive number"));
        }

        try {
            accountStore.transfer(sourceId, destId, amount);
        } catch (AccountNotFoundException e) {
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of("error", "one or both accounts not found"));
        } catch (InsufficientFundsException e) {
            return ResponseEntity.status(HttpStatus.UNPROCESSABLE_ENTITY).body(Map.of("error", "insufficient funds"));
        } catch (SameAccountException e) {
            return ResponseEntity.badRequest().body(Map.of("error", "source and destination accounts must differ"));
        }

        Transaction tx = store.create(sourceId, destId, req.getAmount());
        return ResponseEntity.status(HttpStatus.CREATED).body(tx);
    }

    @Operation(summary = "Get a transaction", description = "Retrieves a transaction by its ID")
    @GetMapping("/transactions/{id}")
    public ResponseEntity<?> getTransaction(@PathVariable String id) {
        try {
            return ResponseEntity.ok(store.getById(id));
        } catch (TransactionNotFoundException e) {
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of("error", "transaction not found"));
        }
    }
}
