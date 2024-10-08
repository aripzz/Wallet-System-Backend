-- migrate:up
CREATE TYPE status_enum AS ENUM ('success', 'pending', 'failed');

CREATE TABLE
    Transactions (
        id SERIAL PRIMARY KEY,
        users_id INT NOT NULL,
        products_id INT NOT NULL,
        amount DECIMAL(10, 2) NOT NULL,
        status status_enum NOT NULL DEFAULT 'pending',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (users_id) REFERENCES users (id),
        FOREIGN KEY (products_id) REFERENCES Products (id)
    );

-- migrate:down
DROP TABLE Transactions;

DROP TYPE status_enum;