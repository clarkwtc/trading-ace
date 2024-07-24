package eth

import (
    "context"
    "fmt"
    "github.com/ethereum/go-ethereum/ethclient"
    "tradingACE/main/infrastructure/server"
)

type EthClientManager struct {
    context context.Context
    Client  *ethclient.Client
    URL     string
}

func NewEthClientManager() *EthClientManager {
    ctx := context.Background()
    url := server.SystemConfig.Ethereum.URL
    client, err := ethclient.DialContext(ctx, url)
    if err != nil {
        fmt.Println(err)
    }
    return &EthClientManager{ctx, client, url}
}

func (manager *EthClientManager) ReconnectEthClient(url string) {
    manager.Close()

    ctx := context.Background()
    client, err := ethclient.DialContext(ctx, url)
    if err != nil {
        return
    }

    manager.context = ctx
    manager.Client = client
    manager.URL = url
}

func (manager *EthClientManager) Close() {
    manager.Client.Close()
}
