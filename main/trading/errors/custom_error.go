package errors

type ICustomError interface {
    Error() string
    error() string
}
type CustomError struct {
    Message string
}

func (error *CustomError) Error() string {
    return error.Message
}

func (error *CustomError) error() string {
    return error.Message
}
