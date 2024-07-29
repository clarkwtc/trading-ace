package events

import (
    "math/big"
    "tradingACE/main/trading"
)

type SwapEvent struct {
    Id          string
    Address     string
    SwapAmount  *big.Int
    TotalAmount *big.Int
    TotalPoints int
}

func NewSwapEvent(user *trading.User, amount *big.Int) *SwapEvent {
    return &SwapEvent{user.Id.String(), user.Address, amount, user.TotalAmount, user.TotalPoints}
}
