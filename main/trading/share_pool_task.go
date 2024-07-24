package trading

import (
    "math/big"
)

type SharePoolTask struct {
    BaseTask
}

const SharePoolTaskName = "SharePoolTask"

func NewSharePoolTask() *SharePoolTask {
    task := &SharePoolTask{}
    task.BaseTask = BaseTask{SharePoolTaskName, 10000, task}
    return task
}

func (task *SharePoolTask) getRewardPoint(addedAmount *big.Int, allUserAmount *big.Int) int {
    totalAmountFloat := new(big.Float).SetInt(addedAmount)
    allUserAmountFloat := new(big.Float).SetInt(allUserAmount)
    sharePoolRatio := new(big.Float).Quo(totalAmountFloat, allUserAmountFloat)
    rewardPoint := new(big.Float).SetInt64(int64(task.RewardPoint))

    result, _ := new(big.Float).Mul(sharePoolRatio, rewardPoint).Int64()
    return int(result)
}

func (task *SharePoolTask) Complete(user *User, allUserAmount *big.Int) {
    if task.IsTargetTask(user) && task.isCompletedPredecessorTasks(user) {
        task.reward(user, allUserAmount)
    }
}

func (task *SharePoolTask) isCompletedPredecessorTasks(user *User) bool {
    for _, taskRecord := range user.Tasks {
        if taskRecord.Name == OnBoardingTaskName && taskRecord.Status == Completed {
            return true
        }
    }
    return false
}

func (task *SharePoolTask) reward(user *User, allUserAmount *big.Int) {
    rewardPoint := task.getRewardPoint(user.TotalAmount, allUserAmount)
    user.AddPoints(task.Name, rewardPoint)
    user.AddRewardRecord(task.Name, rewardPoint)
    user.CompleteTask(task.Name)
}
