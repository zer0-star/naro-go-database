package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type City struct {
	ID          int    `json:"id" db:"ID"`
	Name        string `json:"name" db:"Name"`
	CountryCode string `json:"countryCode" db:"CountryCode"`
	District    string `json:"district" db:"District"`
	Population  int    `json:"population" db:"Population"`
}

func main() {
	var cityName string

	if len(os.Args) < 2 {
		cityName = "Tokyo"
	} else {
		cityName = os.Args[1]
	}

	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE")))

	if err != nil {
		log.Fatalf("Cannot connect to database: %s", err)
	}

	fmt.Println("Connected!")

	var city City

	if err := db.Get(&city, fmt.Sprintf("SELECT * FROM city WHERE Name = '%s'", cityName)); errors.Is(err, sql.ErrNoRows) {
		log.Printf("No such city that Name = '%s'", cityName)
	} else if err != nil {
		log.Fatalf("DB error: %s", err)
	}

	fmt.Printf("The population of the city '%s' is %d\n", cityName, city.Population)
}