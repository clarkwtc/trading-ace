package endpoints

import (
    "fmt"
    "github.com/ethereum/go-ethereum/common"
    "log"
    "math/big"
    "time"
    "tradingACE/main/application"
    "tradingACE/main/infrastructure/contract"
    "tradingACE/main/infrastructure/eth"
    "tradingACE/main/infrastructure/eventhub"
    "tradingACE/main/infrastructure/utils"
    "tradingACE/main/trading"
)

type CampaignResource struct {
    SettlementPointsUsecase *application.SettlementPointsUsecase
    SwapUsecase             *application.SwapUsecase
}

func NewCampaignResource(repository trading.UserRepository, eventHub *eventhub.EventHub) *CampaignResource {
    return &CampaignResource{
        &application.SettlementPointsUsecase{UserRepository: repository, EventHub: eventHub},
        &application.SwapUsecase{UserRepository: repository, EventHub: eventHub},
    }
}

func (resource *CampaignResource) SettlementPoints() {
    ticker := time.NewTicker(24 * time.Hour * 7)
    defer ticker.Stop()

    final := false
    for week := 1; week <= 4; week++ {
        <-ticker.C

        if week == 4 {
            final = true
        }
        resource.SettlementPointsUsecase.Execute(final)
    }
}

func (resource *CampaignResource) WatchSwapEvents() {
    clientManageer := eth.NewClientManager()
    client := clientManageer.Client

    // Uniswap V2 Pair contract address for WETH/USDC
    pairAddress := common.HexToAddress("0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc")

    filterer, err := contract.NewUniswapV2Filterer(pairAddress, client)
    if err != nil {
        log.Fatalf("Failed to create UniswapV2Filterer:: %v", err)
    }

    eventChan := make(chan *contract.UniswapV2Swap)

    sub, err := filterer.WatchSwap(nil, eventChan, nil, nil)
    if err != nil {
        log.Fatalf("Failed to subscribe swap: %v", err)
    }

    fmt.Println("Listening swap events...")

    for {
        select {
        case err := <-sub.Err():
            log.Fatalf("Failed to receive events: %v", err)
        case event := <-eventChan:
            fmt.Printf("Swap Event:\n")
            fmt.Printf("Sender: %s\n", event.Sender.Hex())
            fmt.Printf("USDCAmountIn: %s\n", utils.ForamtUSDC(event.Amount0In))
            fmt.Printf("WETHAmountIn: %s\n", utils.ForamtEther(event.Amount1In))
            fmt.Printf("USDCAmountOut: %s\n", utils.ForamtUSDC(event.Amount0Out))
            fmt.Printf("WETHAmountOut: %s\n", utils.ForamtEther(event.Amount1Out))
            fmt.Printf("To: %s\n", event.To.Hex())

            resource.SwapUsecase.Execute(event.Sender.String(), getAmountVol(event))
        }
    }
}

func getAmountVol(event *contract.UniswapV2Swap) *big.Int {
    amount := new(big.Int).SetInt64(0)
    if event.Amount0In.Cmp(amount) == 1 {
        amount = event.Amount0In
    } else {
        amount = event.Amount0Out
    }
    return amount
}
