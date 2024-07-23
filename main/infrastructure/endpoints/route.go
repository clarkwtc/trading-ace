package endpoints

import (
    "database/sql"
    "github.com/gin-gonic/gin"
    "tradingACE/main/infrastructure/postgres"
    "tradingACE/main/infrastructure/repositories"
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
    }
}
