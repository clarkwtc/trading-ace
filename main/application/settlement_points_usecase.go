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
    users := usecase.UserRepository.FindAllUserTasks()
    if len(users) == 0 {
        return
    }

    campaign := trading.NewCampaign()
    campaign.Users = users
    campaign.Settlement(final)
    usecase.UserRepository.SaveAllUser(users)

    event := events.NewSettlementPointsEvent(campaign.Users)
    eventData, err := json.Marshal(event)
    if err != nil {
        log.Fatalf("Encode event fail: %v", err)
    }
    usecase.EventHub.Publish(eventData)
}
