package application

import (
    "math/big"
    "tradingACE/main/trading"
)

type SettlementPointsUsecase struct {
    UserRepository trading.UserRepository
}

func (usecase *SettlementPointsUsecase) Execute(final bool) {
    users := usecase.UserRepository.FindAllUserTasks()
    if len(users) == 0 {
        return
    }

    campaign := trading.NewCampaign()
    campaign.Users = users
    campaign.Settlement(new(big.Int).SetUint64(100000))
    acceptNextTask(campaign, final)
    usecase.UserRepository.SaveAllUser(users)
}

func acceptNextTask(campaign *trading.Campaign, final bool) {
    if !final {
        for _, user := range campaign.Users {
            user.NextTask(trading.SharePoolTaskName, trading.SharePoolTaskName)
        }
    }
}
