-- migrate:up
CREATE TABLE
    users (
        id SERIAL PRIMARY KEY,
        username VARCHAR(100) NOT NULL,
        email VARCHAR(100) NOT NULL,
        password VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        login_at TIMESTAMP DEFAULT NULL
    );

INSERT INTO
    users (username, email, password)
VALUES
    (
        'admin',
        'admin@example.com',
        '$2a$12$THUvcECQS8Wyh934GvnTweT4yhAUcETMiJbXG1cpeU2kPPjEh.0Ny'
    );

-- migrate:down
drop table users;