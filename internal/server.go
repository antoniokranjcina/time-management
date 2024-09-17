package internal

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"net/http"
	"os"
	"strconv"
	"time"
	"time-management/internal/location/infrastructure/repository"
	http2 "time-management/internal/location/interface/http"
)

type Server struct {
	port int
	db   *sql.DB
}

func NewServer() *http.Server {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}

	db, err := InitializeDB()
	if err != nil {
		panic(err)
	}

	defer CloseDB()

	// Initialize repositories
	locationRepository := repository.NewPgLocationRepository(db)

	// Initialize handlers
	locationHandler := http2.NewLocationHandler(locationRepository)

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      SetupRoutes(locationHandler),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
