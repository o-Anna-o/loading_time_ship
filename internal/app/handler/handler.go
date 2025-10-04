package handler

import (
	"context"
	"loading_time/internal/app/ds"
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

// RegisterHandler регистрирует маршруты
func (h *Handler) RegisterHandler(router *gin.Engine) {
	// GET маршруты
	router.GET("/ships", h.GetShips)
	router.GET("/ship/:id", h.GetShip)

	router.GET("/request_ship", func(ctx *gin.Context) {
		request_ship, err := h.Repository.GetOrCreateUserDraft(1)
		if err != nil {
			logrus.Error(err)
			ctx.HTML(http.StatusInternalServerError, "request_ship.html", gin.H{
				"request_ship": ds.RequestShip{},
				"error":        "Не удалось создать черновик",
			})
			return
		}
		ctx.Redirect(http.StatusFound, "/request_ship/"+strconv.Itoa(request_ship.RequestShipID))
	})

	router.GET("/request_ship/:id", h.GetRequestShip)

	// POST маршруты
	router.POST("/request_ship/add/:ship_id", h.AddShipToRequestShip)
	router.POST("/request_ship/delete/:id", h.DeleteRequestShip)
	router.POST("/request_ship/:id/remove/:ship_id", h.RemoveShipFromRequestShip)
}

// RegisterStatic регистрирует статику
func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/styles", "./resources/styles")
	router.Static("/img", "./resources/img")
}

// errorHandler
func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
