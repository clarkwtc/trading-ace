package endpoints

import (
    "database/sql"
    "github.com/gin-gonic/gin"
    "github.com/spf13/viper"
    "log"
    "time"
    "tradingACE/main/infrastructure/eventhub"
    "tradingACE/main/infrastructure/postgres"
    "tradingACE/main/infrastructure/repositories"
    "tradingACE/main/infrastructure/server"
    "tradingACE/main/trading"
    "tradingACE/main/trading/errors"
)

type Router struct {
    Engine           *gin.Engine
    PostgreSQLClient *sql.DB
    EventHub         *eventhub.EventHub
}

func (router *Router) SetupUserResource() {
    userEndpoints := NewUserResource(repositories.NewUserRepository(postgres.NewUserRepository(router.PostgreSQLClient)))

    userRoutes := router.Engine.Group("/users")
    {
        userRoutes.GET("/:address/getTaskStatus", userEndpoints.GetUserTasksStatus)
        userRoutes.GET("/:address/getPointsHistory", userEndpoints.GetUserPointsHistory)
        userRoutes.GET("/:address/getLeaderboard", userEndpoints.GetLeaderboard)
    }
}

func (router *Router) SetupErrorMiddleware() {
    errorMiddleware := ErrorMiddleware{}
    errorMiddleware.RegisterException(errors.NewNotExistUserErrorHandler())
    router.Engine.Use(errorMiddleware.Execute())
}

func (router *Router) SetupWebsocketConnection() {
    websocketEndpoints := &WebsocketConnection{router.EventHub}

    wsRoutes := router.Engine.Group("/ws")
    {
        wsRoutes.GET("/bindConection", websocketEndpoints.BindConnection)
    }
}

func (router *Router) SetupCampaignTimer() {
    log.Println(viper.Get("campaign_mode"))
    if trading.ParseCampaignMode(viper.GetString("campaign_mode")) != trading.CurrentActiveMode {
        return
    }

    layout := "2006-01-02T15:04:05-07:00"
    startTime, err := time.Parse(layout, server.SystemConfig.Campaign.StartTime)
    if err != nil {
        log.Fatalf("Error parsing start time: %v", err)
    }

    repository := repositories.NewUserRepository(postgres.NewUserRepository(router.PostgreSQLClient))
    campaignEndpoints := NewCampaignResource(repository, router.EventHub)

    now := time.Now()
    duration := startTime.Sub(now)
    time.AfterFunc(duration, func() {
        go campaignEndpoints.SettlementPoints()
    })
    time.AfterFunc(duration, func() {
        go campaignEndpoints.WatchSwapEvents()
    })
}
