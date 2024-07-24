CREATE TABLE User0
(
    id         UUID PRIMARY KEY,
    address    VARCHAR(255) NOT NULL UNIQUE,
    amount     NUMERIC(38, 0) NULL,
    points     INTEGER      NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);