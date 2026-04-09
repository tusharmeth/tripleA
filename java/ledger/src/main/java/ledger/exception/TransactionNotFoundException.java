package ledger.exception;

public class TransactionNotFoundException extends RuntimeException {
    public TransactionNotFoundException() { super("transaction not found"); }
}
