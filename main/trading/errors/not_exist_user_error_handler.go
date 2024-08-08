package errors

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

type NotExistUserErrorHandler struct {
    StatusCode int
    Error      error
    ErrorHandler
}

func NewNotExistUserErrorHandler() *NotExistUserErrorHandler {
    return &NotExistUserErrorHandler{StatusCode: http.StatusNotFound, Error: &NotExistUserError{}}
}

func (handler *NotExistUserErrorHandler) Match(error error) bool {
    switch error.(type) {
    case *NotExistUserError:
        return true
    default:
        return false
    }
}

func (handler *NotExistUserErrorHandler) Response(ctx *gin.Context) {
    ctx.JSON(handler.StatusCode, gin.H{"error": NotExistUser})
}
