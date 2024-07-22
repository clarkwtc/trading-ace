package server

import (
    "github.com/spf13/viper"
    "log"
    "time"
    "tradingACE/main/trading"
)

type Config struct {
    Database struct {
        Host     string
        User     string
        Password string
        Name     string
        Port     int
        SSLMode  string
        TimeZone string
    }
    Server struct {
        Port int
    }

    Campaign struct {
        StartTime string
    }
}

var SystemConfig Config

func InitConfig() {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")

    err := viper.ReadInConfig()
    if err != nil {
        log.Fatal("Could not read config", err)
    }

    err = viper.Unmarshal(&SystemConfig)
    if err != nil {
        log.Fatal("Unable to decode config", err)
    }

    setCampaignMode()
}

func setCampaignMode() {
    campaignMode := trading.PastBacktestMode
    layout := "2006-01-02T15:00:00"
    startTime, err := time.Parse(layout, SystemConfig.Campaign.StartTime)
    if err != nil {
        log.Fatalf("Error parsing start time: %v", err)
    }

    now := time.Now()
    if now.After(startTime) {
        campaignMode = trading.CurrentActiveMode
    }

    viper.Set("campaign_mode", trading.ParseCampaignModeName(campaignMode))
}
