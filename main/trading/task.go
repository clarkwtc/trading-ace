package trading

type Task interface {
    IsTargetTask(user *User) bool
}

type BaseTask struct {
    Name        string
    RewardPoint int
    Task
}

func (task *BaseTask) IsTargetTask(user *User) bool {
    return user.CurrentTask.Name == task.Name && user.CurrentTask.Status == OnGoing
}
