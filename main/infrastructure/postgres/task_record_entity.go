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

func ToTaskRecord(taskRecordEntity *TaskRecordEntity, user *trading.User, task trading.Task) *trading.TaskRecord {
    amount, _ := new(big.Int).SetString(taskRecordEntity.Amount, 10)
    taskRecord := &trading.TaskRecord{Id: uuid.MustParse(taskRecordEntity.Id), User: user, Task: task, Status: trading.ParseTaskStatus(taskRecordEntity.Status), SwapAmount: amount, EarnPoints: taskRecordEntity.Points, CompletedTime: taskRecordEntity.UpdatedAt}
    task.SetTaskRecord(taskRecord)
    return taskRecord
}

func ToTaskRecordEntity(taskRecord *trading.TaskRecord, user *trading.User, tasksMap map[string]*TaskEntity) *TaskRecordEntity {
    now := time.Now()
    return &TaskRecordEntity{taskRecord.Id.String(), user.Id.String(), tasksMap[taskRecord.Task.GetName()].Id, trading.ParseTaskStatusName(taskRecord.Status), taskRecord.SwapAmount.String(), taskRecord.EarnPoints, now, now}
}
