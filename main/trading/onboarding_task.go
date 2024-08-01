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
    task.BaseTask = BaseTask{Name: OnBoardingTaskName, RewardPoint: 100, Task: task}
    return task
}

func (task *OnBoardingTask) getRewardPoint() int {
    return task.RewardPoint
}

func (task *OnBoardingTask) Complete(amount *big.Int) {
    if task.IsTargetTask() && task.isRequiredAmount(amount) {
        task.reward()
        task.TaskRecord.Completed()
        task.nextTask()
    }
}

func (task *OnBoardingTask) isRequiredAmount(amount *big.Int) bool {
    sumAmount := new(big.Int).Add(task.GetUser().TotalAmount, amount)
    requiredAmount := new(big.Int).SetInt64(requiredSwapAmount)
    return sumAmount.Cmp(utils.ToUSDC(requiredAmount)) >= 0
}

func (task *OnBoardingTask) reward() {
    task.TaskRecord.SetEarnPoints(task.getRewardPoint())
}

func (task *OnBoardingTask) nextTask() {
    user := task.GetUser()
    user.AcceptTask(NewTaskRecord(user, NewSharePoolTask(), task.TaskRecord.SwapAmount, 0))
}
