package trading

import (
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
    Name   string
    Status TaskStatus
    Amount *big.Int
    Points int
}

func NewOnBoardingTaskRecord(amount *big.Int, points int) *TaskRecord {
    return &TaskRecord{OnBoardingTaskName, OnGoing, amount, points}
}

func NewSharePoolTaskRecord(amount *big.Int, points int) *TaskRecord {
    return &TaskRecord{SharePoolTaskName, OnGoing, amount, points}
}

func (t *TaskRecord) Clone() *TaskRecord {
    return &TaskRecord{
        Name:   t.Name,
        Status: t.Status,
        Amount: t.Amount,
        Points: t.Points,
    }
}
