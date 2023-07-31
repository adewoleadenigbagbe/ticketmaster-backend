package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const MOVIEDB_URL string = "https://api.themoviedb.org/3/movie/popular?language=en-US&page=1"
const API_KEY string = "6a4af6431ecf275b09f733a9ed14fe96"
const AUTHORIZATION = "Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiI2YTRhZjY0MzFlY2YyNzViMDlmNzMzYTllZDE0ZmU5NiIsInN1YiI6IjY0YWU3ZGVjNjZhMGQzMDEwMGRiYTFhYiIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.WS39L-os2iWGQyRJAflD_VzuWLda4BvpWkBHcXOgbG0"

var workerPoolSize = 4
var pages = make(chan int, workerPoolSize)
var movielist = make([]MovieData, 0)
var movies = []entities.Movie{}
var genres = []enums.Genre{
	enums.Action, enums.Adventure, enums.Animation, enums.Comedy,
	enums.Crime, enums.Documentary, enums.Drama, enums.Family,
	enums.Fantasy, enums.History, enums.Horror, enums.Music,
	enums.Mystery, enums.Romance, enums.ScienceFiction, enums.TVMovie,
	enums.Thriller, enums.War, enums.Western,
}
var seats = []enums.SeatType{
	enums.Gold, enums.Premium, enums.Standard,
}

// create a time alias
type JsonReleaseDate time.Time

// Implement Marshaler and Unmarshaler interface
func (j *JsonReleaseDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*j = JsonReleaseDate(t)
	return nil
}

func (j JsonReleaseDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(j))
}

type MovieData struct {
	Adult            bool            `json:"adult"`
	BackDropPath     string          `json:"backdrop_path"`
	GenreIDs         []int           `json:"genre_ids"`
	ID               int             `json:"id"`
	OriginalLanguage string          `json:"original_language"`
	OriginalTitle    string          `json:"original_title"`
	Overview         string          `json:"overview"`
	Popularity       float32         `json:"popularity"`
	PosterPath       string          `json:"poster_path"`
	ReleaseDate      JsonReleaseDate `json:"release_date"`
	Title            string          `json:"title"`
	Video            bool            `json:"video"`
	VoteAverage      float32         `json:"vote_average"`
	VoteCount        int             `json:"vote_count"`
}

type ResponseData struct {
	Page         int         `json:"page"`
	TotalPages   int         `json:"total_pages"`
	TotalResults int         `json:"total_results"`
	MovieDatas   []MovieData `json:"results"`
}

func main() {
	dsn := "root:P@ssw0r1d@tcp(127.0.0.1:3306)/?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			NoLowerCase: true,
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to the Database")

	dbName := "ticketmasterDB"
	createCommand := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", dbName)
	useDBCommand := fmt.Sprintf("USE %s;", dbName)

	db.Exec(createCommand)
	db.Exec(useDBCommand)

	err = createDataBaseEntities(db, &entities.City{},
		&entities.Show{}, &entities.Cinema{}, &entities.CinemaHall{},
		&entities.CinemaSeat{}, &entities.Show{})
	//&entities.Movie{})

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("All Tables are sucessfully created in the DB")

	// maxpage := 500
	// go AllocateJobs(maxpage)
	// CreateWorkerThread(workerPoolSize)

	// //sort the data
	// sort.Sort(byUUID(movies))
	// for _, movie := range movies {
	// 	tx := db.Create(movie)
	// 	if tx.Error != nil {
	// 		continue
	// 	}
	// }

	folderPath := "jsondata"
	getJsonData(folderPath, db)
}

func getMovieData(page int) ResponseData {
	req, _ := http.NewRequest("GET", MOVIEDB_URL, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("api_key", API_KEY)
	req.Header.Add("Authorization", AUTHORIZATION)

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

// AllocateJobs - create jobs to be done - sender
func AllocateJobs(totalPages int) {
	for i := 1; i <= totalPages; i++ {
		pages <- i
	}

	close(pages)
}

// create workers
func CreateWorkerThread(noOfWorkers int) []MovieData {
	var wg sync.WaitGroup
	for i := 1; i <= noOfWorkers; i++ {
		//means add the number of worker semaphore
		wg.Add(1)
		go worker(&wg)
	}
	wg.Wait()

	return movielist
}

// receive jobs
func worker(wg *sync.WaitGroup) {
	for page := range pages {
		fmt.Println("Page #", page)
		resp := getMovieData(page)
		AddMovieToList(resp.MovieDatas)
	}
	wg.Done()
}

func AddMovieToList(movieDatasResponse []MovieData) {
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

func getJsonData(folderPath string, db *gorm.DB) {
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	targetFolderPath := "ticketmaster-backend"
	index := strings.Index(currentWorkingDirectory, targetFolderPath)
	if index == -1 {
		log.Fatal("Target folder not found")
	}

	path := filepath.Join(currentWorkingDirectory[:index], targetFolderPath, folderPath, "\\*.json")
	files, err := filepath.Glob(path)

	if err != nil {
		log.Fatalln(err)
	}

	cities := []struct {
		Name    string `json:"city"`
		State   string `json:"state"`
		ZipCode int    `json:"zip_code"`
	}{}

	cinemas := []struct {
		Name        string `json:"city"`
		CinemaHalls int    `json:"cinemahalls"`
	}{}

	cinemahalls := []struct {
		Name       string `json:"city"`
		TotalSeats int    `json:"totalseats"`
	}{}

	for _, file := range files {
		if filepath.Ext(file) == ".json" {

			content, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}

			switch filepath.Base(file) {
			case "city.json":
				err = json.Unmarshal(content, &cities)
			case "cinema.json":
				err = json.Unmarshal(content, &cinemas)
			case "cinemahall.json":
				err = json.Unmarshal(content, &cinemahalls)
			default:
				log.Fatalln("File not available for processing")
			}

			if err != nil {
				continue
			}
		}
	}

	cityEntities := []entities.City{}
	for _, city := range cities {
		cityentity := entities.City{
			Id:           sequentialguid.New().String(),
			Name:         city.Name,
			State:        city.State,
			Zipcode:      sql.NullString{String: strconv.Itoa(city.ZipCode), Valid: true},
			IsDeprecated: false,
		}
		cityEntities = append(cityEntities, cityentity)
	}

	//sort the cities
	sort.Sort(utilities.ByCityID(cityEntities))

	cinemaEntities := []entities.Cinema{}
	for _, cinema := range cinemas {
		cinemaentity := entities.Cinema{
			Id:                sequentialguid.New().String(),
			Name:              cinema.Name,
			TotalCinemalHalls: cinema.CinemaHalls,
			CityId:            cityEntities[rand.Intn(len(cityEntities))].Id,
			IsDeprecated:      false,
		}
		cinemaEntities = append(cinemaEntities, cinemaentity)
	}

	//sort the cinemas
	sort.Sort(utilities.ByCinemaID(cinemaEntities))

	cinemaHallEntities := []entities.CinemaHall{}
	for _, cinemaHall := range cinemahalls {
		cinemahallentity := entities.CinemaHall{
			Id:           sequentialguid.New().String(),
			Name:         cinemaHall.Name,
			TotalSeat:    cinemaHall.TotalSeats,
			CinemaId:     cinemaEntities[rand.Intn(len(cinemaEntities))].Id,
			IsDeprecated: false,
		}
		cinemaHallEntities = append(cinemaHallEntities, cinemahallentity)
	}

	// sort the cinemahall
	sort.Sort(utilities.ByCinemaHallID(cinemaHallEntities))

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
	sort.Sort(utilities.ByCinemaSeatID(cinemaSeatsEntities))

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
