package endpoints

import (
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    "log"
    "net/http"
    "strings"
    "tradingACE/main/infrastructure/eventhub"
)

type WebsocketConnection struct {
    EventHub *eventhub.EventHub
}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func (socket *WebsocketConnection) BindConnection(ctx *gin.Context) {
    connection, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
    if err != nil {
        ctx.JSON(http.StatusHTTPVersionNotSupported, gin.H{"error": err.Error()})
        return
    }

    client := eventhub.NewClient(connection)
    socket.EventHub.Subscribe(client)
    go socket.handleReadMessage(client)
    defer socket.close(client)
    socket.handleWriteMessage(client)
}

func (socket *WebsocketConnection) close(client *eventhub.Client) {
    socket.EventHub.UnSubscribe(client)
    err := client.Connection.Close()
    if err != nil {
        log.Printf("Could not to close websocket connection: %v", err)
        return
    }
}

func (socket *WebsocketConnection) handleReadMessage(client *eventhub.Client) {
    for {
        _, message, err := client.Connection.ReadMessage()
        if err != nil {
            log.Printf("Read message fail: %v", err)
            break
        }

        switch strings.ToLower(string(message)) {
        case "unsubscribe":
            socket.close(client)
            break
        default:
            socket.EventHub.Publish(message)
        }
    }
}

func (socket *WebsocketConnection) handleWriteMessage(client *eventhub.Client) {
    for {
        select {
        case message := <-client.Message:
            err := client.Connection.WriteMessage(websocket.TextMessage, message)
            if err != nil {
                log.Println("Write message fail:", err)
                break
            }
        }
    }
}
