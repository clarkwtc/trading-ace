package trading

import (
    "github.com/google/uuid"
    "log"
    "math/big"
    "strings"
    "time"
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
        log.Printf("Not support status %v", status)
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
        log.Printf("Not support status %v", modeName)
    }
    return -1
}

type TaskRecord struct {
    Id            uuid.UUID
    User          *User
    Task          Task
    Status        TaskStatus
    SwapAmount    *big.Int
    EarnPoints    int
    CompletedTime time.Time
}

func NewTaskRecord(user *User, task Task, amount *big.Int, points int) *TaskRecord {
    taskRecord := &TaskRecord{Id: uuid.New(), User: user, Task: task, Status: OnGoing, SwapAmount: amount, EarnPoints: points}
    task.SetTaskRecord(taskRecord)
    return taskRecord
}

func (taskRecord *TaskRecord) SetEarnPoints(point int) {
    taskRecord.EarnPoints = point
    taskRecord.User.AddPoints(point)
}

func (taskRecord *TaskRecord) AddSwapAmount(amount *big.Int) {
    taskRecord.SwapAmount.Add(taskRecord.SwapAmount, amount)
}

func (taskRecord *TaskRecord) Completed() {
    taskRecord.Status = Completed
}
