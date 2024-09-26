package internal

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"net/http"
	"os"
	"strconv"
	"time"
	empRepo "time-management/internal/employee/infrastructure/repository"
	empHttp "time-management/internal/employee/interface/http"
	locRepo "time-management/internal/location/infrastructure/repository"
	locHttp "time-management/internal/location/interface/http"
	repRepo "time-management/internal/report/infrastructure/repository"
	repHttp "time-management/internal/report/interface/http"
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
	locationRepository := locRepo.NewPgLocationRepository(db)
	employeeRepository := empRepo.NewPgEmployeeRepository(db)
	reportRepository := repRepo.NewPgReportRepository(db)

	// Initialize handlers
	locationHandler := locHttp.NewLocationHandler(locationRepository)
	employeeHandler := empHttp.NewEmployeeHandler(employeeRepository)
	reportHandler := repHttp.NewReportHandler(reportRepository)

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      SetupRoutes(locationHandler, employeeHandler, reportHandler),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
