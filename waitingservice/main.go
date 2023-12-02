package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Wolechacho/ticketmaster-backend/common"
	db "github.com/Wolechacho/ticketmaster-backend/database"
	"github.com/Wolechacho/ticketmaster-backend/enums"
	"github.com/labstack/echo/v4"
	"github.com/muesli/cache2go"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"

	"github.com/samber/lo"
)

type JsonSeatResponse struct {
	ShowId        string   `json:"showId"`
	CinemaSeatIds []string `json:"cinemaSeatIds"`
}

func main() {
	afterOneHourTicker := time.NewTicker(1 * time.Hour)
	afterThreeMinTicker := time.NewTicker(3 * time.Minute)
	eventChan := make(chan bool)

	// Accessing a new cache table for the first time will create it.
	cache := cache2go.Cache("waitingServiceCache")

	// connect to nats
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	//connect to the db
	db := db.ConnectToDatabase()

	subscribeToAvailableSeat(nc, cache, eventChan)

	//create an server sent event http server
	e := echo.New()

	e.GET("/subscribe", func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
		c.Response().Header().Set("Cache-Control", "no-cache")
		c.Response().Header().Set("Connection", "keep-alive")
		c.Response().WriteHeader(http.StatusOK)

		for {
			select {
			case <-c.Request().Context().Done():
				fmt.Println("Client closed connection")
				return nil
			case <-eventChan:
				seats := getAvailableSeat(cache, db)
				broadcast(c, seats)
				time.Sleep(1 * time.Second)
			case t := <-afterOneHourTicker.C:
				fmt.Println("Tick at", t)
				cache.Flush()
			case w := <-afterThreeMinTicker.C:
				fmt.Println("Tick at", w)
				availableSeats := []JsonSeatResponse{}
				if cache.Count() > 0 {
					cache.Foreach(func(key interface{}, item *cache2go.CacheItem) {
						message, ok := item.Data().(common.SeatAvailableMessage)
						if ok {
							jsonAvailableSeat := JsonSeatResponse{
								ShowId:        key.(string),
								CinemaSeatIds: message.CinemaSeatIds,
							}
							availableSeats = append(availableSeats, jsonAvailableSeat)
						}
					})
					broadcast(c, availableSeats)
				}

			default:
			}
		}
	})

	if err := e.Start(":8207"); err != nil && err != http.ErrServerClosed {
		e.Logger.Fatal("shutting down the server")
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	defer nc.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

}

func subscribeToAvailableSeat(nc *nats.Conn, cache *cache2go.CacheTable, ch chan bool) {
	nc.Subscribe(common.SeatAvailableEvent, func(m *nats.Msg) {
		buf := bytes.NewBuffer(m.Data)
		dec := gob.NewDecoder(buf)
		var message common.SeatAvailableMessage
		err := dec.Decode(&message)
		fmt.Println("Receiving Message:", message)
		if err != nil {
			log.Fatal("decode error:", err)
		}

		//add to cache
		if !cache.Exists(message.ShowId) {
			cache.Add(message.ShowId, 0, message)
		} else {
			value, _ := cache.Value(message.ShowId)
			cacheItem, _ := value.Data().(common.SeatAvailableMessage)
			cache.Delete(message.ShowId)
			message.CinemaSeatIds = append(message.CinemaSeatIds, cacheItem.CinemaSeatIds...)
			cache.Add(message.ShowId, 0, message)
		}
		ch <- true
	})
}

func getAvailableSeat(cache *cache2go.CacheTable, db *gorm.DB) []JsonSeatResponse {
	availableSeats := []JsonSeatResponse{}
	seatsToDelete := []common.SeatAvailableMessage{}
	if cache.Count() > 0 {
		cache.Foreach(func(key interface{}, item *cache2go.CacheItem) {
			message, ok := item.Data().(common.SeatAvailableMessage)
			if ok {
				//check the db one more time if the seat are still available
				showSeatsQuery, err := db.Table("showseats").
					Where("ShowId = ?", key).
					Where("CinemaSeatId IN ?", message.CinemaSeatIds).
					Where("Status = ?", enums.ExpiredSeat).
					Where("IsDeprecated = ?", false).
					Select("showseats.CinemaSeatId").
					Rows()
				if err != nil {
					return
				}

				defer showSeatsQuery.Close()

				ids := []string{}
				for showSeatsQuery.Next() {
					var cinemaSeatId string
					err = showSeatsQuery.Scan(&cinemaSeatId)

					if err != nil {
						fmt.Println(err)
						return
					}
					ids = append(ids, cinemaSeatId)
				}

				//get distinct ids
				ids = lo.Uniq(ids)

				//say if we have 3 available seat in the cache and 2 is available in the db
				//since one has already been booked/cancelled it no longer needs to be sent to the client,
				//remove from the cache and add the 2 availabe seats in the db
				intersectIds := lo.Intersect(message.CinemaSeatIds, ids)
				if len(intersectIds) > 0 {
					seat := common.SeatAvailableMessage{
						ShowId:        key.(string),
						CinemaSeatIds: intersectIds,
					}
					seatsToDelete = append(seatsToDelete, seat)

					jsonAvailableSeat := JsonSeatResponse{
						ShowId:        key.(string),
						CinemaSeatIds: intersectIds,
					}
					availableSeats = append(availableSeats, jsonAvailableSeat)
				}
			}
		})
	}

	for _, seat := range seatsToDelete {
		if cache.Exists(seat.ShowId) {
			cache.Delete(seat.ShowId)
			cache.Add(seat.ShowId, 0, seat)
		}
	}
	return availableSeats
}

func broadcast(c echo.Context, seats []JsonSeatResponse) error {
	resp, err := json.Marshal(seats)
	if err != nil {
		return nil
	}
	s := fmt.Sprintf("data: %s\n\n", string(resp))
	fmt.Fprint(c.Response(), s)
	c.Response().Flush()
	return nil
}
