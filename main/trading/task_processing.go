package trading

type TaskProcessingStatus int

const (
    OnGoing TaskProcessingStatus = iota
    Completed
)

func IsTaskProcessingStatus(status TaskProcessingStatus) bool {
    switch status {
    case OnGoing, Completed:
        return true
    default:
        return false
    }
}

type TaskProcessing struct {
    TaskName string
    Status   TaskProcessingStatus
}
