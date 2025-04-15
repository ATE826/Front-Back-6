package handlers

import (
	"front-back_6/models"
	"front-back_6/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterInput struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Server struct {
	db    *gorm.DB
	cache *cache.Cache
}

func NewServer(db *gorm.DB) *Server {
	return &Server{db: db}
}

func (s *Server) Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminLogin := os.Getenv("ADMIN_LOGIN")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	var role string
	if input.Email == adminLogin && input.Password == adminPassword {
		role = "admin"
	} else {
		role = "user"
	}

	user := models.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  input.Password,
		Role:      role,
	}

	user.HashPassword()

	if err := s.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created"})
}

func (s *Server) LoginCheck(email, password string) (string, error) {
	var err error

	user := models.User{}

	if err = s.db.Model(models.User{}).Where("email = ?", email).Take(&user).Error; err != nil {
		return "", err
	}

	err = user.VerifyPassword(password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	token, err := utils.GenerateToken(user)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Server) Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := s.db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	if err := user.VerifyPassword(input.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	token, err := utils.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// Устанавливаем куку
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "auth_token",         // имя куки
		Value:    token,                // значение куки
		Path:     "/",                  // путь, на котором кука будет доступна
		HttpOnly: true,                 // кука доступна только через HTTP(S) протокол
		SameSite: http.SameSiteLaxMode, // кука доступна только на том же сайте
		Secure:   false,                // ставь true если HTTPS
		MaxAge:   3600,                 // время жизни куки в секундах
	})

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func (s *Server) GetData(c *gin.Context) {
	// Проверка наличия данных в кэше
	data, found := s.cache.Get("user_data")
	if found {
		// Если данные есть в кэше, отправляем их
		c.JSON(http.StatusOK, gin.H{"data": data})
		return
	}

	// Если данных нет в кэше, получаем их из базы данных
	// Пример: данные, которые могут быть кэшированы
	var users []models.User
	if err := s.db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve data"})
		return
	}

	// Генерация данных для отправки
	// Здесь может быть любая логика для обработки данных перед отправкой
	dataToCache := gin.H{"users": users}

	// Сохраняем данные в кэш на 1 минуту
	s.cache.Set("user_data", dataToCache, cache.DefaultExpiration)

	// Отправляем данные пользователю
	c.JSON(http.StatusOK, gin.H{"data": dataToCache})
}
