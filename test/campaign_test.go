package test

import (
    "github.com/google/uuid"
    "github.com/spf13/viper"
    "github.com/stretchr/testify/assert"
    "math/big"
    "testing"
    "tradingACE/main/infrastructure/utils"
    "tradingACE/main/trading"
)

func TestParseCampaignModeName(t *testing.T) {
    tests := []struct {
        campaignMode     trading.CampaignMode
        campaignModeName string
    }{
        {trading.PastBacktestMode, "PastBacktestMode"},
        {trading.CurrentActiveMode, "CurrentActiveMode"},
    }

    for _, task := range tests {
        // Given
        campaignMode := task.campaignMode

        // When
        campaignModeName := trading.ParseCampaignModeName(campaignMode)

        // Then
        assert.Equal(t, task.campaignModeName, campaignModeName)
    }
}

func TestParseCampaignMode(t *testing.T) {
    tests := []struct {
        campaignMode     trading.CampaignMode
        campaignModeName string
    }{
        {trading.PastBacktestMode, "PastBacktestMode"},
        {trading.CurrentActiveMode, "CurrentActiveMode"},
    }

    for _, task := range tests {
        // Given
        campaignModeName := task.campaignModeName
        // When
        campaignMode := trading.ParseCampaignMode(campaignModeName)

        // Then
        assert.Equal(t, task.campaignMode, campaignMode)
    }
}

func TestNewCampaignOnPastBacktestMode(t *testing.T) {
    tests := []struct {
        campaignMode trading.CampaignMode
    }{
        {trading.PastBacktestMode},
        {trading.CurrentActiveMode},
    }

    for _, task := range tests {
        // Given
        viper.Set("campaign_mode", trading.ParseCampaignModeName(task.campaignMode))
        // When
        campaign := trading.NewCampaign()

        // Then
        assert.Equal(t, task.campaignMode, campaign.Mode)
        assert.Equal(t, 0, len(campaign.Users))
    }
}

func TestSwapOnNotCompletedTask(t *testing.T) {
    // Given
    viper.Set("campaign_mode", trading.ParseCampaignModeName(trading.PastBacktestMode))
    campaign := trading.NewCampaign()
    address := uuid.New().String()
    amount := utils.ToUSDC(new(big.Int).SetInt64(100))

    // When
    campaign.Swap(address, amount)

    // Then
    user := campaign.GetUserByAddress(address)
    assert.Equal(t, amount, user.TotalAmount)
    assert.Equal(t, 0, user.TotalPoints)
    assert.Equal(t, 1, user.CountTaskRecord())

    taskRecord := user.GetTaskRecords()[0]
    assert.Equal(t, trading.OnBoardingTaskName, taskRecord.Task.GetName())
    assert.Equal(t, trading.OnGoing, taskRecord.Status)
    assert.Equal(t, amount, taskRecord.SwapAmount)
    assert.Equal(t, 0, taskRecord.EarnPoints)
}

func TestSwapOnCompletedTask(t *testing.T) {
    // Given
    viper.Set("campaign_mode", trading.ParseCampaignModeName(trading.PastBacktestMode))
    campaign := trading.NewCampaign()
    address := uuid.New().String()
    amount := utils.ToUSDC(new(big.Int).SetInt64(1001))

    // When
    campaign.Swap(address, amount)

    // Then
    user := campaign.GetUserByAddress(address)
    assert.Equal(t, amount, user.TotalAmount)
    assert.Equal(t, 100, user.TotalPoints)
    assert.Equal(t, 2, user.CountTaskRecord())

    taskRecord := user.GetTaskRecords()[0]
    assert.Equal(t, trading.OnBoardingTaskName, taskRecord.Task.GetName())
    assert.Equal(t, trading.Completed, taskRecord.Status)
    assert.Equal(t, amount, taskRecord.SwapAmount)
    assert.Equal(t, 100, taskRecord.EarnPoints)

    taskRecord = user.GetTaskRecords()[1]
    assert.Equal(t, trading.SharePoolTaskName, taskRecord.Task.GetName())
    assert.Equal(t, trading.OnGoing, taskRecord.Status)
    assert.Equal(t, amount, taskRecord.SwapAmount)
    assert.Equal(t, 0, taskRecord.EarnPoints)
}

func TestSettlement(t *testing.T) {
    // Given
    viper.Set("campaign_mode", trading.ParseCampaignModeName(trading.PastBacktestMode))
    campaign := trading.NewCampaign()
    address := uuid.New().String()
    amount := utils.ToUSDC(new(big.Int).SetInt64(2000))
    campaign.Swap(address, amount)

    address2 := uuid.New().String()
    amount2 := utils.ToUSDC(new(big.Int).SetInt64(2000))
    campaign.Swap(address2, amount2)

    // When
    campaign.Settlement(false)

    // Then
    user := campaign.GetUserByAddress(address)
    user2 := campaign.GetUserByAddress(address2)
    assert.Equal(t, 3, user.CountTaskRecord())
    assert.Equal(t, amount, user.TotalAmount)
    assert.Equal(t, 5100, user.TotalPoints)
    assert.Equal(t, 5100, user2.TotalPoints)

    taskRecord := user.GetTaskRecords()[0]
    assert.Equal(t, trading.OnBoardingTaskName, taskRecord.Task.GetName())
    assert.Equal(t, trading.Completed, taskRecord.Status)
    assert.Equal(t, amount, taskRecord.SwapAmount)
    assert.Equal(t, 100, taskRecord.EarnPoints)

    taskRecord = user.GetTaskRecords()[1]
    assert.Equal(t, trading.SharePoolTaskName, taskRecord.Task.GetName())
    assert.Equal(t, trading.Completed, taskRecord.Status)
    assert.Equal(t, amount, taskRecord.SwapAmount)
    assert.Equal(t, 5000, taskRecord.EarnPoints)

    taskRecord = user.GetTaskRecords()[2]
    assert.Equal(t, trading.SharePoolTaskName, taskRecord.Task.GetName())
    assert.Equal(t, trading.OnGoing, taskRecord.Status)
    assert.Equal(t, new(big.Int).SetInt64(0), taskRecord.SwapAmount)
    assert.Equal(t, 0, taskRecord.EarnPoints)
}
