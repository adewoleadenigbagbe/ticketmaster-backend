package models

import (
	"encoding/json"
	"fmt"
	"log"
)

type DatabaseConfig struct {
	DatabaseName string `json:"dbName"`
	User         string `json:"user"`
	Password     string `json:"password"`
	IpAddress    string `json:"ip"`
	Port         string `json:"port"`
	Charset      string `json:"charset"`
	ParseTime    string `json:"parseTime"`
	Loc          string `json:"loc"`
}

func CreateDbConfig(content []byte) *DatabaseConfig {
	dbConfig := &DatabaseConfig{}
	err := json.Unmarshal(content, dbConfig)
	if err != nil {
		log.Fatalln(err)
	}

	return dbConfig
}

func (dbConfig *DatabaseConfig) GetDsn() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=%s&parseTime=%s&loc=%s", dbConfig.User, dbConfig.Password,
		dbConfig.IpAddress, dbConfig.Port, dbConfig.Charset, dbConfig.ParseTime, dbConfig.Loc)

	return dsn
}

type MovieApiConfig struct {
	Url    string `json:"url"`
	ApiKey string `json:"apiKey"`
	Auth   string `json:"auth"`
}

func CreateMovieApiConfig(content []byte) MovieApiConfig {
	movieApiConfig := MovieApiConfig{}
	err := json.Unmarshal(content, &movieApiConfig)
	if err != nil {
		log.Fatalln(err)
	}

	return movieApiConfig
}
