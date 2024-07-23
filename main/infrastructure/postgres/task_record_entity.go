package postgres

import (
    "math/big"
    "time"
    "tradingACE/main/trading"
)

type TaskRecordEntity struct {
    Name      string
    UserId    int
    TaskId    int
    Status    string
    Amount    []byte
    Points    int
    CreatedAt time.Time
    UpdatedAt time.Time
}

func ToTaskRecord(entity *TaskRecordEntity) *trading.TaskRecord {
    return &trading.TaskRecord{Name: entity.Name, Status: trading.ParseTaskStatus(entity.Status), Amount: new(big.Int).SetBytes(entity.Amount), Points: entity.Points}
}
