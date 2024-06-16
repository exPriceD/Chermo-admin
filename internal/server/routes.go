package server

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/exPriceD/Chermo-admin/internal/entities"
	"github.com/exPriceD/Chermo-admin/internal/mail"
	"github.com/exPriceD/Chermo-admin/internal/models"
	"github.com/exPriceD/Chermo-admin/internal/parser"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.POST("/login", s.LoginHandler)

	r.GET("/museums", s.getMuseumsHandler)

	r.GET("/events", s.getEventsHandler)
	r.GET("/event/:id", s.getEventWithSchedule)
	r.POST("/event/:id/register", s.RegisterVisitorHandler)

	r.GET("/events/register/confirm/:id/:time_slot_id", s.confirmRegistration)
	r.GET("/events/register/cancel/:id/:time_slot_id", s.cancelRegistration)

	protectedAPI := r.Group("/api/v1")
	protectedAPI.Use(AuthMiddleware())
	{
		protectedAPI.GET("/events", s.getEventsForAdminPanelHandler)
		protectedAPI.GET("/event/:id", s.getEventByIDHandler)

		protectedAPI.GET("/event/:id/schedule", s.GetEventScheduleHandler)
		protectedAPI.POST("/event/:id/schedule", s.CreateScheduleHandler)

		protectedAPI.GET("/event/:id/visitors", s.getVisitorsHandler)
		protectedAPI.GET("/event/:id/visitors/date", s.getVisitorsDatesHandler)
		protectedAPI.GET("/event/:id/visitors/date/time", s.getVisitorsTimesHandler)

		protectedAPI.POST("/create_user", s.CreateUserHandler)
	}

	return r
}

func (s *Server) getVisitorsHandler(c *gin.Context) {
	eventID := c.Param("id")
	eventIDInt, err := strconv.Atoi(eventID)

	dates, err := s.scheduleRepo.GetEventDates(eventIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dates)
}

func (s *Server) getVisitorsDatesHandler(c *gin.Context) {
	eventID := c.Param("id")
	eventDate := c.Query("date")
	fmt.Println(eventDate)
	eventIDInt, _ := strconv.Atoi(eventID)

	timeSlots, err := s.scheduleRepo.GetTimeSlots(eventIDInt, eventDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, timeSlots)
}

func (s *Server) getVisitorsTimesHandler(c *gin.Context) {
	eventID := c.Param("id")
	eventIDInt, _ := strconv.Atoi(eventID)

	eventDate := c.Query("date")
	startTime := c.Query("time")
	visitors, err := s.visitorsRepo.GetVisitors(eventIDInt, eventDate, startTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, visitors)
}

func (s *Server) confirmRegistration(c *gin.Context) {
	visitorID := c.Param("id")
	timeSlotID := c.Param("time_slot_id")

	timeSlotIDInt, _ := strconv.Atoi(timeSlotID)

	visitorIDInt, err := strconv.Atoi(visitorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid visitor ID"})
		return
	}

	exists, err := s.visitorsRepo.VisitorExists(visitorIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve visitor"})
		return
	}
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Visitor not found"})
		return
	}

	exists, err = s.scheduleRepo.TimeSlotExists(timeSlotIDInt)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Timeslot not found"})
		return
	}

	err = s.scheduleRepo.UpdateRegistrationStatus(visitorIDInt, timeSlotIDInt, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not confirm registration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Запись подтверждена"})
}

func (s *Server) cancelRegistration(c *gin.Context) {
	visitorID := c.Param("id")
	timeSlotID := c.Param("time_slot_id")

	timeSlotIDInt, _ := strconv.Atoi(timeSlotID)

	visitorIDInt, err := strconv.Atoi(visitorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid visitor ID"})
		return
	}

	exists, err := s.visitorsRepo.VisitorExists(visitorIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve visitor"})
		return
	}
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Visitor not found"})
		return
	}

	exists, err = s.scheduleRepo.TimeSlotExists(timeSlotIDInt)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Timeslot not found"})
		return
	}

	err = s.scheduleRepo.UpdateRegistrationStatus(visitorIDInt, timeSlotIDInt, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not cancel registration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Запись отменена"})
}

func (s *Server) getEventWithSchedule(c *gin.Context) {
	eventID := c.Param("id")
	eventIDInt, err := strconv.Atoi(eventID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	event, err := s.eventsRepo.GetEventByID(eventIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve event"})
		return
	}

	schedules, err := s.scheduleRepo.GetScheduleByEventID(eventIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve schedule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": event, "slots": schedules})

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
	var request models.ScheduleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startDate, err := time.Parse("02.01.2006", request.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
		return
	}

	endDate, err := time.Parse("02.01.2006", request.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
		return
	}

	err = s.scheduleRepo.CreateSchedule(startDate, endDate, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Could not insert in db %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Schedule created successfully"})
}

func (s *Server) RegisterVisitorHandler(c *gin.Context) {
	var request entities.RegistrationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isAvailable, err := s.scheduleRepo.IsTimeSlotAvailable(request.TimeslotID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not check time slot availability"})
		return
	}

	if !isAvailable {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Time slot is not available"})
		return
	}

	visitor := entities.Visitor{
		FirstName:  request.FirstName,
		LastName:   request.LastName,
		Patronymic: request.Patronymic,
		Phone:      request.Phone,
		Email:      request.Email,
	}

	visitorID, err := s.visitorsRepo.AddVisitor(visitor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Could not add visitor %v", err)})
		return
	}

	err = s.scheduleRepo.RegisterVisitor(request.TimeslotID, visitorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Could not register visitor %v", err)})
		return
	}

	// Получаем информацию о мероприятии
	event, err := s.eventsRepo.GetEventByTimeslotID(request.TimeslotID)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve event information"})
		return
	}

	// Формируем текст письма
	mailText := fmt.Sprintf(
		"Для подтверждения регистрации на мероприятие перейдите по ссылке: http://localhost:8080/events/register/confirm/%d/%d\n\n"+
			"Информация о мероприятии:\n"+
			"Название: %s\n"+
			"Место: %s\n"+
			"Дата: %s\n"+
			"Время: %s - %s\n",
		visitorID, request.TimeslotID, event.Title, event.Museum, event.Date.Format("02.01.2006"), event.StartTime.Format("15:04"), event.EndTime.Format("15:04"),
	)

	err = mail.SendEmail(os.Getenv("MAIL_USERNAME"), visitor.Email, "Подтверждение регистрации на мероприятие", mailText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not send confirmation email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}
