package application

import (
    "math/big"
    "tradingACE/main/trading"
)

type SwapUsecase struct {
    UserRepository trading.UserRepository
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
}
