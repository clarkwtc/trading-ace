CREATE TABLE IF NOT EXISTS TaskRecord
(
    id         uuid primary key,
    user_id    uuid         NOT NULL,
    task_id    uuid         NOT NULL,
    status     VARCHAR(255) NOT NULL,
    amount     NUMERIC(38, 0) NULL,
    points     INTEGER      NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES User0 (id),
    FOREIGN KEY (task_id) REFERENCES Task (id)
);