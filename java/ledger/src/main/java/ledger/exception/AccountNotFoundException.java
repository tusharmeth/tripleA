package ledger.exception;

public class AccountNotFoundException extends RuntimeException {
    public AccountNotFoundException() { super("account not found"); }
}
