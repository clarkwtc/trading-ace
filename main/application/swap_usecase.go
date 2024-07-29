package application

import (
    "encoding/json"
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
    campaign.AddUsers(user)
    campaign.Swap(address, amount)
    usecase.UserRepository.SaveAllUser(campaign.Users)

    event := events.NewSwapEvent(campaign.Users[0], amount)
    eventData, _ := json.Marshal(event)
    usecase.EventHub.Publish(eventData)
}
