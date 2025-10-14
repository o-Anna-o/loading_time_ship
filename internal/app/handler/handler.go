package handler

import (
	"context"
	"loading_time/internal/app/repository"

	"loading_time/internal/app/handler/api"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository            *repository.Repository
	ShipAPIHandler        *api.ShipHandler
	RequestShipAPIHandler *api.RequestShipHandler
	UserAPIHandler        *api.UserHandler
	MinioClient           *minio.Client
}

func NewHandler(rep *repository.Repository) *Handler {
	// Инициализация MinIO клиента
	endpoint := "localhost:9000"
	accessKey := "minio_login_001"
	secretKey := "minio_login_001"
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		logrus.Errorf("MinIO connection error: %v", err)
		// Продолжаем работу даже если MinIO не доступен
		minioClient = nil
	} else {
		logrus.Info("MinIO connected successfully")
	}

	// Проверяем существование бакета
	if minioClient != nil {
		bucketName := "loading-time-img"
		exists, err := minioClient.BucketExists(context.Background(), bucketName)
		if err != nil {
			logrus.Errorf("Bucket check error: %v", err)
		} else if exists {
			logrus.Infof("Bucket %s is available", bucketName)
		} else {
			logrus.Warnf("Bucket %s does not exist", bucketName)
		}
	}

	// Создаем API хендлеры и передаем MinIO клиент
	shipAPIHandler := &api.ShipHandler{
		Repository:  rep,
		MinioClient: minioClient,
	}

	requestShipAPIHandler := &api.RequestShipHandler{Repository: rep}
	userAPIHandler := &api.UserHandler{Repository: rep}

	return &Handler{
		Repository:            rep,
		ShipAPIHandler:        shipAPIHandler,
		RequestShipAPIHandler: requestShipAPIHandler,
		UserAPIHandler:        userAPIHandler,
		MinioClient:           minioClient,
	}
}
func (h *Handler) SetupRoutes(router *gin.Engine) {
	router.GET("/ships", h.GetShips)
	router.GET("/ship/:id", h.GetShip)
	router.GET("/request_ship", h.CreateOrRedirectRequestShip)
	router.GET("/request_ship/:id", h.GetRequestShip)
	router.POST("/request_ship/add/:ship_id", h.AddShipToRequestShip)
	router.POST("/request_ship/calculate_loading_time/:id", h.CalculateLoadingTime)

	// API маршруты
	apiGroup := router.Group("/api")
	{
		// Домен услуги (контейнеровозы)
		apiGroup.GET("/ships", h.ShipAPIHandler.GetShipsAPI)
		apiGroup.GET("/ships/:id", h.ShipAPIHandler.GetShipAPI)
		apiGroup.POST("/ships", h.ShipAPIHandler.CreateShipAPI)
		apiGroup.PUT("/ships/:id", h.ShipAPIHandler.UpdateShipAPI)
		apiGroup.POST("/ships/:id/image", h.ShipAPIHandler.AddShipImageAPI)
		apiGroup.DELETE("/ships/:id", h.ShipAPIHandler.DeleteShipAPI)
		apiGroup.POST("/ships/:id/add-to-ship-bucket", h.ShipAPIHandler.AddShipToRequestShipAPI)

		apiGroup.PUT("/request_ship/:id/ships/:ship_id", h.RequestShipAPIHandler.UpdateShipInRequestAPI)
		apiGroup.POST("/request_ship/:id/ships/:ship_id", h.RequestShipAPIHandler.DeleteShipFromRequestShipAPI)
		apiGroup.DELETE("/request_ship/:id/ships/:ship_id", h.RequestShipAPIHandler.DeleteShipFromRequestShipAPI)
		apiGroup.PUT("/request_ship/:id/formation", h.RequestShipAPIHandler.FormRequestShipAPI)
		apiGroup.POST("/request_ship/:id/completion", h.RequestShipAPIHandler.CompleteRequestShipAPI)

		apiGroup.GET("/request_ship/:id", h.RequestShipAPIHandler.GetRequestShipAPI)
		apiGroup.PUT("/request_ship/:id", h.RequestShipAPIHandler.UpdateRequestShipAPI)
		apiGroup.POST("/request_ship/:id", h.RequestShipAPIHandler.DeleteRequestShipAPI)
		apiGroup.DELETE("/request_ship/:id", h.RequestShipAPIHandler.DeleteRequestShipAPI)
		apiGroup.GET("/request_ship/basket", h.RequestShipAPIHandler.GetRequestShipBasketAPI)
		apiGroup.GET("/request_ship", h.RequestShipAPIHandler.GetRequestShipsAPI)

		apiGroup.POST("/users/register", h.UserAPIHandler.RegisterUserAPI)
		apiGroup.GET("/users/profile", h.UserAPIHandler.GetUserProfileAPI)
		apiGroup.PUT("/users/profile", h.UserAPIHandler.UpdateUserProfileAPI)
		apiGroup.POST("/users/login", h.UserAPIHandler.LoginUserAPI)
		apiGroup.POST("/users/logout", h.UserAPIHandler.LogoutUserAPI)
	}
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.Static("/styles", "./resources/styles")
	router.Static("/img", "./resources/img")
}

func (h *Handler) errorHandler(c *gin.Context, code int, err error) {
	logrus.Error(err.Error())
	c.JSON(code, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
