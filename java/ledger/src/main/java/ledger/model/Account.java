package ledger.model;

import com.fasterxml.jackson.annotation.JsonProperty;
import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.time.Instant;

@Entity
@Table(name = "accounts")
@Getter
@NoArgsConstructor
public class Account {

    @Id
    @Column(name = "id")
    @JsonProperty("account_id")
    private String id;

    @Column(name = "initial_balance", nullable = false)
    @JsonProperty("initial_balance")
    private double initialBalance;

    @Setter
    @Column(name = "balance", nullable = false)
    @JsonProperty("balance")
    private double balance;

    @Column(name = "created_at", nullable = false, updatable = false)
    @JsonProperty("created_at")
    private Instant createdAt;

    @Setter
    @Column(name = "updated_at", nullable = false)
    @JsonProperty("updated_at")
    private Instant updatedAt;

    public Account(String id, double initialBalance, Instant createdAt) {
        this.id = id;
        this.initialBalance = initialBalance;
        this.balance = initialBalance;
        this.createdAt = createdAt;
        this.updatedAt = createdAt;
    }
}
