package models

import (
    "tradingACE/main/trading"
    "tradingACE/main/trading/events"
)

type TaskStatusViewModel struct {
    Name   string `json:"name"`
    Status string `json:"status"`
    Amount int    `json:"amount"`
    Points int    `json:"points"`
}

func ToTaskStatusViewModel(event *events.UserEvent) []*TaskStatusViewModel {
    user := event.User
    result := make([]*TaskStatusViewModel, 0)
    for _, task := range user.Tasks {
        item := &TaskStatusViewModel{
            Name: task.Name, Status: trading.ParseTaskStatusName(task.Status),
            Amount: int(task.Amount.Int64()), Points: task.Points}

        result = append(result, item)
    }
    return result
}
