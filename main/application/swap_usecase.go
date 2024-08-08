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
    user, err := usecase.UserRepository.FindUserTasksByAddress(address)
    if err != nil {
        log.Printf("FindUserTasksByAddress fail: %v", err)
        return
    }
    
    if user == nil {
        return
    }

    campaign := trading.NewCampaign()
    campaign.AddUsers(user)
    campaign.Swap(address, amount)
    err = usecase.UserRepository.SaveAllUser(campaign.Users)
    if err != nil {
        log.Printf("SaveAllUser fail: %v", err)
        return
    }

    event := events.NewSwapEvent(campaign.Users[0], amount)
    eventData, err := json.Marshal(event)
    if err != nil {
        log.Printf("Encode event fail: %v", err)
        return
    }
    usecase.EventHub.Publish(eventData)
}
