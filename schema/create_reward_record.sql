CREATE TABLE RewardRecord
(
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER NOT NULL,
    task_id    INTEGER NOT NULL,
    points     INT     NOT NULL,
    created_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES User0 (id),
    FOREIGN KEY (task_id) REFERENCES Task (id)
);