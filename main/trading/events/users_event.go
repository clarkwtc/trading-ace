package events

import "tradingACE/main/trading"

type UsersEvent struct {
	Users []*trading.User
}