package test

import (
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "math/big"
    "testing"
    "tradingACE/main/trading"
)

func TestParseTaskStatusName(t *testing.T) {
    tests := []struct {
        taskStatus     trading.TaskStatus
        taskStatusName string
    }{
        {trading.OnGoing, "OnGoing"},
        {trading.Completed, "Completed"},
    }

    for _, task := range tests {
        // Given
        taskStatus := task.taskStatus

        // When
        taskStatusName := trading.ParseTaskStatusName(taskStatus)

        // Then
        assert.Equal(t, task.taskStatusName, taskStatusName)
    }
}

func TestParseTaskStatus(t *testing.T) {
    tests := []struct {
        taskStatus     trading.TaskStatus
        taskStatusName string
    }{
        {trading.OnGoing, "OnGoing"},
        {trading.Completed, "Completed"},
    }

    for _, task := range tests {
        // Given
        taskStatusName := task.taskStatusName

        // When
        taskStatus := trading.ParseTaskStatus(taskStatusName)

        // Then
        assert.Equal(t, task.taskStatus, taskStatus)
    }
}

func TestNewTaskRecord(t *testing.T) {
    // Given
    user := trading.NewUser(uuid.New().String())
    task := trading.NewOnBoardingTask()
    amount := new(big.Int).SetInt64(10)
    points := 1

    // When
    taskRecord := trading.NewTaskRecord(user, task, amount, points)

    // Then
    assert.Equal(t, true, taskRecord.Id != uuid.Nil)
    assert.Equal(t, trading.OnBoardingTaskName, taskRecord.Task.GetName())
    assert.Equal(t, trading.OnGoing, taskRecord.Status)
    assert.Equal(t, amount, taskRecord.SwapAmount)
    assert.Equal(t, points, taskRecord.EarnPoints)
}

func TestSetEarnPoints(t *testing.T) {
    // Given
    amount := new(big.Int).SetInt64(0)
    user := trading.NewUser(uuid.New().String())
    task := trading.NewOnBoardingTask()
    taskRecord := trading.NewTaskRecord(user, task, amount, 0)
    points := 100

    // When
    taskRecord.SetEarnPoints(points)

    // Then
    assert.Equal(t, points, taskRecord.EarnPoints)
}

func TestAddSwapAmount(t *testing.T) {
    // Given
    amount := new(big.Int).SetInt64(10)
    points := 1
    user := trading.NewUser(uuid.New().String())
    task := trading.NewOnBoardingTask()
    taskRecord := trading.NewTaskRecord(user, task, amount, points)
    addAmount := new(big.Int).SetInt64(10)

    // When
    taskRecord.AddSwapAmount(addAmount)

    // Then
    assert.Equal(t, amount.Add(amount, addAmount), taskRecord.SwapAmount)
}

func TestCompleted(t *testing.T) {
    // Given
    amount := new(big.Int).SetInt64(10)
    points := 1
    user := trading.NewUser(uuid.New().String())
    task := trading.NewOnBoardingTask()
    taskRecord := trading.NewTaskRecord(user, task, amount, points)

    // When
    taskRecord.Completed()

    // Then
    assert.Equal(t, trading.Completed, taskRecord.Status)
}
