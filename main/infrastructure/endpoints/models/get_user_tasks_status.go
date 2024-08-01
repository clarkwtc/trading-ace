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
    if user == nil {
        return result
    }
    for _, taskRecord := range user.GetTaskRecords() {
        item := &TaskStatusViewModel{
            Name: taskRecord.Task.GetName(), Status: trading.ParseTaskStatusName(taskRecord.Status),
            Amount: int(taskRecord.SwapAmount.Int64()), Points: taskRecord.EarnPoints}

        result = append(result, item)
    }
    return result
}
