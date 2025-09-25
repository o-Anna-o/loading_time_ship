package handler

import (
	"context"
	"loading_time/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository  *repository.Repository
	MinioClient *minio.Client
}

func NewHandler(r *repository.Repository) *Handler {
	// Инициализация MinIO
	endpoint := "localhost:9000"
	accessKey := "minio_login_001"
	secretKey := "minio_login_001"
	useSSL := false

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		logrus.Errorf("MinIO connection error: %v", err)
		return &Handler{
			Repository:  r,
			MinioClient: nil,
		}
	}

	logrus.Info("MinIO connected successfully")

	// проверка подключения к бакету loading-time-img
	bucketName := "loading-time-img"
	exists, err := client.BucketExists(context.Background(), bucketName)
	if err != nil {
		logrus.Errorf("Bucket check error: %v", err)
	} else if exists {
		logrus.Infof("Bucket %s is available", bucketName)
	} else {
		logrus.Warnf("Bucket %s does not exist", bucketName)
	}

	return &Handler{
		Repository:  r,
		MinioClient: client,
	}
}

func (h *Handler) GetShips(ctx *gin.Context) {
	var ships []repository.Ship
	var err error

	nameQuery := ctx.Query("name_query")

	if nameQuery == "" {
		ships, err = h.Repository.GetShips()
	} else {
		ships, err = h.Repository.GetShipsByName(nameQuery)
	}

	if err != nil {
		logrus.Error(err)
	}
	// для подсчета кораблей в заявке
	request, err := h.Repository.GetRequest(1)
	requestCount := 0
	if err == nil {
		for _, shipInRequest := range request.Ships {
			requestCount += shipInRequest.Count
		}
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"ships":         ships,
		"name_query":    nameQuery,
		"request_count": requestCount, // количество кораблей в заявке
	})
}

func (h *Handler) GetShip(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	ship, err := h.Repository.GetShip(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "ship.html", gin.H{
		"ship": ship,
	})
}

func (h *Handler) GetRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")

	if idStr == "" {
		ctx.HTML(http.StatusOK, "request.html", gin.H{
			"request": repository.Request{
				ID:    0,
				Ships: []repository.ShipInRequest{},
			},
		})
		return
	}

	requestID, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Invalid request ID")
		return
	}

	request, err := h.Repository.GetRequest(requestID)
	if err != nil {
		ctx.String(http.StatusNotFound, err.Error())
		return
	}

	ctx.HTML(http.StatusOK, "request.html", gin.H{
		"request": request,
	})
}
