package application

import (
    "tradingACE/main/trading"
    "tradingACE/main/trading/events"
)

type GetUserTasksStatusQuery struct {
    UserRepository trading.UserRepository
}

func (query *GetUserTasksStatusQuery) Execute(address string) *events.UserEvent {
    user := query.UserRepository.FindUserTasksByAddress(address)
    if user == nil {
        return &events.UserEvent{}
    }

    return &events.UserEvent{User: user}
}
