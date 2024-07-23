package postgres

import (
    "github.com/google/uuid"
    "time"
    "tradingACE/main/trading"
)

type RewardRecordEntity struct {
    Id        string
    UserId    string
    TaskId    string
    Points    int
    CreatedAt time.Time
}

func ToRewardRecord(rewardRecordEntity *RewardRecordEntity, taskEntity *TaskEntity) *trading.RewardRecord {
    return &trading.RewardRecord{Id: uuid.MustParse(rewardRecordEntity.Id), TaskName: taskEntity.Name, Points: rewardRecordEntity.Points, RewardTime: rewardRecordEntity.CreatedAt}
}

func ToRewardRecordEntity(user *trading.User, reward *trading.RewardRecord, tasksMap map[string]*TaskEntity) *RewardRecordEntity {
    now := time.Now()
    return &RewardRecordEntity{reward.Id.String(), user.Id.String(), tasksMap[reward.TaskName].Id, reward.Points, now}
}
