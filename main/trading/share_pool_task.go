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
    task.BaseTask = BaseTask{Name: SharePoolTaskName, RewardPoint: 10000, Task: task}
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

func (task *SharePoolTask) Complete(allUserAmount *big.Int, final bool) {
    if task.IsTargetTask() && task.isCompletedPredecessorTasks() {
        task.reward(allUserAmount)
        task.TaskRecord.Completed()
        if !final {
            task.nextTask()
        }
    }
}

func (task *SharePoolTask) isCompletedPredecessorTasks() bool {
    return task.GetUser().GetTaskRecord(OnBoardingTaskName, Completed) != nil
}

func (task *SharePoolTask) reward(allUserAmount *big.Int) {
    rewardPoint := task.getRewardPoint(task.GetUser().TotalAmount, allUserAmount)
    task.TaskRecord.SetEarnPoints(rewardPoint)
}

func (task *SharePoolTask) nextTask() {
    user := task.GetUser()
    user.AcceptTask(NewTaskRecord(user, NewSharePoolTask(), new(big.Int).SetInt64(0), 0))
}
