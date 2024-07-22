package trading

import (
    "math/big"
)

type SharePoolTask struct {
    BaseTask
}

const sharePoolTaskName = "SharePoolTask"

func NewSharePoolTask() *SharePoolTask {
    task := &SharePoolTask{}
    task.BaseTask = BaseTask{sharePoolTaskName, 10000, task}
    return task
}

func (task *SharePoolTask) getRewardPoint(addedAmount *big.Int, allUsersSwapAmount *big.Int) int {
    totalAmountFloat := new(big.Float).SetInt(addedAmount)
    allUsersSwapAmountFloat := new(big.Float).SetInt(allUsersSwapAmount)
    sharePoolRatio := new(big.Float).Quo(totalAmountFloat, allUsersSwapAmountFloat)
    rewardPoint := new(big.Float).SetInt64(int64(task.RewardPoint))

    result, _ := new(big.Float).Mul(sharePoolRatio, rewardPoint).Int64()
    return int(result)
}

func (task *SharePoolTask) Complete(user *User, allUsersSwapAmount *big.Int) {
    if task.IsTargetTask(user) && task.isCompletedPredecessorTasks(user) {
        task.reward(user, allUsersSwapAmount)
    }
}

func (task *SharePoolTask) isCompletedPredecessorTasks(user *User) bool {
    for _, record := range user.PointHistory {
        if record.TaskName == "OnBoardingTask" {
            return true
        }
    }
    return false
}

func (task *SharePoolTask) reward(user *User, allUsersSwapAmount *big.Int) {
    diffAmount := new(big.Int).Sub(user.TotalAmount, user.GetTotalAmountOfPointHistory())
    rewardPoint := task.getRewardPoint(diffAmount, allUsersSwapAmount)
    user.AddPoints(rewardPoint)
    user.AddRewardRecord(task.Name, rewardPoint, diffAmount)
}
