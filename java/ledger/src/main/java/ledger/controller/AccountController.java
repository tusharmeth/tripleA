package ledger.controller;

import ledger.dto.CreateAccountRequest;
import ledger.exception.AccountExistsException;
import ledger.model.Account;
import ledger.store.AccountStore;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import ledger.exception.AccountNotFoundException;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.Map;

@RestController
@Tag(name = "accounts")
@RequiredArgsConstructor
public class AccountController {

    private final AccountStore store;

    @Operation(summary = "Create an account", description = "Creates a new account with a given ID and initial balance")
    @PostMapping("/accounts")
    public ResponseEntity<?> createAccount(@RequestBody CreateAccountRequest req) {
        if (req.getAccountId() == null || req.getAccountId().isBlank()) {
            return ResponseEntity.badRequest().body(Map.of("error", "account_id is required"));
        }
        if (req.getInitialBalance() == null || req.getInitialBalance() < 0) {
            return ResponseEntity.badRequest().body(Map.of("error", "initial_balance must be non-negative"));
        }

        try {
            Account acc = store.create(req.getAccountId(), req.getInitialBalance());
            return ResponseEntity.status(HttpStatus.CREATED).body(acc);
        } catch (AccountExistsException e) {
            return ResponseEntity.status(HttpStatus.CONFLICT).body(Map.of("error", "account already exists"));
        }
    }

    @Operation(summary = "Get an account", description = "Retrieves an account by its ID")
    @GetMapping("/accounts/{id}")
    public ResponseEntity<?> getAccount(@PathVariable String id) {
        try {
            return ResponseEntity.ok(store.getById(id));
        } catch (AccountNotFoundException e) {
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of("error", "account not found"));
        }
    }
}
