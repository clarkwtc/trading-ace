package trading

type UserRepository interface {
    FindUserTasksByAddress(address string) (*User, error)
    FindAllUserTasks() ([]*User, error)
    SaveAllUser(users []*User) error
}
