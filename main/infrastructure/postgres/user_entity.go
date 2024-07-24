package postgres

import (
    "github.com/google/uuid"
    "math/big"
    "time"
    "tradingACE/main/trading"
)

type UserEntity struct {
    Id        string
    Address   string
    Amount    string
    Points    int
    CreatedAt time.Time
    UpdatedAt time.Time
}

func ToUser(userEntity *UserEntity) *trading.User {
    amount, _ := new(big.Int).SetString(userEntity.Amount, 10)
    return &trading.User{Id: uuid.MustParse(userEntity.Id), Address: userEntity.Address, TotalAmount: amount, TotalPoints: userEntity.Points}
}

func ToUserEntity(user *trading.User) *UserEntity {
    now := time.Now()
    return &UserEntity{Id: user.Id.String(), Address: user.Address, Amount: user.TotalAmount.String(), Points: user.TotalPoints, CreatedAt: now, UpdatedAt: now}
}
