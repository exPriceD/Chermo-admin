package server

import (
	"fmt"
	"github.com/exPriceD/Chermo-admin/internal/repositories/events"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/exPriceD/Chermo-admin/internal/database"
)

type Server struct {
	port       int
	eventsRepo *events.Repository
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
	Srv.eventsRepo = eventsRepo
	return server
}
