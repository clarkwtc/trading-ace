package trading

import (
    "math/big"
    "tradingACE/main/infrastructure/utils"
)

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

func (task *OnBoardingTask) Complete(user *User, amount *big.Int) {
    if task.IsTargetTask(user) && task.isRequiredAmount(user, amount) {
        task.reward(user)
    }
}

func (task *OnBoardingTask) isRequiredAmount(user *User, amount *big.Int) bool {
    sumAmount := new(big.Int).Add(user.TotalAmount, amount)
    requiredAmount := new(big.Int).SetInt64(requiredSwapAmount)
    return sumAmount.Cmp(utils.ToUSDC(requiredAmount)) >= 0
}

func (task *OnBoardingTask) reward(user *User) {
    user.AddPoints(task.Name, task.getRewardPoint())
    user.AddRewardRecord(task.Name, task.getRewardPoint())
    user.CompleteTask(task.Name)
    user.NextTask(task.Name, SharePoolTaskName)
}
