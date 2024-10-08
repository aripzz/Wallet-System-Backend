-- migrate:up
CREATE TABLE
    Wallets (
        id SERIAL PRIMARY KEY,
        user_id INT,
        wallet_types_id INT NOT NULL,
        active BOOLEAN NOT NULL DEFAULT TRUE,
        balance DECIMAL(10, 2),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES Users (id),
        FOREIGN KEY (wallet_types_id) REFERENCES wallet_types (id)
    );

INSERT INTO
    Wallets (user_id, wallet_types_id, balance, active)
VALUES
    (1, 1, 10000, true),
    (1, 2, 10000, true),
    (1, 3, 10000, true);

-- migrate:down
TRUNCATE TABLE Wallets;

DROP TABLE Wallets;