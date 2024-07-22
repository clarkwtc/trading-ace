package main

import (
    "fmt"
    "github.com/google/uuid"
    "math/big"
    "time"
    "tradingACE/main/trading"
)

func main() {
    campaign := trading.NewCampaign(time.Now())

    uid := uuid.New().String()
    amount0 := new(big.Int).SetInt64(2000)
    amount1 := new(big.Int).SetInt64(900)
    campaign.Swap(uid, &trading.Event{Amount0Out: amount0, Amount1Out: amount1, To: uid})
    user := campaign.GetUserByAddress(uid)
    campaign.Settlement(new(big.Int).SetInt64(2000))

    fmt.Println(user)
}
