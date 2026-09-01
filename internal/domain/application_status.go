package domain

type ApplicationStatus string

const (
	ApplicationStatusNew        ApplicationStatus = "new"
	ApplicationStatusInProgress ApplicationStatus = "in_progress"
	ApplicationStatusSuccess    ApplicationStatus = "success"
	ApplicationStatusRejected   ApplicationStatus = "rejected"
)

func (s ApplicationStatus) IsValid() bool {
	switch s {
	case ApplicationStatusNew,
		ApplicationStatusInProgress,
		ApplicationStatusSuccess,
		ApplicationStatusRejected:
		return true
	default:
		return false
	}
}
