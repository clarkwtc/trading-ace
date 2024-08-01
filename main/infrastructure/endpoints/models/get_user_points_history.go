package models

import (
    "time"
    "tradingACE/main/trading/events"
)

type PointsHistoryViewModel struct {
    TaskName   string    `json:"task_name"`
    Points     int       `json:"points"`
    RewardTime time.Time `json:"reward_time"`
}

func ToPointsHistoryViewModel(event *events.UserEvent) []*PointsHistoryViewModel {
    user := event.User
    result := make([]*PointsHistoryViewModel, 0)
    if user == nil {
        return result
    }
    for _, taskRecord := range user.GetTaskRecords() {
        item := &PointsHistoryViewModel{taskRecord.Task.GetName(), taskRecord.EarnPoints, taskRecord.CompletedTime}

        result = append(result, item)
    }
    return result
}
