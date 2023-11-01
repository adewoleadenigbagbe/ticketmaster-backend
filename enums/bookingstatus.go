package enums

type BookingStatus int

const (
	PendingBook BookingStatus = 1
	Booked      BookingStatus = 2
	Cancelled   BookingStatus = 3
)
