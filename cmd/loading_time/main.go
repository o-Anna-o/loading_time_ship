package main

// docker-compose up -d
// go run cmd/loading_time/main.go

// swag init -g cmd/loading_time/main.go
// go run cmd/loading_time/main.go
// и перейти на http://localhost:8080/swagger/index.html

import (
	"fmt"
	"time"

	"loading_time/internal/app/config"
	"loading_time/internal/app/dsn"
	"loading_time/internal/app/handler"
	"loading_time/internal/app/pkg"
	"loading_time/internal/app/repository"
	"loading_time/internal/app/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	_ "loading_time/docs" // Swagger docs

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-contrib/cors"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите: Bearer {ваш_JWT_токен}
func main() {
	utils.InitRedis()
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		logrus.Infof("Incoming request: %s %s", c.Request.Method, c.Request.URL.Path)
	})

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://192.168.1.68:3000",
			"https://192.168.1.68:3000",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.LoadHTMLGlob("templates/*.html")

	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	postgresString := dsn.FromEnv()
	fmt.Println(postgresString)

	rep, errRep := repository.New(postgresString, conf.RedisEndpoint, conf.RedisPassword, conf.JwtKey)
	if errRep != nil {
		logrus.Fatalf("error initializing repository: %v", errRep)
	}

	http.DefaultClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
		},
	}

	hand := handler.NewHandler(rep)

	router.Use(func(c *gin.Context) {
		if m := c.PostForm("_method"); m != "" {
			logrus.Infof("Overriding method to %s for %s", m, c.Request.URL.Path)
			c.Request.Method = m
		}
		c.Next()
	})

	// Добавляем маршрут для Swagger с логом
	router.GET("/swagger/*any", func(c *gin.Context) {
		logrus.Infof("Serving Swagger UI for path: %s", c.Request.URL.Path)
		ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
	})

	application := pkg.NewApp(conf, router, hand)
	application.RunApp()
}
