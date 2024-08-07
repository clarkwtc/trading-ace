package trading

import (
    "github.com/spf13/viper"
    "log"
    "math/big"
    "strings"
)

type CampaignMode int

const (
    PastBacktestMode CampaignMode = iota
    CurrentActiveMode
)

func ParseCampaignModeName(mode CampaignMode) string {
    switch mode {
    case PastBacktestMode:
        return "PastBacktestMode"
    case CurrentActiveMode:
        return "CurrentActiveMode"
    default:
        log.Fatalf("Not support mode %v", mode)
    }
    return ""
}

func ParseCampaignMode(modeName string) CampaignMode {
    modeName = strings.ToLower(modeName)
    switch modeName {
    case strings.ToLower("PastBacktestMode"):
        return PastBacktestMode
    case strings.ToLower("CurrentActiveMode"):
        return CurrentActiveMode
    default:
        log.Fatalf("Not support mode %v", modeName)
    }
    return -1
}

type Campaign struct {
    Mode  CampaignMode
    Users []*User
}

func NewCampaign() *Campaign {
    campaignMode := ParseCampaignMode(viper.GetString("campaign_mode"))
    return &Campaign{Mode: campaignMode}
}

func (campaign *Campaign) AddUsers(user *User) {
    if user == nil {
        return
    }
    campaign.Users = append(campaign.Users, user)
}

func (campaign *Campaign) Swap(address string, amount *big.Int) {
    user := campaign.GetUserByAddress(address)
    if user == nil {
        user = NewUser(address)
    }

    campaign.AddUsers(user)
    user.AddAmount(amount)
    taskRecord := user.GetTaskRecordByStatus(OnGoing)
    taskRecord.AddSwapAmount(amount)

    switch task := taskRecord.Task.(type) {
    case *OnBoardingTask:
        task.Complete(amount)
    }
}

func (campaign *Campaign) GetUserByAddress(address string) *User {
    for _, user := range campaign.Users {
        if user.Address == address {
            return user
        }
    }
    return nil
}

func (campaign *Campaign) Settlement(final bool) {
    totalAmount := campaign.sumAllUserAmount()

    for _, user := range campaign.Users {
        taskRecord := user.GetTaskRecordByStatus(OnGoing)

        switch task := taskRecord.Task.(type) {
        case *SharePoolTask:
            task.Complete(totalAmount, final)
        }
    }
}

func (campaign *Campaign) sumAllUserAmount() *big.Int {
    sum := new(big.Int).SetInt64(0)
    for _, user := range campaign.Users {
        sum.Add(sum, user.TotalAmount)
    }
    return sum
}
