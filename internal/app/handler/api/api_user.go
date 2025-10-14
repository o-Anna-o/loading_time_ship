package api

import (
	"loading_time/internal/app/ds"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	Repository interface {
		DB() *gorm.DB
	}
}

// RegisterUserAPI - POST /api/users/register - регистрация
func (h *UserHandler) RegisterUserAPI(c *gin.Context) {
	var user ds.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"description": err.Error(),
		})
		return
	}

	db := h.Repository.DB()

	err := db.Create(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":      "error",
			"description": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   user,
	})
}

// GetUserProfileAPI - GET /api/users/profile - поля пользователя
func (h *UserHandler) GetUserProfileAPI(c *gin.Context) {
	const fixedUserID = 1

	var user ds.User

	db := h.Repository.DB()

	err := db.First(&user, fixedUserID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":      "error",
			"description": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   user,
	})
}

// _________________________________________________________________________
// UpdateUserProfileAPI - PUT /api/users/profile - обновление профиля
func (h *UserHandler) UpdateUserProfileAPI(c *gin.Context) {
	const fixedUserID = 1

	var updates struct {
		FIO                 string  `json:"fio"`
		Contacts            string  `json:"contacts"`
		CargoWeight         float64 `json:"cargo_weight"`
		Containers20ftCount int     `json:"containers_20ft_count"`
		Containers40ftCount int     `json:"containers_40ft_count"`
	}

	if err := c.BindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"description": err.Error(),
		})
		return
	}

	db := h.Repository.DB()
	err := db.Model(&ds.User{}).Where("user_id = ?", fixedUserID).Updates(updates).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":      "error",
			"description": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Profile updated successfully",
	})
}

// LoginUserAPI - POST /api/users/login - аутентификация
func (h *UserHandler) LoginUserAPI(c *gin.Context) {
	var credentials struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"description": err.Error(),
		})
		return
	}

	// проверка
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Login successful",
		"data": gin.H{
			"user_id": 1,
			"fio":     "Агапова Анна Денисовна",
		},
	})
}

// LogoutUserAPI - POST /api/users/logout - деавторизация
func (h *UserHandler) LogoutUserAPI(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Logout successful",
	})
}
