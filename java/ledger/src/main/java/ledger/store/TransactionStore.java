package ledger.store;

import ledger.exception.TransactionNotFoundException;
import ledger.model.Transaction;
import ledger.repository.TransactionRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Component;
import org.springframework.transaction.annotation.Transactional;

import java.security.SecureRandom;
import java.time.Instant;
import java.util.HexFormat;

@Component
@RequiredArgsConstructor
public class TransactionStore {

    private static final SecureRandom RANDOM = new SecureRandom();

    private final TransactionRepository transactionRepository;

    @Transactional
    public Transaction create(String sourceId, String destId, String amount) {
        String id = generateId();
        Transaction tx = new Transaction(id, sourceId, destId, amount, Instant.now());
        return transactionRepository.save(tx);
    }

    @Transactional(readOnly = true)
    public Transaction getById(String id) {
        return transactionRepository.findById(id)
                .orElseThrow(TransactionNotFoundException::new);
    }

    private static String generateId() {
        byte[] bytes = new byte[8];
        RANDOM.nextBytes(bytes);
        return HexFormat.of().formatHex(bytes);
    }
}
