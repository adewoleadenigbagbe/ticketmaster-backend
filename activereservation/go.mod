module github.com/Wolechacho/ticketmaster-backend/activereservation

go 1.20

replace github.com/Wolechacho/ticketmaster-backend/shared => ../shared

require (
	github.com/Wolechacho/ticketmaster-backend/shared v1.0.0
	github.com/fatih/color v1.17.0
	github.com/muesli/cache2go v0.0.0-20221011235721-518229cd8021
	github.com/nats-io/nats.go v1.31.0
	github.com/samber/lo v1.39.0
	gorm.io/gorm v1.25.10
)

require (
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.17.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.23.0 // indirect
	golang.org/x/exp v0.0.0-20220303212507-bbda1eaf7a17 // indirect
	golang.org/x/sys v0.20.0 // indirect
	gorm.io/driver/mysql v1.5.6 // indirect
)
