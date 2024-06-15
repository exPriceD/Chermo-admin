package server

import (
	"fmt"
	"github.com/exPriceD/Chermo-admin/internal/repositories/auth"
	"github.com/exPriceD/Chermo-admin/internal/repositories/events"
	"github.com/exPriceD/Chermo-admin/internal/repositories/museum"
	"github.com/exPriceD/Chermo-admin/internal/repositories/schedule"
	"github.com/exPriceD/Chermo-admin/internal/repositories/visitors"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/exPriceD/Chermo-admin/internal/database"
)

type Server struct {
	port         int
	eventsRepo   *events.Repository
	authRepo     *auth.Repository
	museumRepo   *museum.Repository
	visitorsRepo *visitors.Repository
	scheduleRepo *schedule.Repository
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	Srv := &Server{
		port: port,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", Srv.port),
		Handler:      Srv.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	db, err := database.InitDB()
	if err != nil {
		panic(fmt.Sprintf("cannot connect to database: %s", err))
	}

	eventsRepo := events.NewEventsRepository(db)
	authRepo := auth.NewAuthRepository(db)
	museumRepo := museum.NewMuseumRepository(db)
	visitorsRepo := visitors.NewVisitorsRepository(db)
	scheduleRepo := schedule.NewScheduleRepository(db)

	Srv.eventsRepo = eventsRepo
	Srv.authRepo = authRepo
	Srv.museumRepo = museumRepo
	Srv.visitorsRepo = visitorsRepo
	Srv.scheduleRepo = scheduleRepo
	return server
}
