package ledger.exception;

public class AccountExistsException extends RuntimeException {
    public AccountExistsException() { super("account already exists"); }
}
