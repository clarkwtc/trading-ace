package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "log"
    "tradingACE/main/infrastructure/endpoints"
    "tradingACE/main/infrastructure/eventhub"
    "tradingACE/main/infrastructure/postgres"
    "tradingACE/main/infrastructure/server"
)

func main() {
    server.InitConfig()
    client := postgres.Init()
    eventHub := eventhub.NewEventHub()
    router := endpoints.Router{Engine: gin.Default(), PostgreSQLClient: client, EventHub: eventHub}
    router.SetupErrorMiddleware()
    router.SetupUserResource()
    router.SetupWebsocketConnection()
    go router.SetupCampaignTimer()
    err := router.Engine.Run(fmt.Sprintf(":%d", server.SystemConfig.Server.Port))
    if err != nil {
        log.Fatalf("Could not to run web engine: %v", err)
        return
    }
}
