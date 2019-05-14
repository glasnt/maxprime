package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// maxNumber to check for primality - increase for a longer test
	defaultMaxNumber = 500000
	// for purposes of demo, limit the highest number user can pass in
	defaultMaxNumberCeiling = 99999999
)

var (
	// allow to override the max number ceiling
	// getEnvAsInt("MAX_NUMBER_CEILING", defaultMaxNumberCeiling)
	maxNumberCeiling = defaultMaxNumberCeiling
)

func healthHandler(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

func homeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", map[string]interface{}{
		"release": release,
		"max":     defaultMaxNumber,
		"ceiling": maxNumberCeiling,
	})
}

func primeArgHandler(c *gin.Context) {

	maxVar := c.Param("max")
	logger.Printf("max == %s", maxVar)
	if maxVar == "" {
		log.Println("Error on nil max parameter")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Null Argument",
			"status":  "BadRequest",
		})
		return
	}

	maxNum, err := strconv.Atoi(maxVar)
	if err != nil {
		logger.Printf("Error while parsing max parameter: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Argument",
			"status":  "BadRequest",
		})
		return
	}

	if maxNum > maxNumberCeiling {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("Request (%d) over max number: %d", maxNum, maxNumberCeiling),
			"status":  "BadRequest",
		})
		return
	}

	c.JSON(http.StatusOK, getPrimeResponse(maxNum))
}

func defaultPrimeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, getPrimeResponse(defaultMaxNumber))
}

// PrimeResponse represents body of the prime request response
type PrimeResponse struct {
	ID       string `json:"id"`
	Ts       string `json:"ts"`
	Duration string `json:"dur"`
	Release  string `json:"rel"`
	Prime    *prime `json:"prime"`
}

func getPrimeResponse(maxNum int) *PrimeResponse {

	s := time.Now()
	p := calcPrime(maxNum)
	d := time.Since(s)

	resp := &PrimeResponse{
		ID:       newID(),
		Duration: fmt.Sprintf("%s", d%time.Second),
		Prime:    p,
		Ts:       time.Now().UTC().String(),
		Release:  release,
	}

	return resp
}

func newID() string {
	id, err := uuid.NewRandom()
	if err != nil {
		logger.Fatalf("Error while getting id: %v\n", err)
	}
	return id.String()
}

func getEnv(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	if fallbackValue == "" {
		logger.Fatalf("Required env var (%s) not set", key)
	}
	return fallbackValue
}
