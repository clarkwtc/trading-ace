package trading

type UserRepository interface {
    FindUserTasksByAddress(address string) *User
    FindAllUserTasks() []*User
    SaveAllUser(users []*User)
}
