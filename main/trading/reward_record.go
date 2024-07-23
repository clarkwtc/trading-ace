package trading

import (
    "github.com/google/uuid"
    "time"
)

type RewardRecord struct {
    Id         uuid.UUID
    TaskName   string
    Points     int
    RewardTime time.Time
}
