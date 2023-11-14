package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Wolechacho/ticketmaster-backend/common"
	db "github.com/Wolechacho/ticketmaster-backend/database"
	"github.com/Wolechacho/ticketmaster-backend/enums"
	"github.com/muesli/cache2go"
	"github.com/nats-io/nats.go"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type DbQuery struct {
	ItemKey       string
	ShowId        string
	CinemaSeatIds []string
	UserId        string
}

type ShowSeatDTO struct {
	Id           string
	Status       enums.ShowSeatStatus
	CinemaSeatId string
	ShowId       string
	SeatNumber   int
	UserId       string
}

func main() {
	var err error
	var wg sync.WaitGroup
	wg.Add(1)

	// Accessing a new cache table for the first time will create it.
	cache := cache2go.Cache("reservationCache")

	// connect to nats
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	//connec to the db
	db := db.ConnectToDatabase()

	addMessagetoCache(cache, nc, db)
	setSeatStatusAfterExpiration(cache, db)
	wg.Wait()
}

func addMessagetoCache(cache *cache2go.CacheTable, nc *nats.Conn, db *gorm.DB) {
	nc.Subscribe(common.EventSubject, func(m *nats.Msg) {
		buf := bytes.NewBuffer(m.Data)
		dec := gob.NewDecoder(buf)
		var message common.BookingMessage
		err := dec.Decode(&message)
		fmt.Println("Receiving Message:", message)
		if err != nil {
			log.Fatal("decode error:", err)
		}

		//make a unique , append with nano time
		now := time.Now().UnixNano()
		key := fmt.Sprint(now, message.UserId)

		//add to cache
		if !cache.Exists(key) {
			cache.Add(key, 0, message)
		}
	})
}

func setSeatStatusAfterExpiration(cache *cache2go.CacheTable, db *gorm.DB) {
	ticker := time.NewTicker(5 * time.Minute)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
				//loop through the cache , find all items that are expired
				// update the show seat status and push message to the broker about an expired
				//expired reserve seat now available , while doing so
				// check if the seat is already book while seating in cache
				// do nothing if seat is already book
				checkAndSetExpiredItems(cache, db)
			}
		}
	}()
	<-done
	ticker.Stop()
	fmt.Println("Ticker stopped")
}

func setStatusToAvailable(db *gorm.DB, cache *cache2go.CacheTable, showId string, filter DbQuery) {
	showSeatsQuery, err := db.Table("showseats").
		Where("CinemaSeatId IN ?", filter.CinemaSeatIds).
		Where("showseats.ShowId = ?", showId).
		Where("showseats.IsDeprecated = ?", false).
		Joins("join cinemaSeats on showseats.CinemaSeatId = cinemaSeats.Id").
		Joins("join bookings on showseats.BookingId = bookings.Id").
		Select("showseats.Id, showseats.Status, showseats.CinemaSeatId,showseats.ShowId,cinemaSeats.SeatNumber,bookings.UserId").
		Rows()

	if err != nil {
		return
	}

	defer showSeatsQuery.Close()

	showSeats := []ShowSeatDTO{}
	for showSeatsQuery.Next() {
		var showSeatDTO ShowSeatDTO
		err = showSeatsQuery.Scan(&showSeatDTO.Id,
			&showSeatDTO.Status,
			&showSeatDTO.CinemaSeatId,
			&showSeatDTO.ShowId,
			&showSeatDTO.SeatNumber,
			&showSeatDTO.UserId)

		if err != nil {
			fmt.Println(err)
			return
		}
		showSeats = append(showSeats, showSeatDTO)
	}

	//update
	_ = db.Transaction(func(tx *gorm.DB) error {
		for _, showSeat := range showSeats {
			if showSeat.Status == enums.Reserved || showSeat.Status == enums.PendingAssignment {
				var dbErr error
				dbErr = db.Transaction(func(tx *gorm.DB) error {
					dbErr = db.Table("showseats").
						Where("ShowId = ? AND CinemaSeatId = ?", showSeat.ShowId, showSeat.CinemaSeatId).
						Where("IsDeprecated = ?", false).
						Update("status", enums.Available).Error

					if dbErr != nil {
						return dbErr
					}

					dbErr = db.Table("bookings").
						Where("ShowId = ? AND UserId = ?", showSeat.ShowId, showSeat.UserId).
						Update("status", enums.Expired).Error

					if dbErr != nil {
						return dbErr
					}
					return nil
				})
			}
		}

		//no error
		cache.Delete(filter.ItemKey)
		return nil
	})
}

func checkAndSetExpiredItems(cache *cache2go.CacheTable, db *gorm.DB) {
	if cache.Count() > 0 {
		filters := []DbQuery{}
		now := time.Now()
		cache.Foreach(func(key interface{}, item *cache2go.CacheItem) {
			message, ok := item.Data().(common.BookingMessage)
			if ok {
				if now.After(message.BookingDateTime) && now.After(message.ExpiryDateTime) {
					itemKey, _ := key.(string)
					query := DbQuery{
						ItemKey:       itemKey,
						ShowId:        message.ShowId,
						CinemaSeatIds: message.CinemaSeatIds,
						UserId:        message.UserId,
					}
					filters = append(filters, query)
				}
			}
		})

		//group them
		mapQueries := lo.GroupBy(filters, func(item DbQuery) string {
			return item.ShowId
		})

		if len(mapQueries) > 0 {
			for key, val := range mapQueries {
				for _, v := range val {
					setStatusToAvailable(db, cache, key, v)
				}
			}

		}
	}
}
