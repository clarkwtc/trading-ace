CREATE TABLE RewardRecord
(
    id         uuid primary key,
    user_id    uuid NOT NULL,
    task_id    uuid NOT NULL,
    points     INT  NOT NULL,
    created_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES User0 (id),
    FOREIGN KEY (task_id) REFERENCES Task (id)
);