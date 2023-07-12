package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const MOVIEDB_URL string = "https://api.themoviedb.org/3/movie/popular?language=en-US&page=1"
const API_KEY string = "6a4af6431ecf275b09f733a9ed14fe96"
const AUTHORIZATION = "Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiI2YTRhZjY0MzFlY2YyNzViMDlmNzMzYTllZDE0ZmU5NiIsInN1YiI6IjY0YWU3ZGVjNjZhMGQzMDEwMGRiYTFhYiIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.WS39L-os2iWGQyRJAflD_VzuWLda4BvpWkBHcXOgbG0"

var workerPoolSize = 4
var pages = make(chan int, workerPoolSize)
var movielist = make([]MovieData, 0)

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
	//newTime := JsonReleaseDate(time.Now())
	//fmt.Printf("%+v\n", newTime)
	resp := getMovieData(1)
	AddMovieToList(resp.MovieDatas)

	go AllocateJobs(resp.TotalPages)
	CreateWorkerThread(workerPoolSize)

	fmt.Println(movielist)
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
	fmt.Println("Number of total pages", totalPages)
	for i := 2; i <= totalPages; i++ {
		fmt.Println("Job started running #", i)
		pages <- i
	}

	close(pages)
}

// create workers
func CreateWorkerThread(noOfWorkers int) []MovieData {
	var wg sync.WaitGroup
	for i := 1; i <= noOfWorkers; i++ {
		fmt.Println("Worker #", i)
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

func AddMovieToList(movies []MovieData) {
	movielist = append(movielist, movies...)
}
