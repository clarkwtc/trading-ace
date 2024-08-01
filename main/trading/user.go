package trading

import (
    "github.com/google/uuid"
    "math/big"
    "strings"
)

type User struct {
    Id          uuid.UUID
    Address     string
    TotalAmount *big.Int
    TotalPoints int
    taskRecords []*TaskRecord
}

func NewUser(adrress string) *User {
    tasks := make([]*TaskRecord, 0)
    user := &User{Id: uuid.New(), Address: adrress, TotalAmount: &big.Int{}, TotalPoints: 0, taskRecords: tasks}
    user.taskRecords = append(tasks, NewTaskRecord(user, NewOnBoardingTask(), new(big.Int).SetInt64(0), 0))
    return user
}

func (user *User) CountTaskRecord() int {
    return len(user.taskRecords)
}

func (user *User) GetTaskRecords() []*TaskRecord {
    return user.taskRecords
}

func (user *User) SetTaskRecords(taskRecords []*TaskRecord){
    user.taskRecords = taskRecords
}

func (user *User) GetTaskRecord(name string, status TaskStatus) *TaskRecord {
    for _, taskRecord := range user.taskRecords {
        if task, ok := taskRecord.Task.(Task); ok {
            if task.GetName() == name && taskRecord.Status == status {
                return taskRecord
            }
        }
    }
    return nil
}

func (user *User) GetTaskRecordByStatus(status TaskStatus) *TaskRecord {
    for _, taskRecord := range user.taskRecords {
        if taskRecord.Status == status {
            return taskRecord
        }
    }
    return nil
}

func (user *User) GetTaskRecordByName(name string) []*TaskRecord {
    var taskRecords []*TaskRecord
    for _, taskRecord := range user.taskRecords {
        if task, ok := taskRecord.Task.(Task); ok {
            if strings.ToLower(task.GetName()) == strings.ToLower(name) {
                taskRecords = append(taskRecords, taskRecord)
            }
        }
    }
    return taskRecords
}

func (user *User) AddPoints(point int) {
    user.TotalPoints += point
}

func (user *User) AddAmount(amount *big.Int) {
    user.TotalAmount = new(big.Int).Add(user.TotalAmount, amount)
}

func (user *User) AcceptTask(taskRecord *TaskRecord) {
    user.taskRecords = append(user.taskRecords, taskRecord)
}
