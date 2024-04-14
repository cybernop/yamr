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
	RecordedOn time.Time `json:"recorded_on"`
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

func interal_error(c *gin.Context, msg string, err error) {
	fmt.Fprintf(os.Stderr, msg+": %v\n", err)
	c.JSON(http.StatusInternalServerError, response_error{msg})
}

func bad_request(c *gin.Context, msg string) {
	fmt.Fprint(os.Stderr, msg+"\n")
	c.JSON(http.StatusBadRequest, response_error{msg})
}

func postReading(c *gin.Context) {
	var newReading reading

	// Call BindJSON to bind the received JSON to
	// newReading.
	if err := c.BindJSON(&newReading); err != nil {
		return
	}

	// Add the new reading to the slice
	readings = append(readings, newReading)
	c.IndentedJSON(http.StatusCreated, newReading)
}
