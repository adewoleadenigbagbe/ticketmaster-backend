package enums

type ShowSeatStatus int

const (
	Available         ShowSeatStatus = 1
	Reserved          ShowSeatStatus = 2
	PendingAssignment ShowSeatStatus = 3
	Assigned          ShowSeatStatus = 4
)
