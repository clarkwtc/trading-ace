package application

import (
    "log"
    "tradingACE/main/trading"
    "tradingACE/main/trading/events"
)

type GetLeaderboardQuery struct {
    UserRepository trading.UserRepository
}

func (query *GetLeaderboardQuery) Execute(taskName string) *events.UsersEvent {
    users, err := query.UserRepository.FindAllUserTasks()
    if err != nil {
        log.Printf("FindAllUserTasks fail: %v", err)
        return &events.UsersEvent{}
    }

    if len(users) == 0 {
        return &events.UsersEvent{}
    }

    filterTaskName(users, taskName)

    return &events.UsersEvent{Users: users}
}

func filterTaskName(users []*trading.User, taskName string) {
    for _, user := range users {
        user.SetTaskRecords(user.GetTaskRecordByName(taskName))
    }
}
