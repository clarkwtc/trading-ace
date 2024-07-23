package trading

import "math/big"

type User struct {
    Address      string
    TotalAmount  *big.Int
    TotalPoints  int
    CurrentTask  *TaskRecord
    Tasks        []*TaskRecord
    PointHistory []*RewardRecord
}

func NewUser(adrress string) *User {
    return &User{Address: adrress, TotalAmount: &big.Int{}, TotalPoints: 0, CurrentTask: NewOnBoardingTaskRecord(&big.Int{}, 0)}
}

func (user *User) AddPoints(point int) {
    user.TotalPoints += point
    if user.CurrentTask.Status == OnGoing {
        user.CurrentTask.Points += user.CurrentTask.Points
    }
}

func (user *User) AddAmount(amount *big.Int) {
    user.TotalAmount = new(big.Int).Add(user.TotalAmount, amount)
    if user.CurrentTask.Status == OnGoing {
        user.CurrentTask.Amount = new(big.Int).Add(user.CurrentTask.Amount, amount)
    }
}

func (user *User) AddRewardRecord(taskName string, point int) {
    user.PointHistory = append(user.PointHistory, &RewardRecord{taskName, point})
}

func (user *User) ChangeSharePoolTask() {
    if user.CurrentTask.Name == SharePoolTaskName {
        return
    }
    user.CurrentTask = NewSharePoolTaskRecord(user.CurrentTask.Amount, user.CurrentTask.Points)
}

func (user *User) CompleteTask() {
    user.CurrentTask.Status = Completed
    user.Tasks = append(user.Tasks, user.CurrentTask.Clone())
}
