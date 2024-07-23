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
    Mode           CampaignMode
    OnBoardingTask *OnBoardingTask
    SharePoolTask  *SharePoolTask
    Users          []*User
}

func NewCampaign() *Campaign {
    campaignMode := ParseCampaignMode(viper.GetString("campaign_mode"))
    return &Campaign{Mode: campaignMode, OnBoardingTask: NewOnBoardingTask(), SharePoolTask: NewSharePoolTask()}
}

func (campaign *Campaign) AddUsers(user *User) {
    campaign.Users = append(campaign.Users, user)
}

func (campaign *Campaign) Swap(address string, evnet *Event) {
    user := campaign.GetUserByAddress(address)
    if user == nil {
        user = NewUser(address)
    }
    campaign.Users = append(campaign.Users, user)

    user.AddAmount(campaign.OnBoardingTask.Name, evnet.Amount0Out)
    campaign.OnBoardingTask.Complete(user, evnet)
}

func (campaign *Campaign) Settlement(allUsersSwapAmount *big.Int) {
    for _, user := range campaign.Users {
        campaign.SharePoolTask.Complete(user, allUsersSwapAmount)
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
