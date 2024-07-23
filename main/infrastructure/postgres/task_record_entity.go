package postgres

import (
    "github.com/google/uuid"
    "math/big"
    "time"
    "tradingACE/main/trading"
)

type TaskRecordEntity struct {
    Id        string
    UserId    string
    TaskId    string
    Status    string
    Amount    string
    Points    int
    CreatedAt time.Time
    UpdatedAt time.Time
}

func ToTaskRecord(taskRecordEntity *TaskRecordEntity, taskEntity *TaskEntity) *trading.TaskRecord {
    amount, _ := new(big.Int).SetString(taskRecordEntity.Amount, 10)
    return &trading.TaskRecord{Id: uuid.MustParse(taskRecordEntity.Id), Name: taskEntity.Name, Status: trading.ParseTaskStatus(taskRecordEntity.Status), Amount: amount, Points: taskRecordEntity.Points}
}

func ToTaskRecordEntity(user *trading.User, task *trading.TaskRecord, tasksMap map[string]*TaskEntity) *TaskRecordEntity {
    now := time.Now()
    return &TaskRecordEntity{task.Id.String(), user.Id.String(), tasksMap[task.Name].Id, trading.ParseTaskStatusName(task.Status), task.Amount.String(), task.Points, now, now}
}
