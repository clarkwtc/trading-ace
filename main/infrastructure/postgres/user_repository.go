package postgres

import (
    "database/sql"
    "log"
    "tradingACE/main/trading"
)

type UserRepository struct {
    client *sql.DB
}

func NewUserRepository(client *sql.DB) *UserRepository {
    return &UserRepository{client}
}

func (repository *UserRepository) FindUserTasksByAddress(address string) *trading.User {
    query := `
    SELECT t.name as name, tr.status as status, tr.amount as amount, tr.points as points FROM TaskRecord tr
    LEFT JOIN User0 u ON tr.user_id = u.id
    LEFT JOIN Task t ON tr.task_id = t.id
    WHERE u.address = $1
    Limit 100
    `
    rows, err := repository.client.Query(query, address)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    user := trading.NewUser(address)
    var taskRecords []*trading.TaskRecord

    for rows.Next() {
        taskRecord := &TaskRecordEntity{}
        err := rows.Scan(&taskRecord.Name, &taskRecord.Status, &taskRecord.Amount, &taskRecord.Points)
        if err != nil {
            log.Fatal(err)
        }
        taskRecords = append(taskRecords, ToTaskRecord(taskRecord))
    }
    user.Tasks = taskRecords

    return user
}
