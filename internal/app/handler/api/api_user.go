package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"loading_time/internal/app/ds"
	"loading_time/internal/app/repository"
	"loading_time/internal/app/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// =========================================================
// 🔐 JWT + роли
// =========================================================

var jwtKey = []byte("super_secret_key")

type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Middleware проверки токена и роли
func (h *UserHandler) AuthMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется авторизация"})
			c.Abort()
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный токен"})
			c.Abort()
			return
		}

		// Проверка роли
		allowed := false
		for _, role := range allowedRoles {
			if claims.Role == role {
				allowed = true
				break
			}
		}
		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Недостаточно прав"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// =========================================================
// 👤 USERS (регистрация / вход / профиль)
// =========================================================

type UserHandler struct {
	Repository *repository.Repository
}

// @Summary      Регистрация пользователя
// @Description  Создаёт нового пользователя
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      ds.User  true  "Пользователь"
// @Success      201   {object}  ds.User
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /api/users/register [post]
func (h *UserHandler) RegisterUserAPI(c *gin.Context) {
	var user ds.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if user.Role == "" {
		user.Role = "creator"
	}

	// Не хешируем здесь пароль — это делает repository.CreateUser
	if err := h.Repository.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.Password = "" // не отдаём пароль
	c.JSON(http.StatusCreated, user)
}

// @Summary      Вход пользователя
// @Description  Аутентификация и выдача JWT
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        credentials  body      object{login=string,password=string}  true  "Логин и пароль"
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]string
// @Router       /api/users/login [post]
func (h *UserHandler) LoginUserAPI(c *gin.Context) {
	//для поиска ошибки авторизации
	bodyBytes, _ := io.ReadAll(c.Request.Body)
	logrus.Infof("LoginUserAPI: Content-Type=%s, BodyRaw=%s", c.Request.Header.Get("Content-Type"), string(bodyBytes))
	c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	//
	var cred struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&cred); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Repository.GetUserByLogin(cred.Login)
	if err != nil {
		logrus.Infof("LoginUserAPI: пользователь %s не найден, ошибка: %v", cred.Login, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные данные"})
		return
	}

	// Лог перед bcrypt
	logrus.Infof("LoginUserAPI DEBUG: RAW='%s' LEN=%d | HASH='%s' LEN=%d",
		cred.Password, len(cred.Password),
		user.Password, len(user.Password))
	logrus.Infof("LoginUserAPI: найден пользователь %s с хешем %s, введённый пароль: %s", user.Login, user.Password, cred.Password)
	// debug: показать что реально сравниваем
	fmt.Printf("RAW: '%s' | HASH: '%s'\n", cred.Password, user.Password)

	// Проверка пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(cred.Password)); err != nil {
		logrus.Infof("LoginUserAPI: пароль не совпадает для пользователя %s", cred.Login)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные данные"})
		return
	}

	// Генерация JWT
	tokenString, err := utils.GenerateJWT(user.UserID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании токена"})
		return
	}

	// Сохраняем сессию в Redis на 2 часа
	_ = utils.SetSession(tokenString, user.UserID, user.Role, 2*time.Hour)

	c.JSON(http.StatusOK, gin.H{
		"message": "Успешный вход",
		"token":   tokenString,
		"role":    user.Role,
	})
}

// @Summary      Выход пользователя
// @Description  Удаление JWT и сессии
// @Tags         users
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /api/users/logout [post]
func (h *UserHandler) LogoutUserAPI(c *gin.Context) {
	token, _ := c.Cookie("jwt")
	if token != "" {
		_ = utils.DeleteSession(token)
	}

	c.SetCookie("jwt", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// @Summary      Профиль пользователя
// @Description  Получение данных профиля авторизованного пользователя
// @Tags         users
// @Produce      json
// @Success      200  {object}  ds.User
// @Failure      404  {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/users/profile [get]
func (h *UserHandler) GetUserProfileAPI(c *gin.Context) {
	userID := c.GetInt("user_id")
	user, err := h.Repository.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}
	user.Password = "" // не отдаём пароль
	c.JSON(http.StatusOK, user)
}

// UpdateUserProfileAPI

// @Summary      Обновление профиля пользователя
// @Description  Обновляет данные авторизованного пользователя
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      ds.User  true  "Пользователь"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/users/profile [put]
func (h *UserHandler) UpdateUserProfileAPI(c *gin.Context) {
	uid, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var user ds.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	switch v := uid.(type) {
	case float64:
		user.UserID = int(v)
	case int:
		user.UserID = v
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id type"})
		return
	}
	if err := h.Repository.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated"})
}
