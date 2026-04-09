package ledger.exception;

public class SameAccountException extends RuntimeException {
    public SameAccountException() { super("source and destination accounts must differ"); }
}
