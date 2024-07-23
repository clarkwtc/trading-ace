package trading

import "math/big"

type OnBoardingTask struct {
    BaseTask
}

const OnBoardingTaskName = "OnBoardingTask"
const requiredSwapAmount = 1000

func NewOnBoardingTask() *OnBoardingTask {
    task := &OnBoardingTask{}
    task.BaseTask = BaseTask{OnBoardingTaskName, 100, task}
    return task
}

func (task *OnBoardingTask) getRewardPoint() int {
    return task.RewardPoint
}

func (task *OnBoardingTask) Complete(user *User, event *Event) {
    if task.IsTargetTask(user) && task.isRequiredAmount(user, event) {
        task.reward(user)
    }
}

func (task *OnBoardingTask) isRequiredAmount(user *User, event *Event) bool {
    sumAmount := new(big.Int).Add(user.TotalAmount, event.Amount0Out)
    requiredAmount := new(big.Int).SetInt64(requiredSwapAmount)
    return sumAmount.Cmp(requiredAmount) >= 0
}

func (task *OnBoardingTask) reward(user *User) {
    user.AddPoints(task.Name, task.getRewardPoint())
    user.AddRewardRecord(task.Name, task.getRewardPoint())
    user.CompleteTask(task.Name)
    user.NextTask(task.Name, SharePoolTaskName)
}
