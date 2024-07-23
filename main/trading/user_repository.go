package trading

type UserRepository interface {
    FindUserTasksByAddress(address string) *User
    FindUserRewardByAddress(address string) *User
}
