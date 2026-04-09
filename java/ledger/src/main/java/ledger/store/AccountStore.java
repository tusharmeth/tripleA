package ledger.store;

import ledger.exception.AccountExistsException;
import ledger.exception.AccountNotFoundException;
import ledger.exception.InsufficientFundsException;
import ledger.exception.SameAccountException;
import ledger.model.Account;
import org.springframework.stereotype.Component;

import java.time.Instant;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.locks.ReentrantReadWriteLock;

@Component
public class AccountStore {

    private final ConcurrentHashMap<String, Account> accounts = new ConcurrentHashMap<>();
    private final ReentrantReadWriteLock lock = new ReentrantReadWriteLock();

    public Account create(String id, double initialBalance) {
        lock.writeLock().lock();
        try {
            if (accounts.containsKey(id)) {
                throw new AccountExistsException();
            }
            Instant now = Instant.now();
            Account acc = new Account(id, initialBalance, now);
            accounts.put(id, acc);
            return acc;
        } finally {
            lock.writeLock().unlock();
        }
    }

    public Account getById(String id) {
        lock.readLock().lock();
        try {
            Account acc = accounts.get(id);
            if (acc == null) throw new AccountNotFoundException();
            return acc;
        } finally {
            lock.readLock().unlock();
        }
    }

    // Atomically debits source and credits destination.
    public void transfer(String sourceId, String destId, double amount) {
        if (sourceId.equals(destId)) throw new SameAccountException();

        lock.writeLock().lock();
        try {
            Account src = accounts.get(sourceId);
            if (src == null) throw new AccountNotFoundException();

            Account dst = accounts.get(destId);
            if (dst == null) throw new AccountNotFoundException();

            if (src.getBalance() < amount) throw new InsufficientFundsException();

            Instant now = Instant.now();
            src.setBalance(src.getBalance() - amount);
            src.setUpdatedAt(now);
            dst.setBalance(dst.getBalance() + amount);
            dst.setUpdatedAt(now);
        } finally {
            lock.writeLock().unlock();
        }
    }
}
