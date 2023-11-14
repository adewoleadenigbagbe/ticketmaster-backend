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
	"github.com/labstack/echo/v4"
	"github.com/muesli/cache2go"
	"github.com/nats-io/nats.go"
	"github.com/samber/lo"
)

type JsonSeatResponse struct {
	ShowId        string   `json:"showId"`
	CinemaSeatIds []string `json:"cinemaSeatIds"`
}
type SeatModel struct {
	CinemaSeatId string
}

func main() {
	ticker := time.NewTicker(1 * time.Hour)
	eventChan := make(chan bool)
	// Accessing a new cache table for the first time will create it.
	cache := cache2go.Cache("waitingServiceCache")
	// connect to nats
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	//connec to the db
	db := db.ConnectToDatabase()

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
		}

		eventChan <- true
	})

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
				if cache.Count() > 0 {
					resp := []JsonSeatResponse{}
					cache.Foreach(func(key interface{}, item *cache2go.CacheItem) {
						message, ok := item.Data().(common.SeatAvailableMessage)
						if ok {
							showSeatsQuery, err := db.Table("showseats").
								Where("ShowId = ?", key).
								Where("CinemaSeatId IN ?", message.CinemaSeatIds).
								//TODO: change this to enum state when merge
								Where("Status = ?", 1).
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
							seat := JsonSeatResponse{
								ShowId:        message.ShowId,
								CinemaSeatIds: ids,
							}
							resp = append(resp, seat)

							//remove the seatIds tha is not available in the cache
							leftIds, _ := lo.Difference(message.CinemaSeatIds, ids)
							cache.Delete(message.ShowId)
							availableSeat := common.SeatAvailableMessage{
								ShowId:        message.ShowId,
								CinemaSeatIds: leftIds,
							}
							cache.Add(message.ShowId, 0, availableSeat)
						}
					})

					data, err := json.Marshal(resp)
					if err != nil {
						return err
					}

					fmt.Fprint(c.Response(), data)
					c.Response().Flush()
				}

				time.Sleep(1 * time.Second)
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
				cache.Flush()
			default:
				fmt.Fprint(c.Response(), "data: hi\n\n")
				c.Response().Flush()
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
