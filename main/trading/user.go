package trading

import (
    "github.com/google/uuid"
    "math/big"
    "time"
)

type User struct {
    Id           uuid.UUID
    Address      string
    TotalAmount  *big.Int
    TotalPoints  int
    Tasks        []*TaskRecord
    PointHistory []*RewardRecord
}

func NewUser(adrress string) *User {
    tasks := make([]*TaskRecord, 0)
    tasks = append(tasks, NewOnBoardingTaskRecord(new(big.Int).SetInt64(0), 0))
    return &User{Id: uuid.New(), Address: adrress, TotalAmount: &big.Int{}, TotalPoints: 0, Tasks: tasks}
}

func (user *User) GetTask(name string, status TaskStatus) *TaskRecord {
    for _, task := range user.Tasks {
        if task.Name == name && task.Status == status {
            return task
        }
    }
    return nil
}

func (user *User) AddPoints(taskName string, point int) {
    user.TotalPoints += point
    task := user.GetTask(taskName, OnGoing)
    if task == nil {
        return
    }
    task.AddPoints(point)
}

func (user *User) AddAmount(taskName string, amount *big.Int) {
    user.TotalAmount = new(big.Int).Add(user.TotalAmount, amount)
    task := user.GetTask(taskName, OnGoing)
    if task == nil {
        return
    }
    task.AddAmount(amount)
}

func (user *User) AddRewardRecord(taskName string, point int) {
    user.PointHistory = append(user.PointHistory, &RewardRecord{uuid.New(), taskName, point, time.Now()})
}

func (user *User) AddTask(taskRecord *TaskRecord) {
    user.Tasks = append(user.Tasks, taskRecord)
}

func (user *User) NextTask(previousTask string, newTask string) {
    task := user.GetTask(previousTask, Completed)
    if task == nil {
        return
    }
    user.Tasks = append(user.Tasks, &TaskRecord{uuid.New(), newTask, OnGoing, task.Amount, 0})
}

func (user *User) CompleteTask(taskName string) {
    task := user.GetTask(taskName, OnGoing)
    if task == nil {
        return
    }
    task.Status = Completed
}
