package enums

type BookingStatus int

const (
	BookNone    BookingStatus = 0
	Available   BookingStatus = 1
	Reserved    BookingStatus = 2
	PendingBook BookingStatus = 3
	Booked      BookingStatus = 4
	Cancelled   BookingStatus = 5
)
