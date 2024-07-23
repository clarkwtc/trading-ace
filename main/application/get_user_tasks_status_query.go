package application

import (
    "tradingACE/main/trading"
    "tradingACE/main/trading/events"
)

type GetUserTasksStatusQuery struct {
    UserRepository trading.UserRepository
}

func (usecase *GetUserTasksStatusQuery) Execute(address string) *events.UserEvent {
    user := usecase.UserRepository.FindUserTasksByAddress(address)
    if user == nil {
        return nil
    }

    return &events.UserEvent{User: user}
}
