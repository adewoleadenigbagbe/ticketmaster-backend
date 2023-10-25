package enums

type ShowSeatStatus int

const (
	Available   ShowSeatStatus = 1
	Reserved    ShowSeatStatus = 2
	PendingBook ShowSeatStatus = 3
	Booked      ShowSeatStatus = 4
	Cancelled   ShowSeatStatus = 5
)
