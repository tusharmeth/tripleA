package ledger.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Getter;

@Getter
public class CreateTransactionRequest {

    @JsonProperty("source_account_id")
    private String sourceAccountId;

    @JsonProperty("destination_account_id")
    private String destinationAccountId;

    @JsonProperty("amount")
    private String amount;
}
