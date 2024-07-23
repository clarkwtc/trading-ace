package endpoints

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "tradingACE/main/application"
    "tradingACE/main/infrastructure/endpoints/models"
    "tradingACE/main/trading"
)

type UserResource struct {
    GetUserTasksStatusQuery *application.GetUserTasksStatusQuery
}

func NewUserResource(userRepository trading.UserRepository) *UserResource {
    return &UserResource{
        &application.GetUserTasksStatusQuery{UserRepository: userRepository},
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
