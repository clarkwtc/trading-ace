package endpoints

import (
    "github.com/gin-gonic/gin"
    "tradingACE/main/trading/errors"
)

type ErrorMiddleware struct {
    errors.ErrorHandler
}

func (handler *ErrorMiddleware) RegisterException(exceptionHandler errors.IErrorHandler) {
    exceptionHandler.SetNext(handler.Concrete)
    handler.Concrete = exceptionHandler
}

func (handler *ErrorMiddleware) Execute() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        if len(ctx.Errors) > 0 {
            err := ctx.Errors.Last().Err
            handler.Handle(ctx, err)
            return
        }

        ctx.Next()

        if len(ctx.Errors) > 0 {
            err := ctx.Errors.Last().Err
            handler.Handle(ctx, err)
        }
    }
}
