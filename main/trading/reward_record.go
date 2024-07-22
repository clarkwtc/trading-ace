package trading

import "math/big"

type RewardRecord struct {
    TaskName string
    Point    int
    Amount   *big.Int
}
