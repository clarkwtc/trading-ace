package events

import (
    "math/big"
    "tradingACE/main/trading"
)

type SettlementPointsUser struct {
    Id          string
    Address     string
    TotalAmount *big.Int
    TotalPoints int
}

type SettlementPointsEvent struct {
    Users []*SettlementPointsUser
    Count int
}

func NewSettlementPointsEvent(users []*trading.User) *SettlementPointsEvent {
    var settlementPointsUsers []*SettlementPointsUser
    for _, user := range users {
        settlementPointsUsers = append(settlementPointsUsers, &SettlementPointsUser{user.Id.String(), user.Address, user.TotalAmount, user.TotalPoints})
    }
    return &SettlementPointsEvent{settlementPointsUsers, len(settlementPointsUsers)}
}
