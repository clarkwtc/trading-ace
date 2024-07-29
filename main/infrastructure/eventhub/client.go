package eventhub

import "github.com/gorilla/websocket"

type Client struct {
    Connection *websocket.Conn
    Message    chan []byte
}

func NewClient(connection *websocket.Conn) *Client {
    return &Client{connection, make(chan []byte)}
}
