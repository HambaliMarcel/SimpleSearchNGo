package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type SearchResult struct {
	Results []string `json:"results"`
	Exists  bool     `json:"exists"`
}

func main() {
	r := gin.Default()

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	defer client.Close()

	r.Use(cors.Default())

	r.POST("/search", func(c *gin.Context) {
		var requestBody map[string]string
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		query := requestBody["query"]

		// Convert query to lowercase for case-insensitive search
		query = strings.ToLower(query)

		// Check if the key exists in Redis
		exists, err := client.Exists(ctx, query).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if exists == 0 {
			// Key doesn't exist, return appropriate response
			c.JSON(http.StatusOK, SearchResult{
				Exists:  false,
				Results: nil,
			})
			return
		}

		// Key exists, retrieve the result
		result, err := client.Get(ctx, query).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		fmt.Println("Raw result from Redis:", result)

		// Check if the result is valid JSON
		var searchResult SearchResult
		if err := json.Unmarshal([]byte(result), &searchResult); err != nil {
			// Handle the case where the result is not valid JSON
			c.JSON(http.StatusOK, SearchResult{
				Exists:  true,
				Results: []string{result},
			})
			return
		}

		c.JSON(http.StatusOK, searchResult)

	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	address := fmt.Sprintf(":%s", port)
	r.Run(address)
}
