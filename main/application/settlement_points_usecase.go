package application

import (
    "encoding/json"
    "log"
    "tradingACE/main/infrastructure/eventhub"
    "tradingACE/main/trading"
    "tradingACE/main/trading/events"
)

type SettlementPointsUsecase struct {
    UserRepository trading.UserRepository
    EventHub       *eventhub.EventHub
}

func (usecase *SettlementPointsUsecase) Execute(final bool) {
    users, err := usecase.UserRepository.FindAllUserTasks()
    if err != nil {
        log.Printf("FindAllUserTasks fail: %v", err)
        return
    }
    
    if len(users) == 0 {
        return
    }

    campaign := trading.NewCampaign()
    campaign.Users = users
    campaign.Settlement(final)
    err = usecase.UserRepository.SaveAllUser(users)
    if err != nil {
        log.Printf("SaveAllUser fail: %v", err)
        return
    }

    event := events.NewSettlementPointsEvent(campaign.Users)
    eventData, err := json.Marshal(event)
    if err != nil {
        log.Printf("Encode event fail: %v", err)
        return
    }
    usecase.EventHub.Publish(eventData)
}
