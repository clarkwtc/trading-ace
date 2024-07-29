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
        assert.Equal(t, trading.OnBoardingTaskName, campaign.OnBoardingTask.Name)
        assert.Equal(t, trading.SharePoolTaskName, campaign.SharePoolTask.Name)
        assert.Equal(t, 0, len(campaign.Users))
    }
}

func TestAddUsers(t *testing.T) {
    // Given
    viper.Set("campaign_mode", trading.ParseCampaignModeName(trading.PastBacktestMode))
    campaign := trading.NewCampaign()
    user := trading.NewUser(uuid.New().String())

    // When
    campaign.AddUsers(user)

    // Then
    assert.Equal(t, 1, len(campaign.Users))
    assert.Equal(t, user, campaign.Users[0])
}

func TestSwapOnNotCompletedTask(t *testing.T) {
    // Given
    viper.Set("campaign_mode", trading.ParseCampaignModeName(trading.PastBacktestMode))
    campaign := trading.NewCampaign()
    user := trading.NewUser(uuid.New().String())
    campaign.AddUsers(user)
    amount := utils.ToUSDC(new(big.Int).SetInt64(100))

    // When
    campaign.Swap(user.Address, amount)

    // Then
    assert.Equal(t, amount, user.TotalAmount)
    assert.Equal(t, 0, user.TotalPoints)
    assert.Equal(t, 1, len(user.Tasks))
    assert.Equal(t, trading.OnBoardingTaskName, user.Tasks[0].Name)
    assert.Equal(t, trading.OnGoing, user.Tasks[0].Status)
    assert.Equal(t, amount, user.Tasks[0].Amount)
    assert.Equal(t, 0, user.Tasks[0].Points)
}

func TestSwapOnCompletedTask(t *testing.T) {
    // Given
    viper.Set("campaign_mode", trading.ParseCampaignModeName(trading.PastBacktestMode))
    campaign := trading.NewCampaign()
    user := trading.NewUser(uuid.New().String())
    campaign.AddUsers(user)
    amount := utils.ToUSDC(new(big.Int).SetInt64(1001))

    // When
    campaign.Swap(user.Address, amount)

    // Then
    assert.Equal(t, amount, user.TotalAmount)
    assert.Equal(t, 100, user.TotalPoints)

    assert.Equal(t, 2, len(user.Tasks))
    assert.Equal(t, trading.OnBoardingTaskName, user.Tasks[0].Name)
    assert.Equal(t, trading.Completed, user.Tasks[0].Status)
    assert.Equal(t, amount, user.Tasks[0].Amount)
    assert.Equal(t, 100, user.Tasks[0].Points)

    assert.Equal(t, trading.SharePoolTaskName, user.Tasks[1].Name)
    assert.Equal(t, trading.OnGoing, user.Tasks[1].Status)
    assert.Equal(t, amount, user.Tasks[1].Amount)
    assert.Equal(t, 0, user.Tasks[1].Points)
}

func TestSettlement(t *testing.T) {
    // Given
    viper.Set("campaign_mode", trading.ParseCampaignModeName(trading.PastBacktestMode))
    campaign := trading.NewCampaign()
    user := trading.NewUser(uuid.New().String())
    campaign.AddUsers(user)
    amount := utils.ToUSDC(new(big.Int).SetInt64(2000))
    campaign.Swap(user.Address, amount)

    user2 := trading.NewUser(uuid.New().String())
    campaign.AddUsers(user2)
    amount2 := utils.ToUSDC(new(big.Int).SetInt64(2000))
    campaign.Swap(user2.Address, amount2)

    // When
    campaign.Settlement(false)

    // Then
    assert.Equal(t, amount, user.TotalAmount)
    assert.Equal(t, 5100, user.TotalPoints)
    assert.Equal(t, 5100, user2.TotalPoints)

    assert.Equal(t, 3, len(user.Tasks))
    assert.Equal(t, trading.OnBoardingTaskName, user.Tasks[0].Name)
    assert.Equal(t, trading.Completed, user.Tasks[0].Status)
    assert.Equal(t, amount, user.Tasks[0].Amount)
    assert.Equal(t, 100, user.Tasks[0].Points)

    assert.Equal(t, trading.SharePoolTaskName, user.Tasks[1].Name)
    assert.Equal(t, trading.Completed, user.Tasks[1].Status)
    assert.Equal(t, amount, user.Tasks[1].Amount)
    assert.Equal(t, 5000, user.Tasks[1].Points)

    assert.Equal(t, trading.SharePoolTaskName, user.Tasks[2].Name)
    assert.Equal(t, trading.OnGoing, user.Tasks[2].Status)
    assert.Equal(t, amount, user.Tasks[2].Amount)
    assert.Equal(t, 0, user.Tasks[2].Points)
}

func TestGetUserByAddress(t *testing.T) {
    // Given
    viper.Set("campaign_mode", trading.ParseCampaignModeName(trading.PastBacktestMode))
    campaign := trading.NewCampaign()
    newUser := trading.NewUser(uuid.New().String())
    campaign.AddUsers(newUser)

    // When
    user := campaign.GetUserByAddress(newUser.Address)

    // Then
    assert.Equal(t, newUser, user)
}
