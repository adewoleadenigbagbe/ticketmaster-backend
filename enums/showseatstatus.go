package enums

type ShowSeatStatus int

const (
	Available ShowSeatStatus = 1
	Reserved  ShowSeatStatus = 2
	Booked    ShowSeatStatus = 3
	Cancelled ShowSeatStatus = 4
)
