package api

import (
	"loading_time/internal/app/handler"
	"loading_time/internal/app/repository"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория")
	}

	handler := handler.NewHandler(repo)

	r := gin.Default()
	// добавляем наш html/шаблон
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")
	// слева название папки, в которую выгрузится наша статика
	// справа путь к папке, в которой лежит статика

	r.GET("/home", handler.GetShips)
	r.GET("/ship/:id", handler.GetShip)
	r.GET("/request", handler.GetRequest)
	r.GET("/request/:id", handler.GetRequest)
	r.GET("/request/add/:id", handler.AddToRequest)
	r.GET("/request/:id/remove/:ship_id", handler.RemoveShipFromRequest)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	log.Println("Server down")
}
