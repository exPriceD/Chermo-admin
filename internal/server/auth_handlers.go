package server

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/exPriceD/Chermo-admin/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func (s *Server) LoginHandler(c *gin.Context) {
	var creds models.Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	receivedUser, err := s.authRepo.GetUser(creds.Username)
	err = bcrypt.CompareHashAndPassword([]byte(receivedUser.Password), []byte(creds.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	expTime, _ := strconv.Atoi(os.Getenv("TOKEN_EXPIRATION_TIME"))
	expirationTime := time.Now().Add(time.Duration(expTime) * time.Minute)

	claims := &models.Claims{
		Username: creds.Username,
		Role:     receivedUser.Role,
		MuseumID: receivedUser.MuseumID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func (s *Server) CreateUserHandler(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists || role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var newUser struct {
		Username   string `json:"username"`
		Password   string `json:"password"`
		Role       string `json:"role"`
		MuseumName string `json:"museum_name"`
	}

	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	museum, err := s.museumRepo.GetMuseumByName(newUser.MuseumName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprint("GetMuseumByName error ", err)})
		return
	}

	user := models.User{
		Username: newUser.Username,
		Password: string(hashedPassword),
		Role:     newUser.Role,
		MuseumID: museum.ID,
	}

	if err := s.authRepo.InsertUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprint("Could not create user ", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}
