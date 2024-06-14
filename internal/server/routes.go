package server

import (
	"fmt"
	"github.com/exPriceD/Chermo-admin/internal/entities"
	"github.com/exPriceD/Chermo-admin/internal/models"
	"github.com/exPriceD/Chermo-admin/internal/parser"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.POST("/login", s.LoginHandler)
	// r.GET("/api/get-events", s.getEventsFromChermoHandler)
	r.GET("/api/events", s.getEventsHandler)
	protectedAPI := r.Group("/api/v1")
	protectedAPI.Use(AuthMiddleware())
	{
		protectedAPI.POST("/create_user", s.CreateUserHandler)
		protectedAPI.GET("/protected", s.ProtectedEndpoint)
	}

	return r
}

func (s *Server) getEventsHandler(c *gin.Context) {
	events, err := s.eventsRepo.GetEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

func (s *Server) getEventsFromChermoHandler(c *gin.Context) {
	events, err := parser.FetchEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, event := range events {
		var eventStruct models.Event

		eventStruct.Title = event.Title
		eventStruct.Description = event.Description
		eventStruct.ImageURL = event.Enclosure.URL

		fmt.Println(entities.Event(eventStruct))

		err := s.eventsRepo.InsertEvent(entities.Event(eventStruct))
		if err != nil {
			log.Printf("cannot insert event: %v", err)
			continue
		}
	}
	c.JSON(http.StatusOK, events)
}

func (s *Server) ProtectedEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the protected endpoint!"})
}
