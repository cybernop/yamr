package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type reading struct {
	KindId     int       `json:"kind_id,omitempty"`
	RecordedOn time.Time `json:"recordedOn"`
	Reading    float64   `json:"reading"`
}

type reading_list struct {
	Readings []reading `json:"readings"`
	Kind     string    `json:"kind"`
	Unit     string    `json:"unit"`
}

type kind struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Unit string `json:"unit"`
}

type kind_list struct {
	Kinds []kind `json:"kinds"`
}

type response_error struct {
	Error string `json:"error"`
}

var db_url = "postgres://" + os.Getenv("POSTGRES_USER") + ":" + os.Getenv("POSTGRES_PASSWORD") + "@" + os.Getenv("POSTGRES_HOSTNAME") + ":5432/" + os.Getenv("POSTGRES_DB")

func main() {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	router.GET("/kinds", getKinds)
	router.POST("/kind", postKind)
	router.GET("/readings", getReadings)
	router.POST("/reading", postReading)

	router.Run("localhost:8080")
}

func getKinds(c *gin.Context) {
	// Establish connection
	conn, err := pgx.Connect(context.Background(), db_url)
	if err != nil {
		interal_error(c, "Unable to connect to database", err)
		return
	}
	defer conn.Close(context.Background())

	// Query for kinds
	rows, err := conn.Query(context.Background(), "SELECT kind_id, kind_name, unit FROM kinds")
	if err != nil {
		interal_error(c, "Unable to get readings from database", err)
		return
	}

	// Convert to kinds
	var kk kind_list
	for rows.Next() {
		var k kind
		rows.Scan(&k.Id, &k.Name, &k.Unit)
		kk.Kinds = append(kk.Kinds, k)
	}

	c.IndentedJSON(http.StatusOK, kk)
}

func postKind(c *gin.Context) {
	// Establish connection
	conn, err := pgx.Connect(context.Background(), db_url)
	if err != nil {
		interal_error(c, "Unable to connect to database", err)
		return
	}
	defer conn.Close(context.Background())

	// Call BindJSON to bind the received JSON to
	var k kind
	if err := c.BindJSON(&k); err != nil {
		bad_request_err(c, "Failed to parse argument", err)
		return
	}

	// Build query
	query := "INSERT INTO kinds (kind_name, unit) VALUES (@name, @unit)"
	args := pgx.NamedArgs{
		"name": k.Name,
		"unit": k.Unit,
	}

	// Execute query
	_, err = conn.Exec(context.Background(), query, args)
	if err != nil {
		interal_error(c, "Failed to create new kind", err)
		return
	}

	// Build query
	query = "SELECT kind_id, kind_name FROM kinds ORDER BY kind_id DESC LIMIT 1"

	// Query for readings
	err = conn.QueryRow(context.Background(), query).Scan(&k.Id, &k.Name)
	if err != nil {
		interal_error(c, "Unable to get kinds from database", err)
		return
	}

	c.IndentedJSON(http.StatusCreated, k)
}

func getReadings(c *gin.Context) {
	// Establish connection
	conn, err := pgx.Connect(context.Background(), db_url)
	if err != nil {
		interal_error(c, "Unable to connect to database", err)
		return
	}
	defer conn.Close(context.Background())

	var kind string
	// Query the readings
	if k := c.Request.URL.Query()["kind"]; len(k) == 1 {
		kind = k[0]
		// Filter by kind
		if len(k) == 1 {

		} else {
			bad_request(c, "Multiple 'kind' parameters")
			return
		}
	} else {
		bad_request(c, "Incorrect number of 'kind' parameters")
		return
	}

	// Build query
	query := "SELECT recorded_on, reading FROM readings WHERE kind_id = @kind ORDER BY recorded_on"
	args := pgx.NamedArgs{
		"kind": kind,
	}

	// Query for readings
	rows, err := conn.Query(context.Background(), query, args)
	if err != nil {
		interal_error(c, "Unable to get readings from database", err)
		return
	}

	// Convert query results to structs
	var rr reading_list
	var previousDate time.Time
	for rows.Next() {
		var r reading
		err := rows.Scan(&r.RecordedOn, &r.Reading)
		if err != nil {
			interal_error(c, "Error during type conversion", err)
			return
		}

		// Normalize readings
		if !previousDate.IsZero() {
			diff := r.RecordedOn.Sub(previousDate)
			days := diff.Hours() / 24
			r.Reading = r.Reading / days
			rr.Readings = append(rr.Readings, r)
		}
		previousDate = r.RecordedOn
	}

	// Query for meta information
	query = "SELECT kind_name, unit FROM kinds WHERE kind_id = @kind"
	err = conn.QueryRow(context.Background(), query, args).Scan(&rr.Kind, &rr.Unit)
	if err != nil {
		interal_error(c, "Unable to get kind meta information from database", err)
		return
	}

	c.IndentedJSON(http.StatusOK, rr)
}

func postReading(c *gin.Context) {
	// Establish connection
	conn, err := pgx.Connect(context.Background(), db_url)
	if err != nil {
		interal_error(c, "Unable to connect to database", err)
		return
	}
	defer conn.Close(context.Background())

	// Call BindJSON to bind the received JSON to
	var r reading
	if err := c.BindJSON(&r); err != nil {
		bad_request_err(c, "Failed to parse argument", err)
		return
	}

	// Build query
	query := "INSERT INTO readings (kind_id, recorded_on, reading) VALUES (@kind_id, @recorded_on, @reading)"
	args := pgx.NamedArgs{
		"kind_id":     r.KindId,
		"recorded_on": r.RecordedOn,
		"reading":     r.Reading,
	}

	// Execute query
	_, err = conn.Exec(context.Background(), query, args)
	if err != nil {
		interal_error(c, "Failed to create new reading", err)
		return
	}

	// Build query
	query = "SELECT kind_id, recorded_on, reading FROM readings ORDER BY reading_id DESC LIMIT 1"

	// Query for readings
	err = conn.QueryRow(context.Background(), query).Scan(&r.KindId, &r.RecordedOn, &r.Reading)
	if err != nil {
		interal_error(c, "Unable to get readings from database", err)
		return
	}

	c.IndentedJSON(http.StatusCreated, r)
}

func interal_error(c *gin.Context, msg string, err error) {
	m := fmt.Sprintf(msg+": %v", err)
	fmt.Fprintf(os.Stderr, m+": \n")
	c.JSON(http.StatusInternalServerError, response_error{m})
}

func bad_request(c *gin.Context, msg string) {
	fmt.Fprint(os.Stderr, msg+"\n")
	c.JSON(http.StatusBadRequest, response_error{msg})
}

func bad_request_err(c *gin.Context, msg string, err error) {
	m := fmt.Sprintf(msg+": %v", err)
	fmt.Fprintf(os.Stderr, m+": \n")
	c.JSON(http.StatusBadRequest, response_error{m})
}
