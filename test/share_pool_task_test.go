package test

import (
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "math/big"
    "testing"
    "tradingACE/main/infrastructure/utils"
    "tradingACE/main/trading"
)

func TestNewSharePoolTask(t *testing.T) {
    // Given
    // When
    task := trading.NewSharePoolTask()

    // Then
    assert.Equal(t, "SharePoolTask", task.Name)
    assert.Equal(t, 10000, task.RewardPoint)
}

func TestCompleteOnNotCompletedPredecessorTasks(t *testing.T) {
    // Given
    user := trading.NewUser(uuid.New().String())
    sharePoolTaskRecord := trading.NewTaskRecord(user, trading.NewSharePoolTask(), new(big.Int).SetInt64(0), 0)
    user.AcceptTask(sharePoolTaskRecord)
    sharePoolTask := user.GetTaskRecords()[1].Task.(*trading.SharePoolTask)
    newAmount := utils.ToUSDC(new(big.Int).SetInt64(1000))

    // When
    sharePoolTask.Complete(newAmount, false)

    // Then
    taskRecord := sharePoolTask.TaskRecord
    assert.Equal(t, trading.SharePoolTaskName, taskRecord.Task.GetName())
    assert.Equal(t, trading.OnGoing, taskRecord.Status)
    assert.Equal(t, new(big.Int).SetInt64(0), taskRecord.SwapAmount)
    assert.Equal(t, 0, taskRecord.EarnPoints)
}

func TestCompleteSharePoolTask(t *testing.T) {
    // Given
    user := trading.NewUser(uuid.New().String())
    onBoardingTask := user.GetTaskRecords()[0].Task.(*trading.OnBoardingTask)
    amount := utils.ToUSDC(new(big.Int).SetInt64(1000))
    onBoardingTask.Complete(amount)

    sharePoolTaskRecord1 := user.GetTaskRecords()[1]
    sharePoolTask1 := sharePoolTaskRecord1.Task.(*trading.SharePoolTask)

    // When
    sharePoolTask1.Complete(amount, false)

    // Then
    assert.Equal(t, 3, user.CountTaskRecord())
    assert.Equal(t, new(big.Int).SetInt64(0), user.TotalAmount)
    assert.Equal(t, 100, user.TotalPoints)

    assert.Equal(t, trading.SharePoolTaskName, sharePoolTaskRecord1.Task.GetName())
    assert.Equal(t, trading.Completed, sharePoolTaskRecord1.Status)
    assert.Equal(t, new(big.Int).SetInt64(0), sharePoolTaskRecord1.SwapAmount)
    assert.Equal(t, 0, sharePoolTaskRecord1.EarnPoints)

    sharePoolTaskRecord2 := user.GetTaskRecords()[2]
    assert.Equal(t, trading.SharePoolTaskName, sharePoolTaskRecord2.Task.GetName())
    assert.Equal(t, trading.OnGoing, sharePoolTaskRecord2.Status)
    assert.Equal(t, new(big.Int).SetInt64(0), sharePoolTaskRecord2.SwapAmount)
    assert.Equal(t, 0, sharePoolTaskRecord2.EarnPoints)
}
