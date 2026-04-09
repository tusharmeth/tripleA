package ledger.store;

import ledger.exception.TransactionNotFoundException;
import ledger.model.Transaction;
import org.springframework.stereotype.Component;

import java.security.SecureRandom;
import java.time.Instant;
import java.util.HexFormat;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.locks.ReentrantReadWriteLock;

@Component
public class TransactionStore {

    private static final SecureRandom RANDOM = new SecureRandom();

    private final ConcurrentHashMap<String, Transaction> transactions = new ConcurrentHashMap<>();
    private final ReentrantReadWriteLock lock = new ReentrantReadWriteLock();

    public Transaction create(String sourceId, String destId, String amount) {
        String id = generateId();
        Transaction tx = new Transaction(id, sourceId, destId, amount, Instant.now());

        lock.writeLock().lock();
        try {
            transactions.put(id, tx);
        } finally {
            lock.writeLock().unlock();
        }
        return tx;
    }

    public Transaction getById(String id) {
        lock.readLock().lock();
        try {
            Transaction tx = transactions.get(id);
            if (tx == null) throw new TransactionNotFoundException();
            return tx;
        } finally {
            lock.readLock().unlock();
        }
    }

    private static String generateId() {
        byte[] bytes = new byte[8];
        RANDOM.nextBytes(bytes);
        return HexFormat.of().formatHex(bytes);
    }
}
