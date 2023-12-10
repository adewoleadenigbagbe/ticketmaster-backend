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
	"github.com/google/uuid"
	"github.com/muesli/cache2go"
	"github.com/nats-io/nats.go"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type DbQuery struct {
	ItemKey       string
	ShowId        string
	BookingId     string
	CinemaSeatIds []string
	UserId        string
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
	setSeatStatusAfterExpiration(cache, db, nc)
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

		key := uuid.New().String()
		//add to cache
		if !cache.Exists(key) {
			cache.Add(key, 0, message)
		}
	})
}

func setSeatStatusAfterExpiration(cache *cache2go.CacheTable, db *gorm.DB, nc *nats.Conn) {
	ticker := time.NewTicker(5 * time.Minute)
	tickerForMissedCache := time.NewTicker(10 * time.Minute)

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
				checkAndSetExpiredItems(cache, db, nc)
			case v := <-tickerForMissedCache.C:
				fmt.Println("Tick at", v, "for missed cache")
				UpdateMissedBookingMessage(db, cache)
			}

		}
	}()
	<-done
	ticker.Stop()
	fmt.Println("Ticker stopped")
}

func setStatusToAvailable(db *gorm.DB, cache *cache2go.CacheTable, nc *nats.Conn, bookingId string, filter DbQuery) error {
	var (
		err       error
		showSeats []common.ShowSeatDTO
	)

	showSeatsQuery, err := db.Table("bookings").
		Where("bookings.Id = ?", bookingId).
		Where("bookings.ShowId = ?", filter.ShowId).
		Where("bookings.IsDeprecated = ?", false).
		Joins("join showseats on bookings.Id = showseats.bookingId").
		Where("showseats.ShowId = ?", filter.ShowId).
		Where("showseats.CinemaSeatId IN ?", filter.CinemaSeatIds).
		Where("showseats.IsDeprecated = ?", false).
		Joins("join cinemaSeats on showseats.CinemaSeatId = cinemaSeats.Id").
		Select("showseats.Id, showseats.Status, showseats.CinemaSeatId,showseats.ShowId,cinemaSeats.SeatNumber,bookings.UserId").
		Rows()

	if err != nil {
		return err
	}

	defer showSeatsQuery.Close()

	for showSeatsQuery.Next() {
		var showSeatDTO common.ShowSeatDTO
		err = showSeatsQuery.Scan(&showSeatDTO.Id,
			&showSeatDTO.Status,
			&showSeatDTO.CinemaSeatId,
			&showSeatDTO.ShowId,
			&showSeatDTO.SeatNumber,
			&showSeatDTO.UserId)

		if err != nil {
			return err
		}
		showSeats = append(showSeats, showSeatDTO)
	}

	showSeats = lo.Filter(showSeats, func(item common.ShowSeatDTO, index int) bool {
		return item.Status == enums.Reserved || item.Status == enums.PendingAssignment
	})

	//update
	err = db.Transaction(func(tx *gorm.DB) error {
		err = tx.Table("bookings").
			Where("Id = ?", filter.BookingId).
			Where("IsDeprecated = ?", false).
			Where("ShowId = ? ", filter.ShowId).
			Where("ShowId = ? ", filter.UserId).
			Update("status", enums.Expired).Error

		if err != nil {
			return err
		}

		err = db.Table("showseats").
			Where("bookingId = ?", filter.BookingId).
			Where("ShowId = ? ", filter.ShowId).
			Where("CinemaSeatId IN ?", showSeats).
			Where("IsDeprecated = ?", false).
			Update("status", enums.ExpiredSeat).Error

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	//no error
	cache.Delete(filter.ItemKey)

	//sample message
	cinemaIds := lo.Map(showSeats, func(item common.ShowSeatDTO, index int) string {
		return item.CinemaSeatId
	})
	bk := common.SeatAvailableMessage{
		CinemaSeatIds: cinemaIds,
		ShowId:        filter.ShowId,
	}

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err = encoder.Encode(bk)
	if err != nil {
		log.Fatal("encode error:", err)
	}

	if err := nc.Publish(common.SeatAvailableEvent, buf.Bytes()); err != nil {
		log.Fatal(err)
	}

	nc.Flush()

	return nil
}

func checkAndSetExpiredItems(cache *cache2go.CacheTable, db *gorm.DB, nc *nats.Conn) {
	var (
		err error
	)
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
						BookingId:     message.BookingId,
						CinemaSeatIds: message.CinemaSeatIds,
						UserId:        message.UserId,
					}
					filters = append(filters, query)
				}
			}
		})

		//group them
		mapQueries := lo.GroupBy(filters, func(item DbQuery) string {
			return item.BookingId
		})

		if len(mapQueries) > 0 {
			for key, val := range mapQueries {
				for _, v := range val {
					err = setStatusToAvailable(db, cache, nc, key, v)
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}
	}
}

// UpdateMissedBookingMessage pull missed cache from the DB and push to the cache in case NATS is down
func UpdateMissedBookingMessage(db *gorm.DB, cache *cache2go.CacheTable) error {
	var (
		err error
	)

	bookinqQuery, err := db.Table("bookings").
		Where("bookings.IsDeprecated = ?", false).
		Where("bookings.Status = ?", enums.PendingBook).
		Joins("join showseats on bookings.Id = showseats.bookingId").
		Where("showseats.Status = ? OR showseats.Status = ?", enums.PendingAssignment, enums.Reserved).
		Where("showseats.IsDeprecated = ?", false).
		Joins("join cinemaSeats on showseats.CinemaSeatId = cinemaSeats.Id").
		Select("bookings.Id,bookings.ShowId,bookings.UserId,showseats.CinemaSeatId,showseats.Status").
		Rows()

	if err != nil {
		return err
	}

	defer bookinqQuery.Close()

	var missedBookings []common.BookingDTO
	for bookinqQuery.Next() {
		var booking common.BookingDTO
		err = bookinqQuery.Scan(booking.BookingId, booking.ShowId, booking.UserId, booking.CinemaSeatId, booking.Status)
		if err != nil {
			return err
		}
		missedBookings = append(missedBookings, booking)
	}

	groupBookings := lo.GroupBy(missedBookings, func(item common.BookingDTO) string {
		return item.BookingId
	})

	if cache.Count() == 0 {
		for key, bookings := range groupBookings {
			today := time.Now()
			var message common.BookingMessage
			message.BookingId = key
			message.BookingDateTime = today
			message.ExpiryDateTime = today.Add(5 * time.Minute)

			for _, booking := range bookings {
				if message.BookingId != "" {
					message.ShowId = booking.ShowId
					message.UserId = booking.UserId
					message.Status = booking.Status
				}
				message.CinemaSeatIds = append(message.CinemaSeatIds, booking.CinemaSeatId)
			}

			itemKey := uuid.New().String()
			cache.Add(itemKey, 0, message)
		}
		return nil
	}

	cache.Foreach(func(key interface{}, item *cache2go.CacheItem) {
		message, ok := item.Data().(common.BookingMessage)
		if ok {
			for key, bookings := range groupBookings {
				if key == message.BookingId {
					continue
				}
				today := time.Now()
				var message common.BookingMessage
				message.BookingId = key
				message.BookingDateTime = today
				message.ExpiryDateTime = today.Add(5 * time.Minute)

				for _, booking := range bookings {
					if message.BookingId != "" {
						message.ShowId = booking.ShowId
						message.UserId = booking.UserId
						message.Status = booking.Status
					}
					message.CinemaSeatIds = append(message.CinemaSeatIds, booking.CinemaSeatId)
				}

				itemKey := uuid.New().String()
				cache.Add(itemKey, 0, message)
			}
		}
	})
	return nil
}
