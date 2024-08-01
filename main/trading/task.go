package trading

type Task interface {
    GetName() string
    GetUser() *User
    SetTaskRecord(taskRecord *TaskRecord)
    IsTargetTask() bool
}

type BaseTask struct {
    Name        string
    RewardPoint int
    TaskRecord  *TaskRecord
    Task
}

func (task *BaseTask) GetName() string {
    return task.Name
}

func (task *BaseTask) GetUser() *User {
    return task.TaskRecord.User
}

func (task *BaseTask) SetTaskRecord(taskRecord *TaskRecord) {
    task.TaskRecord = taskRecord
}

func (task *BaseTask) IsTargetTask() bool {
    return task.GetUser().GetTaskRecord(task.Name, OnGoing) != nil
}
