package eventhub

import (
    "sync"
)

type EventHub struct {
    subscriptions map[*Client]bool
    mu            sync.RWMutex
}

func NewEventHub() *EventHub {
    return &EventHub{subscriptions: make(map[*Client]bool)}
}

func (eventHub *EventHub) Subscribe(client *Client) {
    eventHub.mu.Lock()
    defer eventHub.mu.Unlock()

    eventHub.subscriptions[client] = true
}

func (eventHub *EventHub) UnSubscribe(client *Client) {
    eventHub.mu.Lock()
    defer eventHub.mu.Unlock()

    delete(eventHub.subscriptions, client)
}

func (eventHub *EventHub) Publish(message []byte) {
    eventHub.mu.Lock()
    defer eventHub.mu.Unlock()
    for subscriber := range eventHub.subscriptions {
        select {
        case subscriber.Message <- message:
        default:
            close(subscriber.Message)
            delete(eventHub.subscriptions, subscriber)
        }
    }
}
