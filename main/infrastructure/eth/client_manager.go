package eth

import (
    "context"
    "github.com/ethereum/go-ethereum/ethclient"
    "log"
    "tradingACE/main/infrastructure/server"
)

type ClientManager struct {
    context context.Context
    Client  *ethclient.Client
    URL     string
}

func NewClientManager() *ClientManager {
    ctx := context.Background()
    url := server.SystemConfig.Ethereum.URL
    client, err := ethclient.DialContext(ctx, url)
    if err != nil {
        log.Fatalf("Could not connect eth client: %v", err)
    }
    return &ClientManager{ctx, client, url}
}

func (manager *ClientManager) ReconnectEthClient(url string) {
    manager.Close()

    ctx := context.Background()
    client, err := ethclient.DialContext(ctx, url)
    if err != nil {
        log.Fatalf("Could not reconnect eth client: %v", err)
    }

    manager.context = ctx
    manager.Client = client
    manager.URL = url
}

func (manager *ClientManager) Close() {
    manager.Client.Close()
}
