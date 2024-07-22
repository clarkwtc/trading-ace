package trading

import "math/big"

type User struct {
    Address      string
    CurrentTask  *TaskProcessing
    TotalPoints  int
    TotalAmount  *big.Int
    PointHistory []*RewardRecord
}

func NewUser(adrress string, currentTask *TaskProcessing) *User {
    return &User{Address: adrress, CurrentTask: currentTask, TotalPoints: 0, TotalAmount: &big.Int{}}
}

func (user *User) AddPoints(point int) {
    user.TotalPoints += point
}

func (user *User) AddAmount(amount *big.Int) {
    user.TotalAmount = new(big.Int).Add(user.TotalAmount, amount)
}

func (user *User) AddRewardRecord(taskName string, point int, amount *big.Int) {
    user.PointHistory = append(user.PointHistory, &RewardRecord{taskName, point, amount})
}

func (user *User) GetTotalAmountOfPointHistory() *big.Int {
    totalAmount := new(big.Int).SetInt64(0)
    for _, record := range user.PointHistory {
        totalAmount = totalAmount.Add(totalAmount, record.Amount)
    }
    return totalAmount
}

func (user *User) ChangeSharePoolTask(taskName string) {
    user.CurrentTask = &TaskProcessing{taskName, OnGoing}
}

func (user *User) CompleteTask() {
    user.CurrentTask.Status = Completed
}
