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
    assert.Equal(t, 1, len(user.Tasks))
    assert.Equal(t, trading.OnBoardingTaskName, user.Tasks[0].Name)
    assert.Equal(t, trading.OnGoing, user.Tasks[0].Status)
    assert.Equal(t, 0, len(user.PointHistory))
}

func TestGetTask(t *testing.T) {
    // Given
    address := uuid.New().String()
    user := trading.NewUser(address)

    taskName := trading.OnBoardingTaskName
    taskStatus := trading.OnGoing

    // When
    task := user.GetTask(taskName, taskStatus)

    // Then
    assert.Equal(t, true, task != nil)
    assert.Equal(t, taskName, task.Name)
    assert.Equal(t, taskStatus, task.Status)
}

func TestGetTaskByName(t *testing.T) {
    // Given
    address := uuid.New().String()
    user := trading.NewUser(address)

    taskName := trading.OnBoardingTaskName

    // When
    task := user.GetTaskByName(taskName)

    // Then
    assert.Equal(t, true, task != nil)
    assert.Equal(t, taskName, task[0].Name)
}

func TestGetTaskById(t *testing.T) {
    // Given
    address := uuid.New().String()
    user := trading.NewUser(address)
    taskId := user.Tasks[0].Id

    // When
    task := user.GetTaskById(taskId)

    // Then
    assert.Equal(t, true, task != nil)
    assert.Equal(t, taskId, task.Id)
}

func TestAddPoints(t *testing.T) {
    // Given
    address := uuid.New().String()
    user := trading.NewUser(address)

    taskName := trading.OnBoardingTaskName
    points := 200

    // When
    user.AddPoints(taskName, points)

    // Then
    assert.Equal(t, points, user.TotalPoints)
    assert.Equal(t, 1, len(user.Tasks))
    assert.Equal(t, points, user.Tasks[0].Points)
}

func TestAddAmount(t *testing.T) {
    // Given
    address := uuid.New().String()
    user := trading.NewUser(address)

    taskName := trading.OnBoardingTaskName
    amount := new(big.Int).SetInt64(200)

    // When
    user.AddAmount(taskName, amount)

    // Then
    assert.Equal(t, amount, user.TotalAmount)
    assert.Equal(t, 1, len(user.Tasks))
    assert.Equal(t, amount, user.Tasks[0].Amount)
}

func TestAddRewardRecord(t *testing.T) {
    // Given
    address := uuid.New().String()
    user := trading.NewUser(address)
    taskName := trading.OnBoardingTaskName
    point := 100

    // When
    user.AddRewardRecord(taskName, point)

    // Then
    assert.Equal(t, 1, len(user.PointHistory))
    assert.Equal(t, taskName, user.PointHistory[0].TaskName)
    assert.Equal(t, point, user.PointHistory[0].Points)
}

func TestAddTaskRecord(t *testing.T) {
    // Given
    address := uuid.New().String()
    user := trading.NewUser(address)

    amount := new(big.Int).SetInt64(1)
    points := 1
    record := trading.NewSharePoolTaskRecord(amount, points)

    // When
    user.AddTask(record)

    // Then
    assert.Equal(t, 2, len(user.Tasks))
    assert.Equal(t, record, user.Tasks[1])
}

func TestCompleteTask(t *testing.T) {
    // Given
    address := uuid.New().String()
    user := trading.NewUser(address)
    taskName := trading.OnBoardingTaskName

    // When
    user.CompleteTask(taskName)

    // Then
    assert.Equal(t, trading.Completed, user.GetTask(taskName, trading.Completed).Status)
}

func TestNextTask(t *testing.T) {
    // Given
    address := uuid.New().String()
    user := trading.NewUser(address)
    previousTask := user.Tasks[0]
    user.CompleteTask(trading.OnBoardingTaskName)
    nextTask := trading.SharePoolTaskName

    // When
    user.NextTask(previousTask.Id, nextTask)

    // Then
    assert.Equal(t, 2, len(user.Tasks))
    assert.Equal(t, nextTask, user.Tasks[1].Name)
    assert.Equal(t, trading.OnGoing, user.Tasks[1].Status)
}
