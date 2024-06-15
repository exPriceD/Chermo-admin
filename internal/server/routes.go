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
	r.GET("/events", s.getEventsHandler)
	r.GET("/museums", s.getMuseumsHandler)

	protectedAPI := r.Group("/api/v1")
	protectedAPI.Use(AuthMiddleware())
	{
		protectedAPI.GET("/events", s.getEventsForAdminPanelHandler)
		protectedAPI.POST("/create_user", s.CreateUserHandler)
		protectedAPI.GET("/protected", s.ProtectedEndpoint)
	}

	return r
}

func (s *Server) getMuseumsHandler(c *gin.Context) {
	museums, err := s.museumRepo.GetMuseums()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, museums)
}

func (s *Server) getEventsHandler(c *gin.Context) {
	events, err := s.eventsRepo.GetEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

func (s *Server) getEventsForAdminPanelHandler(c *gin.Context) {
	museumID, exists := c.Get("museum_id")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}
	museumIDInt, ok := museumID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid museum ID type"})
		return
	}

	fmt.Println(museumIDInt)

	events, err := s.eventsRepo.GetEventsByMuseum(museumIDInt)
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
		eventStruct.MuseumID = 1

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
