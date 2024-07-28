package endpoints

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "tradingACE/main/application"
    "tradingACE/main/infrastructure/endpoints/models"
    "tradingACE/main/trading"
)

type UserResource struct {
    GetUserTasksStatusQuery   *application.GetUserTasksStatusQuery
    GetUserPointsHistoryQuery *application.GetUserPointsHistoryQuery
    GetLeaderboardQuery       *application.GetLeaderboardQuery
}

func NewUserResource(userRepository trading.UserRepository) *UserResource {
    return &UserResource{
        &application.GetUserTasksStatusQuery{UserRepository: userRepository},
        &application.GetUserPointsHistoryQuery{UserRepository: userRepository},
        &application.GetLeaderboardQuery{UserRepository: userRepository},
    }
}

type UserUri struct {
    Address string `uri:"address"`
}

func (resource *UserResource) GetUserTasksStatus(ctx *gin.Context) {
    var userUri UserUri

    err := ctx.ShouldBindUri(&userUri)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    event := resource.GetUserTasksStatusQuery.Execute(userUri.Address)
    ctx.JSON(http.StatusOK, models.ToTaskStatusViewModel(event))
}

func (resource *UserResource) GetUserPointsHistory(ctx *gin.Context) {
    var userUri UserUri

    err := ctx.ShouldBindUri(&userUri)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    event := resource.GetUserPointsHistoryQuery.Execute(userUri.Address)
    ctx.JSON(http.StatusOK, models.ToPointsHistoryViewModel(event))
}

type LeaderboardQuery struct {
    Task string `form:"task"`
}

func (resource *UserResource) GetLeaderboard(ctx *gin.Context) {
    var leaderboardQuery LeaderboardQuery
    err := ctx.ShouldBindQuery(&leaderboardQuery)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    event := resource.GetLeaderboardQuery.Execute(leaderboardQuery.Task)
    ctx.JSON(http.StatusOK, models.ToLeaderboardViewModel(event))
}
