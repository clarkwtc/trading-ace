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
    assert.Equal(t, trading.SharePoolTaskName, task.Name)
    assert.Equal(t, 10000, task.RewardPoint)
}

func TestCompleteOnNotCompletedPredecessorTasks(t *testing.T) {
    // Given
    user := trading.NewUser(uuid.New().String())
    amount := utils.ToUSDC(new(big.Int).SetInt64(100))
    onBoardingTask := trading.NewOnBoardingTask()
    onBoardingTask.Complete(user, amount)

    task := trading.NewSharePoolTask()

    // When
    task.Complete(user, amount, false)

    // Then
    assert.Equal(t, 1, len(user.Tasks))
    assert.Equal(t, trading.OnBoardingTaskName, user.Tasks[0].Name)
    assert.Equal(t, trading.OnGoing, user.Tasks[0].Status)
    assert.Equal(t, 0, user.Tasks[0].Points)
}

func TestCompleteSharePoolTask(t *testing.T) {
    // Given
    user := trading.NewUser(uuid.New().String())
    amount := utils.ToUSDC(new(big.Int).SetInt64(1000))
    onBoardingTask := trading.NewOnBoardingTask()
    onBoardingTask.Complete(user, amount)

    task := trading.NewSharePoolTask()

    // When
    task.Complete(user, amount, false)

    // Then
    assert.Equal(t, 3, len(user.Tasks))
    assert.Equal(t, trading.OnBoardingTaskName, user.Tasks[0].Name)
    assert.Equal(t, trading.Completed, user.Tasks[0].Status)
    assert.Equal(t, 100, user.Tasks[0].Points)

    assert.Equal(t, trading.SharePoolTaskName, user.Tasks[1].Name)
    assert.Equal(t, trading.Completed, user.Tasks[1].Status)
    assert.Equal(t, 0, user.Tasks[1].Points)

    assert.Equal(t, trading.SharePoolTaskName, user.Tasks[2].Name)
    assert.Equal(t, trading.OnGoing, user.Tasks[2].Status)
    assert.Equal(t, 0, user.Tasks[2].Points)
}
