package application

import (
    "encoding/json"
    "log"
    "math/big"
    "tradingACE/main/infrastructure/eventhub"
    "tradingACE/main/trading"
    "tradingACE/main/trading/events"
)

type SwapUsecase struct {
    UserRepository trading.UserRepository
    EventHub       *eventhub.EventHub
}

func (usecase *SwapUsecase) Execute(address string, amount *big.Int) {
    user := usecase.UserRepository.FindUserTasksByAddress(address)
    if user == nil {
        user = trading.NewUser(address)
    }
    campaign := trading.NewCampaign()
    campaign.Swap(address, amount)
    usecase.UserRepository.SaveAllUser(campaign.Users)

    event := events.NewSwapEvent(campaign.Users[0], amount)
    eventData, err := json.Marshal(event)
    if err != nil {
        log.Fatalf("Encode event fail: %v", err)
    }
    usecase.EventHub.Publish(eventData)
}
