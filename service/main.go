package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type reading struct {
	ID    string  `json:"id"`
	Type  string  `json:"type"`
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

var readings = []reading{
	{ID: "1", Type: "gas", Date: "2023-01-01", Value: 1234.56},
	{ID: "2", Type: "gas", Date: "2023-01-02", Value: 1235.67},
	{ID: "3", Type: "energy", Date: "2023-01-01", Value: 456.8},
}

func main() {
	router := gin.Default()
	router.GET("/readings", getReadings)
	router.POST("/reading", postReading)

	router.Run("localhost:8080")
}

func getReadings(c *gin.Context) {
	q := c.Request.URL.Query()
	if t := q["type"]; len(t) == 1 {
		type_ := t[0]
		var ret []reading

		for _, r := range readings {
			if r.Type == type_ {
				ret = append(ret, r)
			}
		}
		c.IndentedJSON(http.StatusOK, ret)
		return
	}

	c.IndentedJSON(http.StatusOK, readings)
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
