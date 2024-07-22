package trading

import (
    "math/big"
    "time"
)

type CampaignMode int

const (
    PastBacktestMode CampaignMode = iota
    CurrentActiveMode
)

func IsCampaignMode(mode CampaignMode) bool {
    switch mode {
    case PastBacktestMode, CurrentActiveMode:
        return true
    default:
        return false
    }
}

type Campaign struct {
    StartTime      time.Time
    Period         time.Time
    Mode           CampaignMode
    OnBoardingTask *OnBoardingTask
    SharePoolTask  *SharePoolTask
    Users          []*User
}

const week = 7

func NewCampaign(startTime time.Time) *Campaign {
    now := time.Now()
    campaignMode := PastBacktestMode
    if now.After(startTime) {
        campaignMode = CurrentActiveMode
    }

    period := now.AddDate(0, 0, week*4)
    return &Campaign{StartTime: startTime, Period: period, Mode: campaignMode, OnBoardingTask: NewOnBoardingTask(), SharePoolTask: NewSharePoolTask()}
}

func (campaign *Campaign) Swap(address string, evnet *Event) {
    user := campaign.GetUserByAddress(address)
    if user == nil {
        user = NewUser(address, &TaskProcessing{campaign.OnBoardingTask.Name, OnGoing})
    }
    campaign.Users = append(campaign.Users, user)

    user.AddAmount(evnet.Amount0Out)
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
