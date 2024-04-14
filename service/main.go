package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type reading struct {
	ID         int       `json:"id"`
	Kind       string    `json:"kind"`
	RecordedOn time.Time `json:"recordedOn"`
	Reading    float64   `json:"reading"`
}

type response_error struct {
	Error string `json:"error"`
}

var db_url = "postgres://" + os.Getenv("POSTGRES_USER") + ":" + os.Getenv("POSTGRES_PASSWORD") + "@" + os.Getenv("POSTGRES_HOSTNAME") + ":5432/" + os.Getenv("POSTGRES_DB")

func main() {
	router := gin.Default()
	router.GET("/readings", getReadings)
	router.POST("/reading", postReading)

	router.Run("localhost:8080")
}

func getReadings(c *gin.Context) {
	// Establish connection
	conn, err := pgx.Connect(context.Background(), db_url)
	if err != nil {
		interal_error(c, "Unable to connect to database", err)
		return
	}
	defer conn.Close(context.Background())

	// Query the readings
	var rows pgx.Rows
	if k := c.Request.URL.Query()["kind"]; len(k) > 0 {
		// Filter by kind
		if len(k) == 1 {
			rows, err = conn.Query(context.Background(), "SELECT id, kind, recorded_on, reading FROM readings WHERE kind = $1", k[0])
			if err != nil {
				interal_error(c, "Unable to get readings from database", err)
				return
			}
		} else {
			bad_request(c, "Multiple 'kind' parameters")
			return
		}
	} else {
		rows, err = conn.Query(context.Background(), "SELECT id, kind, recorded_on, reading FROM readings")
		if err != nil {
			interal_error(c, "Unable to get readings from database", err)
			return
		}
	}

	// Convert query results to structs
	var rr []reading
	for rows.Next() {
		var r reading
		err := rows.Scan(&r.ID, &r.Kind, &r.RecordedOn, &r.Reading)
		if err != nil {
			interal_error(c, "Error during type conversion", err)
			return
		}

		rr = append(rr, r)
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
	query := "INSERT INTO readings (kind, recorded_on, reading) VALUES (@kind, @recorded_on, @reading)"
	args := pgx.NamedArgs{
		"kind":        r.Kind,
		"recorded_on": r.RecordedOn,
		"reading":     r.Reading,
	}

	// Execute query
	_, err = conn.Exec(context.Background(), query, args)
	if err != nil {
		interal_error(c, "Failed to create new reading", err)
		return
	}

	c.IndentedJSON(http.StatusCreated, nil)
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
