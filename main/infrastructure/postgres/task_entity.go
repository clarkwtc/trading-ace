package postgres

import "tradingACE/main/trading"

type TaskEntity struct {
    Id   string
    Name string
}

func ToTask(entity *TaskEntity) trading.Task {
    switch entity.Name {
    case trading.OnBoardingTaskName:
        return trading.NewOnBoardingTask()
    case trading.SharePoolTaskName:
        return trading.NewSharePoolTask()
    default:
        return nil
    }
}
