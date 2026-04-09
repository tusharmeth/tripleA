package ledger.store;

import ledger.exception.AccountExistsException;
import ledger.exception.AccountNotFoundException;
import ledger.exception.InsufficientFundsException;
import ledger.exception.SameAccountException;
import ledger.model.Account;
import ledger.repository.AccountRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Component;
import org.springframework.transaction.annotation.Transactional;

import java.time.Instant;

@Component
@RequiredArgsConstructor
public class AccountStore {

    private final AccountRepository accountRepository;

    @Transactional
    public Account create(String id, double initialBalance) {
        if (accountRepository.existsById(id)) {
            throw new AccountExistsException();
        }
        Instant now = Instant.now();
        Account acc = new Account(id, initialBalance, now);
        return accountRepository.save(acc);
    }

    @Transactional(readOnly = true)
    public Account getById(String id) {
        return accountRepository.findById(id)
                .orElseThrow(AccountNotFoundException::new);
    }

    // Atomically debits source and credits destination using DB-level pessimistic locks.
    // Accounts are locked in lexicographic ID order to prevent deadlocks.
    @Transactional
    public void transfer(String sourceId, String destId, double amount) {
        if (sourceId.equals(destId)) throw new SameAccountException();

        String firstId  = sourceId.compareTo(destId) < 0 ? sourceId : destId;
        String secondId = sourceId.compareTo(destId) < 0 ? destId   : sourceId;

        Account first  = accountRepository.findByIdForUpdate(firstId)
                .orElseThrow(AccountNotFoundException::new);
        Account second = accountRepository.findByIdForUpdate(secondId)
                .orElseThrow(AccountNotFoundException::new);

        Account src = firstId.equals(sourceId) ? first : second;
        Account dst = firstId.equals(destId)   ? first : second;

        if (src.getBalance() < amount) throw new InsufficientFundsException();

        Instant now = Instant.now();
        src.setBalance(src.getBalance() - amount);
        src.setUpdatedAt(now);
        dst.setBalance(dst.getBalance() + amount);
        dst.setUpdatedAt(now);

        accountRepository.save(src);
        accountRepository.save(dst);
    }
}
