package server

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	_ "github.com/joho/godotenv/autoload"
)

type City struct {
	ID          int    `json:"id" db:"ID"`
	Name        string `json:"name" db:"Name"`
	CountryCode string `json:"countryCode" db:"CountryCode"`
	District    string `json:"district" db:"District"`
	Population  int    `json:"population" db:"Population"`
}

type PostedCity struct {
	Name        *string `json:"name"`
	CountryCode *string `json:"countryCode"`
	District    *string `json:"district"`
	Population  *int    `json:"population"`
}

var (
	db *sqlx.DB
)

func Run(cmd *cobra.Command, args []string) {
	_db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE")))

	if err != nil {
		log.Fatalf("Cannot connect to database: %s", err)
		return
	}

	fmt.Println("Connected!")

	db = _db

	e := echo.New()

	e.GET("/city/:name", getCityHandler)

	e.POST("/city", postCityHandler)

	e.GET("/cities", getCitiesHandler)

	e.Start(":10201")
}

func getCityHandler(c echo.Context) error {
	name := c.Param("name")

	var city City

	if err := db.Get(&city, "SELECT * FROM city WHERE Name = ?", name); errors.Is(err, sql.ErrNoRows) {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("No such city that Name = '%s'", name))
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("DB error: %s", err))
	}

	return c.JSON(http.StatusOK, city)
}

func postCityHandler(c echo.Context) error {
	var posted PostedCity

	err := c.Bind(&posted)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("%s", err))
	}

	if posted.Name == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "name is not provided")
	}

	if posted.CountryCode == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "countryCode is not provided")
	}

	if posted.District == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "district is not provided")
	}

	if posted.Population == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "population is not provided")
	}

	var city City

	if err := db.Get(&city, "INSERT INTO city (Name, CountryCode, District, Population) VALUES (?, ?, ?, ?) RETURNING *", posted.Name, posted.CountryCode, posted.District, posted.Population); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("DB error: %s", err))
	}

	return c.JSON(http.StatusOK, city)
}

func getCitiesHandler(c echo.Context) error {
	limit := 10

	if c.QueryParams().Has("limit") {
		_limit, err := strconv.Atoi(c.QueryParam("limit"))

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("%s", err))
		}

		limit = _limit
	}

	orderBy := "ID"

	if c.QueryParams().Has("orderBy") {
		orderBy = c.QueryParam("orderBy")
		if !(orderBy == "ID" || orderBy == "Name" || orderBy == "CountryCode" || orderBy == "District" || orderBy == "Population") {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unknown column: %s", orderBy))
		}
	}

	order := "ASC"

	if c.QueryParams().Has("order") {
		order = c.QueryParam("order")
		if !(strings.EqualFold(order, "ASC") || strings.EqualFold(order, "DESC")) {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("order must be 'ASC' or 'DESC' (case insensitive), but got: %s", order))
		}
	}

	var cities []City

	if err := db.Select(&cities, "SELECT * FROM city ORDER BY "+orderBy+" "+order+" LIMIT ?", limit); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("DB error: %s", err))
	}

	return c.JSON(http.StatusOK, cities)
}
