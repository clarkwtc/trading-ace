package test

import (
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "math/big"
    "testing"
    "tradingACE/main/infrastructure/utils"
    "tradingACE/main/trading"
)

func TestNewOnBoardingTask(t *testing.T) {
    // Given
    // When
    task := trading.NewOnBoardingTask()

    // Then
    assert.Equal(t, "OnBoardingTask", task.Name)
    assert.Equal(t, 100, task.RewardPoint)
}

func TestCompleteOnNotEnoughAmount(t *testing.T) {
    // Given
    user := trading.NewUser(uuid.New().String())
    onBoardingTask := user.GetTaskRecords()[0].Task.(*trading.OnBoardingTask)
    amount := utils.ToUSDC(new(big.Int).SetInt64(100))

    // When
    onBoardingTask.Complete(amount)

    // Then
    taskRecord := onBoardingTask.TaskRecord
    assert.Equal(t, trading.OnBoardingTaskName, taskRecord.Task.GetName())
    assert.Equal(t, trading.OnGoing, taskRecord.Status)
    assert.Equal(t, new(big.Int).SetInt64(0), taskRecord.SwapAmount)
    assert.Equal(t, 0, taskRecord.EarnPoints)
}

func TestCompleteOnBoardingTask(t *testing.T) {
    // Given
    user := trading.NewUser(uuid.New().String())
    onBoardingTask := user.GetTaskRecords()[0].Task.(*trading.OnBoardingTask)
    amount := utils.ToUSDC(new(big.Int).SetInt64(1000))

    // When
    onBoardingTask.Complete(amount)

    // Then
    onBoardingTaskRecord := onBoardingTask.TaskRecord
    assert.Equal(t, 2, user.CountTaskRecord())
    assert.Equal(t, new(big.Int).SetInt64(0), user.TotalAmount)
    assert.Equal(t, 100, user.TotalPoints)

    assert.Equal(t, trading.OnBoardingTaskName, onBoardingTaskRecord.Task.GetName())
    assert.Equal(t, trading.Completed, onBoardingTaskRecord.Status)
    assert.Equal(t, new(big.Int).SetInt64(0), onBoardingTaskRecord.SwapAmount)
    assert.Equal(t, 100, onBoardingTaskRecord.EarnPoints)

    sharePoolTaskRecord := user.GetTaskRecords()[1].Task.(*trading.SharePoolTask).TaskRecord
    assert.Equal(t, trading.SharePoolTaskName, sharePoolTaskRecord.Task.GetName())
    assert.Equal(t, trading.OnGoing, sharePoolTaskRecord.Status)
    assert.Equal(t, new(big.Int).SetInt64(0), sharePoolTaskRecord.SwapAmount)
    assert.Equal(t, 0, sharePoolTaskRecord.EarnPoints)
}
