package application

import (
    "tradingACE/main/trading"
    "tradingACE/main/trading/events"
)

type GetUserPointsHistoryQuery struct {
    UserRepository trading.UserRepository
}

func (query *GetUserPointsHistoryQuery) Execute(address string) *events.UserEvent {
    user := query.UserRepository.FindUserRewardByAddress(address)
    if user == nil {
        return nil
    }

    return &events.UserEvent{User: user}
}
