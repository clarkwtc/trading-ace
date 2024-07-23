package repositories

import (
    "tradingACE/main/trading"
)

type UserRepository struct {
    access trading.UserRepository
}

func NewUserRepository(access trading.UserRepository) *UserRepository {
    return &UserRepository{access}
}

func (repository *UserRepository) FindUserTasksByAddress(address string) *trading.User {
    return repository.access.FindUserTasksByAddress(address)
}

func (repository *UserRepository) FindUserRewardByAddress(address string) *trading.User {
    return repository.access.FindUserRewardByAddress(address)
}
