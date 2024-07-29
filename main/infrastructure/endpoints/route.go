package endpoints

import (
    "database/sql"
    "fmt"
    "github.com/ethereum/go-ethereum/common"
    "github.com/gin-gonic/gin"
    "github.com/spf13/viper"
    "log"
    "math/big"
    "time"
    "tradingACE/main/application"
    "tradingACE/main/infrastructure/contract"
    "tradingACE/main/infrastructure/eth"
    "tradingACE/main/infrastructure/postgres"
    "tradingACE/main/infrastructure/repositories"
    "tradingACE/main/infrastructure/server"
    "tradingACE/main/infrastructure/utils"
    "tradingACE/main/trading"
)

type Router struct {
    Engine           *gin.Engine
    PostgreSQLClient *sql.DB
}

func (router *Router) SetupUserResource() {
    userEndpoints := NewUserResource(repositories.NewUserRepository(postgres.NewUserRepository(router.PostgreSQLClient)))

    userRoutes := router.Engine.Group("/users")
    {
        userRoutes.GET("/:address/getTaskStatus", userEndpoints.GetUserTasksStatus)
        userRoutes.GET("/:address/getPointsHistory", userEndpoints.GetUserPointsHistory)
        userRoutes.GET("/:address/getLeaderboard", userEndpoints.GetLeaderboard)
    }
}

func (router *Router) StartCampaign() {
    fmt.Println(viper.Get("campaign_mode"))
    if trading.ParseCampaignMode(viper.GetString("campaign_mode")) != trading.CurrentActiveMode {
        return
    }

    layout := "2006-01-02T15:04:05-07:00"
    startTime, err := time.Parse(layout, server.SystemConfig.Campaign.StartTime)
    if err != nil {
        log.Fatalf("Error parsing start time: %v", err)
    }

    now := time.Now()

    duration := startTime.Sub(now)
    time.AfterFunc(duration, func() {
        go router.settlementPoints()
    })
    time.AfterFunc(duration, func() {
        go router.WatchSwapEvents()
    })
}

func (router *Router) settlementPoints() {
    ticker := time.NewTicker(24 * time.Hour * 7)
    defer ticker.Stop()

    repository := repositories.NewUserRepository(postgres.NewUserRepository(router.PostgreSQLClient))
    query := &application.SettlementPointsUsecase{UserRepository: repository}

    final := false
    for week := 1; week <= 4; week++ {
        <-ticker.C

        if week == 4 {
            final = true
        }
        query.Execute(final)
    }
}

func (router *Router) WatchSwapEvents() {
    clientManageer := eth.NewEthClientManager()
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

            router.swap(event)
        }
    }
}

func (router *Router) swap(event *contract.UniswapV2Swap) {
    repository := repositories.NewUserRepository(postgres.NewUserRepository(router.PostgreSQLClient))
    query := &application.SwapUsecase{UserRepository: repository}

    amount := new(big.Int).SetInt64(0)
    if event.Amount0In.Cmp(amount) == 1 {
        amount = event.Amount0In
    } else {
        amount = event.Amount0Out
    }
    query.Execute(event.Sender.String(), amount)
}
