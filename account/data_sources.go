package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type dataSources struct {
	DB          *sqlx.DB
	RedisClient *redis.Client
}

// InitDS establishes connections to fields in dataSources
func initDS() (*dataSources, error) {
	log.Printf("Initializing data sources\n")
	// load env variables - we could pass these in,
	// but this is sort of just a top-level (main package)
	// helper function, so I'll just read them in here
	pgHost := os.Getenv("PG_HOST")
	pgPort := os.Getenv("PG_PORT")
	pgUser := os.Getenv("PG_USER")
	pgPassword := os.Getenv("PG_PASSWORD")
	pgDB := os.Getenv("PG_DB")
	pgSSL := os.Getenv("PG_SSL")

	pgConnString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", pgHost, pgPort, pgUser, pgPassword, pgDB, pgSSL)

	log.Printf("Connecting to Postgresql\n")

	// sleeping 3 seconds because sometimes the db needs a few more seconds after start to accept connections
	time.Sleep(3 * time.Second)

	db, err := sqlx.Open("postgres", pgConnString)
	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	// Verify database connection is working
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}

	// establish redis conection
	rdHost := os.Getenv("REDIS_HOST")
	rdPort := os.Getenv("REDIS_PORT")
	rdAddr := fmt.Sprintf("%s:%s", rdHost, rdPort)

	rdb := redis.NewClient(&redis.Options{
		Addr:     rdAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	if rdb == nil {
		return nil, fmt.Errorf("unknown error when connecting to redis")
	}

	return &dataSources{
		DB:          db,
		RedisClient: rdb,
	}, nil
}

// close all DS-connections to be used in graceful server shutdown
func (d *dataSources) close() error {
	if err := d.DB.Close(); err != nil {
		return fmt.Errorf("error closing postgres: %w", err)
	}

	if err := d.RedisClient.Close(); err != nil {
		return fmt.Errorf("error closing redis: %w", err)
	}

	return nil
}
