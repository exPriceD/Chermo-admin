package server

import (
	"github.com/exPriceD/Chermo-admin/internal/parser"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/", s.HelloWorldHandler)
	r.GET("/health", s.healthHandler)
	r.POST("/login", s.login)

	protectedAPI := r.Group("/api")
	protectedAPI.Use(AuthMiddleware())
	{
		protectedAPI.POST("/create_user", s.CreateUser)
		protectedAPI.GET("/protected", s.ProtectedEndpoint)
		protectedAPI.GET("/events", s.getEventsHandler)
	}

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}

func (s *Server) getEventsHandler(c *gin.Context) {
	items, err := parser.FetchEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (s *Server) ProtectedEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the protected endpoint!"})
}
