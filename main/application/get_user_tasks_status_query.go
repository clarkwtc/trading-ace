package application

import (
    "tradingACE/main/trading"
    "tradingACE/main/trading/errors"
    "tradingACE/main/trading/events"
)

type GetUserTasksStatusQuery struct {
    UserRepository trading.UserRepository
}

func (query *GetUserTasksStatusQuery) Execute(address string) (*events.UserEvent, error) {
    user, err := query.UserRepository.FindUserTasksByAddress(address)
    if user == nil || err != nil {
        return nil, &errors.NotExistUserError{}
    }

    return &events.UserEvent{User: user}, nil
}
