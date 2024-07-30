package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
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
    router.SetupUserResource()
    router.SetupWebsocketConnection()
    go router.SetupCampaignTimer()
    err := router.Engine.Run(fmt.Sprintf(":%d", server.SystemConfig.Server.Port))
    if err != nil {
        return
    }
}
