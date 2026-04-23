package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)


func main() {
	godotenv.Load()

	ctx := context.Background()

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL not set")
	}
	db, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatal("Unable to connect to database: ", err)
	}
	defer db.Close(ctx)
	fmt.Printf("Successfully connected to Postgres\n")

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
	if err := rdb.Ping(ctx).Err(); err != nil{
		log.Fatal("Unable to connect to Redis: ", err)
	}
	fmt.Printf("Successfully connected to Redis\n")

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "online",
			"db":     "connected",
		})
	})

	fmt.Println("Server Chat-Core starting on port 8080")
	r.Run(":8080")
}
