package events

import "tradingACE/main/trading"

type UserEvent struct {
    User *trading.User
}