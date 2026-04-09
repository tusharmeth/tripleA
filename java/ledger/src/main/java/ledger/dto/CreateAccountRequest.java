package ledger.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Getter;

@Getter
public class CreateAccountRequest {

    @JsonProperty("account_id")
    private String accountId;

    @JsonProperty("initial_balance")
    private Double initialBalance;
}
