package test

import (
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "math/big"
    "testing"
    "tradingACE/main/trading"
)

func TestNewUser(t *testing.T) {
    // Given
    address := uuid.New().String()

    // When
    user := trading.NewUser(address)

    // Then
    assert.Equal(t, true, user.Id != uuid.Nil)
    assert.Equal(t, address, user.Address)
    assert.Equal(t, new(big.Int).SetInt64(0), user.TotalAmount)
    assert.Equal(t, 0, user.TotalPoints)
    assert.Equal(t, 1, user.CountTaskRecord())
    taskRecord := user.GetTaskRecords()[0]
    assert.Equal(t, trading.OnBoardingTaskName, taskRecord.Task.(*trading.OnBoardingTask).Name)
    assert.Equal(t, trading.OnGoing, taskRecord.Status)
    assert.Equal(t, 0, taskRecord.EarnPoints)
    assert.Equal(t, new(big.Int).SetInt64(0), taskRecord.SwapAmount)
}

func TestGetTask(t *testing.T) {
    // Given
    address := uuid.New().String()
    user := trading.NewUser(address)

    taskName := trading.OnBoardingTaskName
    taskStatus := trading.OnGoing

    // When
    taskRecord := user.GetTaskRecord(taskName, taskStatus)

    // Then
    assert.Equal(t, true, taskRecord != nil)
    assert.Equal(t, taskName, taskRecord.Task.GetName())
    assert.Equal(t, taskStatus, taskRecord.Status)
}

func TestGetTaskByName(t *testing.T) {
    // Given
    address := uuid.New().String()
    user := trading.NewUser(address)

    taskName := trading.OnBoardingTaskName

    // When
    taskRecord := user.GetTaskRecordByName(taskName)

    // Then
    assert.Equal(t, true, taskRecord != nil)
    assert.Equal(t, taskName, taskRecord[0].Task.GetName())
}

func TestAddPoints(t *testing.T) {
    // Given
    address := uuid.New().String()
    user := trading.NewUser(address)
    points := 200

    // When
    user.AddPoints(points)

    // Then
    assert.Equal(t, points, user.TotalPoints)
}

func TestAddAmount(t *testing.T) {
    // Given
    address := uuid.New().String()
    user := trading.NewUser(address)

    amount := new(big.Int).SetInt64(200)

    // When
    user.AddAmount(amount)

    // Then
    assert.Equal(t, amount, user.TotalAmount)
}

//func TestAddTaskRecord(t *testing.T) {
//    // Given
//    address := uuid.New().String()
//    user := trading.NewUser(address)
//
//    amount := new(big.Int).SetInt64(1)
//    points := 1
//    record := trading.NewSharePoolTaskRecord(amount, points)
//
//    // When
//    user.AddTask(record)
//
//    // Then
//    assert.Equal(t, 2, len(user.taskRecords))
//    assert.Equal(t, record, user.taskRecords[1])
//}
