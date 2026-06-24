package main

import (
	"log"
	"os"

	"time"

	"strconv"

	"github.com/Malachy-Olua/social-platform/internal/db"
	"github.com/Malachy-Olua/social-platform/internal/store"

	"github.com/joho/godotenv"
)

const version = "0.0.1"

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	maxOpenConns, err := strconv.Atoi(os.Getenv("MAX_OPEN_CONNS"))
	if err != nil {
		maxOpenConns = 30 // sensible default
	}

	maxIdleConns, err := strconv.Atoi(os.Getenv("MAX_IDLE_CONNS"))
	if err != nil {
		maxIdleConns = 30 // sensible default
	}

	maxIdleTime, err := strconv.Atoi(os.Getenv("MAX_IDLE_TIME"))
	if err != nil {
		maxIdleTime = 15 // sensible default
	}

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN is not set")
	}

	cfg := config{
		addr: ":9090",
		db: dbConfig{
			dsn:          dsn,
			maxOpenConns: maxOpenConns,
			maxIdleConns: maxIdleConns,
			maxIdleTime:  time.Duration(maxIdleTime) * time.Minute,
		},
		env: os.Getenv("ENV"),
	}

	db, err := db.New(
		cfg.db.dsn,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)

	if err != nil {
		log.Panic(err)
	}

	defer db.Close()
	log.Println("Connected to database...")

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
	}

	// app.storage = store

	mux := app.mount()
	log.Fatal(app.run(mux))
}
