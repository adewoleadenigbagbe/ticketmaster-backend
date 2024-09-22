package enums

type PaymentType int

const (
	Stripe PaymentType = iota + 1
	PayPal
)
