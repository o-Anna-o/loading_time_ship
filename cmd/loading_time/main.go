package main

// go run cmd/loading_time/main.go

// для миграций
// go run cmd/loading_time/migrate/main.go

import (
	"fmt"

	"loading_time/internal/app/config"
	"loading_time/internal/app/dsn"
	"loading_time/internal/app/handler"
	"loading_time/internal/app/pkg"
	"loading_time/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	router := gin.Default()
	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	postgresString := dsn.FromEnv()
	fmt.Println(postgresString)

	rep, errRep := repository.New(postgresString)
	if errRep != nil {
		logrus.Fatalf("error initializing repository: %v", errRep)
	}

	hand := handler.NewHandler(rep)

	application := pkg.NewApp(conf, router, hand)
	application.RunApp()
}
