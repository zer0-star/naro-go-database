package city

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	_ "github.com/joho/godotenv/autoload"
)

type City struct {
	ID          int    `json:"id" db:"ID"`
	Name        string `json:"name" db:"Name"`
	CountryCode string `json:"countryCode" db:"CountryCode"`
	District    string `json:"district" db:"District"`
	Population  int    `json:"population" db:"Population"`
}

type Country struct {
	Name       string `json:"name" db:"Name"`
	Population int    `json:"population" db:"Population"`
}

func Run(cmd *cobra.Command, args []string) {
	cityName := args[0]

	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE")))

	if err != nil {
		log.Fatalf("Cannot connect to database: %s", err)
		return
	}

	fmt.Println("Connected!")

	var (
		city    City
		country Country
	)

	if err := db.Get(&city, "SELECT * FROM city WHERE Name = ?", cityName); errors.Is(err, sql.ErrNoRows) {
		log.Printf("No such city that Name = '%s'", cityName)
		return
	} else if err != nil {
		log.Fatalf("DB error: %s", err)
		return
	}

	if err := db.Get(&country, fmt.Sprintf("SELECT Name, Population FROM country WHERE Code = '%s'", city.CountryCode)); err != nil {
		log.Fatalf("DB error: %s", err)
		return
	}

	fmt.Printf("The population of the city '%s' is %d\n", cityName, city.Population)
	fmt.Printf("The population of the city '%s' is %f%% of the population of its country '%s'\n",
		city.Name, float64(city.Population)/float64(country.Population)*100, country.Name)
}
