package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Wolechacho/ticketmaster-backend/common"
	"github.com/Wolechacho/ticketmaster-backend/enums"
	"github.com/Wolechacho/ticketmaster-backend/models"
	"github.com/muesli/cache2go"
	"github.com/nats-io/nats.go"
	"github.com/samber/lo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type DbQuery struct {
	ShowId        string
	CinemaSeatIds []string
}

type ShowSeatDTO struct {
	Id           string
	Status       enums.BookingStatus
	CinemaSeatId string
	ShowId       string
	SeatNumber   int
}

func main() {
	var err error
	var wg sync.WaitGroup
	wg.Add(1)

	// Accessing a new cache table for the first time will create it.
	cache := cache2go.Cache("myCache")

	// connect to nats
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	//connec to the db
	rootFolderPath := "ticketmaster-backend"
	dbConfigFilePath := "configs\\database.json"

	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	index := strings.Index(currentWorkingDirectory, rootFolderPath)
	if index == -1 {
		log.Fatal("App Root Folder Path not found")
	}

	rootPath := filepath.Join(currentWorkingDirectory[:index], rootFolderPath)

	dbConfigPath := filepath.Join(rootPath, dbConfigFilePath)
	content, err := os.ReadFile(dbConfigPath)
	if err != nil {
		log.Fatalln(err)
	}

	dbConfig := models.CreateDbConfig(content)
	dsn := dbConfig.GetDsn()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			NoLowerCase: true,
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	addMessagetoCache(cache, nc)
	checkExpiredReservation(cache, db)

	wg.Wait()
}

func addMessagetoCache(cache *cache2go.CacheTable, nc *nats.Conn) {
	nc.Subscribe(common.EventSubject, func(m *nats.Msg) {
		fmt.Print("Received a message:")

		buf := bytes.NewBuffer(m.Data)
		dec := gob.NewDecoder(buf)
		var message common.BookingMessage
		err := dec.Decode(&message)
		if err != nil {
			log.Fatal("decode error:", err)
		}

		//add to cache
		cache.Lock()
		cache.Add(message.UserId, 0, message)
		cache.Unlock()

		//cache.SetAboutToDeleteItemCallback()
	})
}

func checkExpiredReservation(cache *cache2go.CacheTable, db *gorm.DB) {
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

				if cache.Count() > 0 {
					filters := []DbQuery{}
					now := time.Now()
					cache.Foreach(func(key interface{}, item *cache2go.CacheItem) {
						message, ok := item.Data().(*common.BookingMessage)
						if ok {
							if now.After(message.BookingDateTime) && now.After(message.ExpiryDateTime) {
								query := DbQuery{
									message.ShowId,
									message.CinemaSeatIds,
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
								setStatusToAvailable(db, key, v.CinemaSeatIds)
							}
						}

					}
				}
			}
		}
	}()
	<-done
	ticker.Stop()
	fmt.Println("Ticker stopped")
}

func setStatusToAvailable(db *gorm.DB, showId string, cinemaIds []string) {
	showSeatsQuery, err := db.Table("showseats").
		Where("CinemaSeatId IN ?", cinemaIds).
		Where("showseats.ShowId = ?", showId).
		Where("showseats.IsDeprecated = ?", false).
		Joins("join cinemaSeats on showseats.CinemaSeatId = cinemaSeats.Id").
		Select("showseats.Id, showseats.Status, showseats.CinemaSeatId,showseats.ShowId,cinemaSeats.SeatNumber").
		Rows()

	if err != nil {
		return
	}

	defer showSeatsQuery.Close()

	showSeats := []ShowSeatDTO{}
	for showSeatsQuery.Next() {
		var showSeatDTO ShowSeatDTO
		err = showSeatsQuery.Scan(&showSeatDTO.Id, &showSeatDTO.Status, &showSeatDTO.CinemaSeatId, &showSeatDTO.ShowId, &showSeatDTO.SeatNumber)
		if err != nil {
			fmt.Println(err)
		}
		showSeats = append(showSeats, showSeatDTO)
	}

	//update
	_ = db.Transaction(func(tx *gorm.DB) error {
		for _, showSeat := range showSeats {
			if showSeat.Status == enums.Reserved || showSeat.Status == enums.PendingBook {
				db.Table("showseats").Where("ShowId = ? AND CinemaSeatId = ?", showSeat.ShowId, showSeat.CinemaSeatId).Update("status", enums.Available)
			}
		}
		return nil
	})
}
