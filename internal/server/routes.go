package server

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/exPriceD/Chermo-admin/internal/entities"
	"github.com/exPriceD/Chermo-admin/internal/models"
	"github.com/exPriceD/Chermo-admin/internal/parser"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
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
		protectedAPI.GET("/event/:id", s.getEventByIDHandler)
		protectedAPI.POST("/create_user", s.CreateUserHandler)

		protectedAPI.GET("/event/:id/schedule", s.GetEventScheduleHandler)
		protectedAPI.POST("/event/:id/schedule", s.CreateScheduleHandler)
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

func (s *Server) getEventByIDHandler(c *gin.Context) {
	eventID := c.Param("id")
	eventIDInt, err := strconv.Atoi(eventID)

	userMuseumID, exists := c.Get("museum_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	event, err := s.eventsRepo.GetEventByID(eventIDInt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve event"})
		}
		return
	}

	if event.MuseumID != userMuseumID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have access to this event"})
		return
	}

	c.JSON(http.StatusOK, event)
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

func (s *Server) GetEventScheduleHandler(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	schedules, err := s.scheduleRepo.GetScheduleByEventID(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve schedule"})
		return
	}

	c.JSON(http.StatusOK, schedules)
}

func (s *Server) CreateScheduleHandler(c *gin.Context) {
	var req models.ScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startDate, err := time.Parse("02.01.2006", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
		return
	}

	endDate, err := time.Parse("02.01.2006", req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
		return
	}

	err = s.scheduleRepo.CreateSchedule(startDate, endDate, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Could not insert in db %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Schedule created successfully"})
}
