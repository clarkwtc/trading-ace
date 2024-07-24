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
    assert.Equal(t, trading.OnBoardingTaskName, task.Name)
    assert.Equal(t, 100, task.RewardPoint)
}

func TestCompleteOnNotEnoughAmount(t *testing.T) {
    // Given
    user := trading.NewUser(uuid.New().String())
    task := trading.NewOnBoardingTask()
    amount := utils.ToUSDC(new(big.Int).SetInt64(100))

    // When
    task.Complete(user, amount)

    // Then
    assert.Equal(t, 1, len(user.Tasks))
    assert.Equal(t, trading.OnBoardingTaskName, user.Tasks[0].Name)
    assert.Equal(t, trading.OnGoing, user.Tasks[0].Status)
    assert.Equal(t, new(big.Int).SetInt64(0), user.Tasks[0].Amount)
    assert.Equal(t, 0, user.Tasks[0].Points)
}

func TestCompleteOnBoardingTask(t *testing.T) {
    // Given
    user := trading.NewUser(uuid.New().String())
    task := trading.NewOnBoardingTask()
    amount := utils.ToUSDC(new(big.Int).SetInt64(1000))

    // When
    task.Complete(user, amount)

    // Then
    assert.Equal(t, 2, len(user.Tasks))
    assert.Equal(t, trading.OnBoardingTaskName, user.Tasks[0].Name)
    assert.Equal(t, trading.Completed, user.Tasks[0].Status)
    assert.Equal(t, 100, user.Tasks[0].Points)

    assert.Equal(t, trading.SharePoolTaskName, user.Tasks[1].Name)
    assert.Equal(t, trading.OnGoing, user.Tasks[1].Status)
    assert.Equal(t, 0, user.Tasks[1].Points)
}
