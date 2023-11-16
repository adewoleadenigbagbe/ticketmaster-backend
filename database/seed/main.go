package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Wolechacho/ticketmaster-backend/database/entities"
	"github.com/Wolechacho/ticketmaster-backend/enums"
	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"github.com/Wolechacho/ticketmaster-backend/helpers/utilities"
	"github.com/Wolechacho/ticketmaster-backend/models"
	"github.com/samber/lo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	RootFolderPath   = "ticketmaster-backend"
	DbConfigFilePath = "configs\\database.json"
	MovieApiFilePath = "configs\\movieapi.json"

	genres = []enums.Genre{
		enums.Action, enums.Adventure, enums.Animation, enums.Comedy,
		enums.Crime, enums.Documentary, enums.Drama, enums.Family,
		enums.Fantasy, enums.History, enums.Horror, enums.Music,
		enums.Mystery, enums.Romance, enums.ScienceFiction, enums.TVMovie,
		enums.Thriller, enums.War, enums.Western,
	}
	seats = []enums.SeatType{
		enums.Gold, enums.Premium, enums.Standard,
	}
	movies = []entities.Movie{}
)

const (
	workerPoolSize = 4
	MaxPage        = 500
)

func main() {
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	index := strings.Index(currentWorkingDirectory, RootFolderPath)
	if index == -1 {
		log.Fatal("App Root Folder Path not found")
	}

	rootPath := filepath.Join(currentWorkingDirectory[:index], RootFolderPath)

	dbConfigPath := filepath.Join(rootPath, DbConfigFilePath)
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

	fmt.Println("SucessFully Connected to the Database")

	//create database entities
	createCommand := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", dbConfig.DatabaseName)
	useDBCommand := fmt.Sprintf("USE %s;", dbConfig.DatabaseName)

	db.Exec(createCommand)
	db.Exec(useDBCommand)

	err = createDataBaseEntities(db, &entities.City{},
		&entities.Show{}, &entities.Cinema{}, &entities.CinemaHall{},
		&entities.CinemaSeat{}, &entities.Show{}, &entities.Movie{})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("All Tables are sucessfully created in the DB")

	filedata := NewFileData(rootPath)
	filedata.GetData(db)

	apiJsonPath := filepath.Join(rootPath, MovieApiFilePath)
	apicontent, err := os.ReadFile(apiJsonPath)
	if err != nil {
		log.Fatalln(err)
	}

	movieApiConfig := models.CreateMovieApiConfig(apicontent)
	apiData := NewApiData()
	apiData.GetData(movieApiConfig, db)
}

type IData interface {
	GetData(db *gorm.DB)
}

type FileData struct {
	JsonFolderPath   string
	TargetFolderPath string
	Converter        func(data []byte, v any) error
	Cities           []struct {
		Name      string  `json:"city"`
		State     string  `json:"state"`
		ZipCode   int     `json:"zip_code"`
		Latitude  float32 `json:"latitude"`
		Longitude float32 `json:"longitude"`
	}
	Cinemas []struct {
		Name        string `json:"name"`
		CinemaHalls int    `json:"cinemahalls"`
	}

	Cinemahalls []struct {
		Name       string `json:"name"`
		TotalSeats int    `json:"totalseats"`
	}
}

func NewFileData(folderPath string) *FileData {
	fileData := &FileData{
		JsonFolderPath:   "jsondata",
		TargetFolderPath: folderPath,
		Converter: func(data []byte, v any) error {
			err := json.Unmarshal(data, v)
			return err
		},
		Cities: []struct {
			Name      string  `json:"city"`
			State     string  `json:"state"`
			ZipCode   int     `json:"zip_code"`
			Latitude  float32 `json:"latitude"`
			Longitude float32 `json:"longitude"`
		}{},

		Cinemas: []struct {
			Name        string `json:"name"`
			CinemaHalls int    `json:"cinemahalls"`
		}{},

		Cinemahalls: []struct {
			Name       string `json:"name"`
			TotalSeats int    `json:"totalseats"`
		}{},
	}

	return fileData
}

func (fileData *FileData) GetData(db *gorm.DB) {
	path := filepath.Join(fileData.TargetFolderPath, fileData.JsonFolderPath, "\\*.json")
	files, err := filepath.Glob(path)

	if err != nil {
		log.Fatalln(err)
	}

	for _, file := range files {
		if filepath.Ext(file) == ".json" {

			content, err := os.ReadFile(file)
			if err != nil {
				continue
			}

			switch filepath.Base(file) {
			case "city.json":
				_ = fileData.Converter(content, &fileData.Cities)
			case "cinema.json":
				_ = fileData.Converter(content, &fileData.Cinemas)
			case "cinemahall.json":
				_ = fileData.Converter(content, &fileData.Cinemahalls)
			default:
				continue
			}
		}
	}

	cityEntities := []entities.City{}
	for _, city := range fileData.Cities {
		cityentity := entities.City{
			Id:      sequentialguid.New().String(),
			Name:    city.Name,
			State:   city.State,
			Zipcode: strconv.Itoa(city.ZipCode),
			Coordinates: entities.Coordinate{
				Longitude: city.Longitude,
				Latitude:  city.Latitude,
			},
			IsDeprecated: false,
		}
		cityEntities = append(cityEntities, cityentity)
	}

	//sort the cities
	cityEntities = lo.UniqBy(cityEntities, func(city entities.City) string {
		return city.Id
	})
	sort.Sort(entities.ByID[entities.City](cityEntities))

	cinemaEntities := []entities.Cinema{}
	for _, cinema := range fileData.Cinemas {
		cinemaentity := entities.Cinema{
			Id:                sequentialguid.New().String(),
			Name:              cinema.Name,
			TotalCinemalHalls: cinema.CinemaHalls,
			CityId:            cityEntities[rand.Intn(len(cityEntities))].Id,
			IsDeprecated:      false,
		}
		cinemaEntities = append(cinemaEntities, cinemaentity)
	}

	// add the address of cinema

	cinemaEntities = lo.UniqBy(cinemaEntities, func(cinema entities.Cinema) string {
		return cinema.Id
	})
	//sort the cinemas
	sort.Sort(entities.ByID[entities.Cinema](cinemaEntities))

	cinemaHallEntities := []entities.CinemaHall{}
	for _, cinemaHall := range fileData.Cinemahalls {
		cinemahallentity := entities.CinemaHall{
			Id:           sequentialguid.New().String(),
			Name:         cinemaHall.Name,
			TotalSeat:    cinemaHall.TotalSeats,
			CinemaId:     cinemaEntities[rand.Intn(len(cinemaEntities))].Id,
			IsDeprecated: false,
		}
		cinemaHallEntities = append(cinemaHallEntities, cinemahallentity)
	}

	cinemaHallEntities = lo.UniqBy(cinemaHallEntities, func(cinemaHall entities.CinemaHall) string {
		return cinemaHall.Id
	})
	// sort the cinemahall
	sort.Sort(entities.ByID[entities.CinemaHall](cinemaHallEntities))

	cinemaSeatsEntities := []entities.CinemaSeat{}
	for _, cinemaHallEntity := range cinemaHallEntities {
		for i := 1; i <= cinemaHallEntity.TotalSeat; i++ {
			cinemaSeat := entities.CinemaSeat{
				Id:           sequentialguid.New().String(),
				SeatNumber:   i,
				Type:         rand.Intn(len(seats)),
				CinemaHallId: cinemaHallEntity.Id,
				IsDeprecated: false,
			}
			cinemaSeatsEntities = append(cinemaSeatsEntities, cinemaSeat)
		}
	}

	//sort the entities cinema seats
	cinemaSeatsEntities = lo.UniqBy(cinemaSeatsEntities, func(cinemaSeat entities.CinemaSeat) string {
		return cinemaSeat.Id
	})
	sort.Sort(entities.ByID[entities.CinemaSeat](cinemaSeatsEntities))

	err = db.Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		for _, city := range cityEntities {
			if err := tx.Create(&city).Error; err != nil {
				// return any error will rollback
				return err
			}
		}

		for _, cinema := range cinemaEntities {
			if err := tx.Create(&cinema).Error; err != nil {
				// return any error will rollback
				return err
			}
		}

		for _, cinemaHall := range cinemaHallEntities {
			if err := tx.Create(&cinemaHall).Error; err != nil {
				// return any error will rollback
				return err
			}
		}

		for _, cinemaSeat := range cinemaSeatsEntities {
			if err := tx.Create(&cinemaSeat).Error; err != nil {
				// return any error will rollback
				return err
			}
		}

		// return nil will commit the whole transaction
		return nil
	})

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("All Data sucessfully saved into the newly created tables")
}

type ApiData struct {
	MovieApiConfig models.MovieApiConfig
	WorkerPoolSize int
	Pages          chan int
	Wg             *sync.WaitGroup
}

func NewApiData() *ApiData {
	apiData := &ApiData{
		WorkerPoolSize: workerPoolSize,
		Pages:          make(chan int, workerPoolSize),
		Wg:             &sync.WaitGroup{},
	}

	return apiData
}

func (apiData *ApiData) GetData(config models.MovieApiConfig, db *gorm.DB) {
	go apiData.allocateJobs(MaxPage)
	apiData.createWorkerThread(config, apiData.WorkerPoolSize)

	//sort the data
	sort.Sort(entities.ByID[entities.Movie](movies))
	for _, movie := range movies {
		tx := db.Create(&movie)
		if tx.Error != nil {
			continue
		}
	}
}

// AllocateJobs - create jobs to be done - sender
func (apiData *ApiData) allocateJobs(totalPages int) {
	for i := 1; i <= totalPages; i++ {
		apiData.Pages <- i
	}

	close(apiData.Pages)
}

// create workers
func (apiData *ApiData) createWorkerThread(movieConfig models.MovieApiConfig, noOfWorkers int) {
	var wg sync.WaitGroup
	for i := 1; i <= noOfWorkers; i++ {
		//means add the number of worker semaphore
		wg.Add(1)
		go apiData.worker(movieConfig, &wg)
	}
	wg.Wait()
}

// receive jobs
func (apiData *ApiData) worker(movieConfig models.MovieApiConfig, wg *sync.WaitGroup) {
	for page := range apiData.Pages {
		fmt.Println("Page #", page)
		resp := getMovieData(movieConfig, page)
		addMovieToList(resp.MovieDatas)
	}
	wg.Done()
}

func getMovieData(movieConfig models.MovieApiConfig, page int) ResponseData {
	url := fmt.Sprintf("%s&page=%d", movieConfig.Url, page)
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("api_key", movieConfig.ApiKey)
	req.Header.Add("Authorization", movieConfig.Auth)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var responseData ResponseData
	err := json.Unmarshal(body, &responseData)

	if err != nil {
		fmt.Println(err)
	}
	return responseData
}

func addMovieToList(movieDatasResponse []MovieData) {
	for _, moviedata := range movieDatasResponse {
		movie := entities.Movie{
			Id:           sequentialguid.New().String(),
			Title:        moviedata.OriginalTitle,
			Description:  sql.NullString{String: moviedata.Overview, Valid: true},
			Duration:     sql.NullInt32{Valid: false},
			ReleaseDate:  time.Time(moviedata.ReleaseDate),
			Genre:        rand.Intn(len(genres)),
			Language:     moviedata.OriginalLanguage,
			Popularity:   moviedata.Popularity,
			VoteCount:    moviedata.VoteCount,
			IsDeprecated: false,
		}
		movies = append(movies, movie)
	}
}

func createDataBaseEntities(db *gorm.DB, entities ...interface{}) error {
	for _, entity := range entities {
		if !db.Migrator().HasTable(entity) {
			err := db.Migrator().CreateTable(entity)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type MovieData struct {
	Adult            bool               `json:"adult"`
	BackDropPath     string             `json:"backdrop_path"`
	GenreIDs         []int              `json:"genre_ids"`
	ID               int                `json:"id"`
	OriginalLanguage string             `json:"original_language"`
	OriginalTitle    string             `json:"original_title"`
	Overview         string             `json:"overview"`
	Popularity       float32            `json:"popularity"`
	PosterPath       string             `json:"poster_path"`
	ReleaseDate      utilities.Datetime `json:"release_date"`
	Title            string             `json:"title"`
	Video            bool               `json:"video"`
	VoteAverage      float32            `json:"vote_average"`
	VoteCount        int                `json:"vote_count"`
}

type ResponseData struct {
	Page         int         `json:"page"`
	TotalPages   int         `json:"total_pages"`
	TotalResults int         `json:"total_results"`
	MovieDatas   []MovieData `json:"results"`
}
