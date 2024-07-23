package trading


type UserRepository interface {
    FindUserTasksByAddress(address string) *User
}
