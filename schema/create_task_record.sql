CREATE TABLE TaskRecord
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    user_id    INTEGER      NOT NULL,
    task_id    INTEGER      NOT NULL,
    status     VARCHAR(255) NOT NULL,
    amount     NUMERIC(38, 0) NULL,
    points     INTEGER      NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES User0 (id),
    FOREIGN KEY (task_id) REFERENCES Task (id)
);