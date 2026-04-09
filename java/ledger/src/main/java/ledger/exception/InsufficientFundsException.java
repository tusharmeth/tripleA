package ledger.exception;

public class InsufficientFundsException extends RuntimeException {
    public InsufficientFundsException() { super("insufficient funds"); }
}
