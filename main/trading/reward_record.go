package trading

import "time"

type RewardRecord struct {
    TaskName   string
    Points     int
    RewardTime time.Time
}
