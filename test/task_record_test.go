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

func TestNewOnBoardingTaskRecord(t *testing.T) {
    // Given
    amount := new(big.Int).SetInt64(10)
    points := 1

    // When
    taskRecord := trading.NewOnBoardingTaskRecord(amount, points)

    // Then
    assert.Equal(t, true, taskRecord.Id != uuid.Nil)
    assert.Equal(t, trading.OnBoardingTaskName, taskRecord.Name)
    assert.Equal(t, trading.OnGoing, taskRecord.Status)
    assert.Equal(t, amount, taskRecord.Amount)
    assert.Equal(t, points, taskRecord.Points)
}

func TestNewSharePoolTaskRecord(t *testing.T) {
    // Given
    amount := new(big.Int).SetInt64(10)
    points := 1

    // When
    taskRecord := trading.NewSharePoolTaskRecord(amount, points)

    // Then
    assert.Equal(t, true, taskRecord.Id != uuid.Nil)
    assert.Equal(t, trading.SharePoolTaskName, taskRecord.Name)
    assert.Equal(t, trading.OnGoing, taskRecord.Status)
    assert.Equal(t, amount, taskRecord.Amount)
    assert.Equal(t, points, taskRecord.Points)
}

func TestTaskRecordAddAmount(t *testing.T) {
    // Given
    amount := new(big.Int).SetInt64(10)
    points := 1
    taskRecord := trading.NewOnBoardingTaskRecord(amount, points)
    addAmount := new(big.Int).SetInt64(10)
    // When
    taskRecord.AddAmount(addAmount)

    // Then
    assert.Equal(t, amount.Add(amount, addAmount), taskRecord.Amount)
}

func TestTaskRecordAddPoints(t *testing.T) {
    // Given
    amount := new(big.Int).SetInt64(10)
    points := 1
    taskRecord := trading.NewOnBoardingTaskRecord(amount, points)
    addPoint := 10
    // When
    taskRecord.AddPoints(addPoint)

    // Then
    assert.Equal(t, points+addPoint, taskRecord.Points)
}
