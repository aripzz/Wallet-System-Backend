-- migrate:up
CREATE TABLE
    wallet_types (
        id SERIAL PRIMARY KEY,
        name VARCHAR(50) NOT NULL UNIQUE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

INSERT INTO
    wallet_types (name)
VALUES
    ('Dana'),
    ('GoPay'),
    ('ShopeePay');

-- migrate:down
DROP TABLE wallet_types;