package postgres

import (
    "time"
    "tradingACE/main/trading"
)

type RewardRecordEntity struct {
    Id        int
    UserId    int
    TaskId    int
    Points    int
    CreatedAt time.Time
}

func ToRewardRecord(rewardRecordEntity *RewardRecordEntity, taskEntity *TaskEntity) *trading.RewardRecord {
    return &trading.RewardRecord{TaskName: taskEntity.Name, Points: rewardRecordEntity.Points, RewardTime: rewardRecordEntity.CreatedAt}
}
