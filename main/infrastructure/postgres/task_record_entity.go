package postgres

import (
    "math/big"
    "time"
    "tradingACE/main/trading"
)

type TaskRecordEntity struct {
    Id        int
    UserId    int
    TaskId    int
    Status    string
    Amount    []byte
    Points    int
    CreatedAt time.Time
    UpdatedAt time.Time
}

func ToTaskRecord(taskRecordEntity *TaskRecordEntity, taskEntity *TaskEntity) *trading.TaskRecord {
    return &trading.TaskRecord{Name: taskEntity.Name, Status: trading.ParseTaskStatus(taskRecordEntity.Status), Amount: new(big.Int).SetBytes(taskRecordEntity.Amount), Points: taskRecordEntity.Points}
}
