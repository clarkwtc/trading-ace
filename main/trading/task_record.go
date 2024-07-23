package trading

import (
    "github.com/google/uuid"
    "log"
    "math/big"
    "strings"
)

type TaskStatus int

const (
    OnGoing TaskStatus = iota
    Completed
)

func ParseTaskStatusName(status TaskStatus) string {
    switch status {
    case OnGoing:
        return "OnGoing"
    case Completed:
        return "Completed"
    default:
        log.Fatalf("Not support status %v", status)
    }
    return ""
}

func ParseTaskStatus(modeName string) TaskStatus {
    modeName = strings.ToLower(modeName)
    switch modeName {
    case strings.ToLower("OnGoing"):
        return OnGoing
    case strings.ToLower("Completed"):
        return Completed
    default:
        log.Fatalf("Not support status %v", modeName)
    }
    return -1
}

type TaskRecord struct {
    Id     uuid.UUID
    Name   string
    Status TaskStatus
    Amount *big.Int
    Points int
}

func NewOnBoardingTaskRecord(amount *big.Int, points int) *TaskRecord {
    return &TaskRecord{uuid.New(), OnBoardingTaskName, OnGoing, amount, points}
}

func NewSharePoolTaskRecord(amount *big.Int, points int) *TaskRecord {
    return &TaskRecord{uuid.New(), SharePoolTaskName, OnGoing, amount, points}
}

func (taskRecord *TaskRecord) Clone() *TaskRecord {
    return &TaskRecord{
        Id:     taskRecord.Id,
        Name:   taskRecord.Name,
        Status: taskRecord.Status,
        Amount: taskRecord.Amount,
        Points: taskRecord.Points,
    }
}

func (taskRecord *TaskRecord) AddPoints(point int) {
    taskRecord.Points += point
}

func (taskRecord *TaskRecord) AddAmount(amount *big.Int) {
    taskRecord.Amount = new(big.Int).Add(taskRecord.Amount, amount)
}
