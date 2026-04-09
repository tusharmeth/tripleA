package ledger.model;

import com.fasterxml.jackson.annotation.JsonProperty;
import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;

import java.time.Instant;

@Entity
@Table(name = "transactions")
@Getter
@NoArgsConstructor
@AllArgsConstructor
public class Transaction {

    @Id
    @Column(name = "id")
    @JsonProperty("transaction_id")
    private String id;

    @Column(name = "source_account_id", nullable = false)
    @JsonProperty("source_account_id")
    private String sourceAccountId;

    @Column(name = "destination_account_id", nullable = false)
    @JsonProperty("destination_account_id")
    private String destinationAccountId;

    @Column(name = "amount", nullable = false)
    @JsonProperty("amount")
    private String amount;

    @Column(name = "created_at", nullable = false, updatable = false)
    @JsonProperty("created_at")
    private Instant createdAt;
}
