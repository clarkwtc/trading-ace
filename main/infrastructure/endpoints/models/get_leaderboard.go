package models

import (
    "sort"
    "tradingACE/main/trading"
    "tradingACE/main/trading/events"
)

type LeaderboardViewModel struct {
    Address string `json:"address"`
    Points  int    `json:"points"`
}

func ToLeaderboardViewModel(event *events.UsersEvent) []*LeaderboardViewModel {
    users := event.Users

    result := make([]*LeaderboardViewModel, 0)

    for _, user := range users {
        result = append(result, &LeaderboardViewModel{user.Address, sumTaskPoints(user)})
    }
    sortUserByPoints(result)
    return result
}

func sumTaskPoints(user *trading.User) int {
    points := 0
    for _, taskRecord := range user.GetTaskRecords() {
        points += taskRecord.EarnPoints
    }

    return points
}

func sortUserByPoints(result []*LeaderboardViewModel) {
    sort.Slice(result, func(i, j int) bool {
        return result[i].Points > result[j].Points
    })
}
