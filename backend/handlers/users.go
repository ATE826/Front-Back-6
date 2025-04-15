package handlers

import (
	"front-back_6/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetCurrentUser(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID type"})
		return
	}

	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	user.Password = "" // Не возвращаем пароль

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (s *Server) LogOut(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Успешный выход из системы"})
}
