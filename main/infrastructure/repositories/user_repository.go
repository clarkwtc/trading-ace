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

func (repository *UserRepository) FindUserTasksByAddress(address string) (*trading.User, error) {
    return repository.access.FindUserTasksByAddress(address)
}

func (repository *UserRepository) FindAllUserTasks() ([]*trading.User, error) {
    return repository.access.FindAllUserTasks()
}

func (repository *UserRepository) SaveAllUser(users []*trading.User) error {
    return repository.access.SaveAllUser(users)
}
