package application

import (
    "tradingACE/main/trading"
    "tradingACE/main/trading/events"
)

type GetUserPointsHistoryQuery struct {
    UserRepository trading.UserRepository
}

func (usecase *GetUserPointsHistoryQuery) Execute(address string) *events.UserEvent {
    user := usecase.UserRepository.FindUserRewardByAddress(address)
    if user == nil {
        return nil
    }

    return &events.UserEvent{User: user}
}
