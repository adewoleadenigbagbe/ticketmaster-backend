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

func (dbConfig DatabaseConfig) NewConfig(content []byte) string {
	err := json.Unmarshal(content, &dbConfig)
	if err != nil {
		log.Fatalln(err)
	}
	//dsn := "root:P@ssw0r1d@tcp(127.0.0.1:3306)/?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=%s&parseTime=%s&loc=%s", dbConfig.User, dbConfig.Password,
		dbConfig.IpAddress, dbConfig.Port, dbConfig.Charset, dbConfig.ParseTime, dbConfig.Loc)

	return dsn
}
