package endpoints

import (
    "database/sql"
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/spf13/viper"
    "log"
    "time"
    "tradingACE/main/application"
    "tradingACE/main/infrastructure/postgres"
    "tradingACE/main/infrastructure/repositories"
    "tradingACE/main/infrastructure/server"
    "tradingACE/main/trading"
)

type Router struct {
    Engine           *gin.Engine
    PostgreSQLClient *sql.DB
}

func (router *Router) SetupUserResource() {
    userEndpoints := NewUserResource(repositories.NewUserRepository(postgres.NewUserRepository(router.PostgreSQLClient)))

    userRoutes := router.Engine.Group("/users")
    {
        userRoutes.GET("/:address/getTaskStatus", userEndpoints.GetUserTasksStatus)
        userRoutes.GET("/:address/getPointsHistory", userEndpoints.GetUserPointsHistory)
    }
}

func (router *Router) SetupErrorMiddleware() {
    errorMiddleware := ErrorMiddleware{}
    router.Engine.Use(errorMiddleware.Execute())
}

func (router *Router) StartCampaign() {
    fmt.Println(viper.Get("campaign_mode"))
    if trading.ParseCampaignMode(viper.GetString("campaign_mode")) != trading.CurrentActiveMode {
        return
    }

    layout := "2006-01-02T15:04:05-07:00"
    startTime, err := time.Parse(layout, server.SystemConfig.Campaign.StartTime)
    if err != nil {
        log.Fatalf("Error parsing start time: %v", err)
    }

    now := time.Now()

    duration := startTime.Sub(now)
    time.AfterFunc(duration, router.settlementPoints())
}

func (router *Router) settlementPoints() func() {
    return func() {
        ticker := time.NewTicker(24 * time.Hour * 7)
        defer ticker.Stop()

        repository := repositories.NewUserRepository(postgres.NewUserRepository(router.PostgreSQLClient))
        query := &application.SettlementPointsUsecase{UserRepository: repository}

        final := false
        for week := 1; week <= 4; week++ {
            <-ticker.C

            if week == 4 {
                final = true
            }
            query.Execute(final)
        }
    }
}
