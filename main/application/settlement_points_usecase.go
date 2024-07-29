package application

import (
    "encoding/json"
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
    eventData, _ := json.Marshal(event)
    usecase.EventHub.Publish(eventData)
}
