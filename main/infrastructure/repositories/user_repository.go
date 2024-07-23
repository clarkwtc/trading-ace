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

//func (repository *UserRepository) FindByIds(ids []uuid.UUID) []*domain.User {
//    return repository.access.FindByIds(ids)
//}
//
//func (repository *UserRepository) FindByUsername(username string) *domain.User {
//    return repository.access.FindByUsername(username)
//}
//
//func (repository *UserRepository) Save(user *domain.User) {
//    repository.access.Save(user)
//}
