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
    router := endpoints.Router{Engine: gin.Default(), PostgreSQLClient: client, EventHub: eventhub.NewEventHub()}
    router.SetupUserResource()
    go router.StartCampaign()
    err := router.Engine.Run(fmt.Sprintf(":%d", server.SystemConfig.Server.Port))
    if err != nil {
        return
    }
}
