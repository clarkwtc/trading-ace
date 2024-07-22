package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "tradingACE/main/infrastructure/endpoints"
    "tradingACE/main/infrastructure/postgres"
    "tradingACE/main/infrastructure/server"
)

func main() {
    server.InitConfig()
    client := postgres.Init()
    router := endpoints.Router{Engine: gin.Default(), PostgreSQLClient: client}
    router.SetupErrorMiddleware()
    router.SetupUserResource()
    err := router.Engine.Run(fmt.Sprintf(":%d", server.SystemConfig.Server.Port))
    if err != nil {
        return
    }
}
