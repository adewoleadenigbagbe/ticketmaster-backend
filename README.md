# Ticketmaster Backend Service

A ticketmaster service is a movie booking system .The service handles ticketing logistics, including seating arrangements, pricing, and customer service. That it !!!.

## Components
The project implements a microservice architecture and is splited in three service which are 

* Rest Api Service
* Reservation Service
* Waiting Service

The Rest api service exposes all endpoints for app functionality (Ticket Booking, Movie Search Payment Charge) as well as other admin operations
The Reservation Service tracks of all active reservations and remove any expired reservation from the system
The Waiting Service tracks all waiting user requests and, as soon as the required number of seats become available (Server Sent Events)

## Implementation Story
Basically the user books seat(s) for a movie show, the seat(s) are reserved for the user for 5mins to make payments, then the seat(s) becomes available again. Reserved seats are queued and handled by the reservation service, it also periodically checks if the reserved seats are expired, any expired seats becomes available again for the user queued in the waiting service. User will receive notification when the seats becomes available again

## Database Table Design
![alt text](<Ticketmaster Database Tables Diagram.jpeg>)

## Application Dependencies
* Go 1.20 >=
* MySql
* Gorm
* Nats.io
* wkhtmltopdf
* Docker (Optional)

## Usage
Make sure you have the mysql service and Nats service running on the machine as specified as the dependencies for this project. Change the configuration as it suit you in the .env

To run locally, cd from the root to either apis, activereservation or waitingservice directory and run 

```
go run main.go
```

To run on Docker, cd to the root and run the docker compose file 

```
docker-compose -f docker-compose.yml up -d
```

## Upcoming Features
Frontend implementation is coming soon. If you are interested in building it , please send an Email


