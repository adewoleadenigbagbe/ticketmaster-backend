module github.com/Wolechacho/ticketmaster-backend/infastructure

go 1.20

replace github.com/Wolechacho/ticketmaster-backend/shared => ../shared

require (
	github.com/SebastiaanKlippert/go-wkhtmltopdf v1.9.2
	github.com/Wolechacho/ticketmaster-backend/shared v1.0.0
	github.com/nats-io/nats.go v1.31.0
	github.com/samber/lo v1.39.0
	github.com/stripe/stripe-go v70.15.0+incompatible
	golang.org/x/crypto v0.23.0
	gorm.io/gorm v1.25.10
)

require (
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.17.2 // indirect
	github.com/labstack/echo/v4 v4.12.0 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/exp v0.0.0-20220303212507-bbda1eaf7a17 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
)
