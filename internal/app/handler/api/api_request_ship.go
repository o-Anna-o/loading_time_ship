package api

import (
	"bytes"
	"fmt"
	"io"
	"loading_time/internal/app/ds"
	"loading_time/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RequestShipHandler struct {
	Repository *repository.Repository
}

// / GetRequestShipBasketAPI - GET /api/request_ship/basket - иконка корзины
//
// @Summary Получить корзину запросов
// @Description Retrieve the count of ships in the user's draft request
// @Tags request_ships
// @Produce json
// @Success 200 {object} object "data: {request_ship_id: int, ships_count: int}"
// @Failure 401 {object} object "message: string"
// @Failure 500 {object} object "error: string"
// @Router /api/request_ship/basket [get]
func (h *RequestShipHandler) GetRequestShipBasketAPI(c *gin.Context) {
	logrus.Info("Вхожу в GetRequestShipBasketAPI, если ошибка, то в логрус basket-4ifd;")
	logrus.Infof("GetRequestShipBasketAPI: Authorization header=%v", c.GetHeader("Authorization"))

	userIDValueCheck, existsCheck := c.Get("user_id")
	logrus.Infof("Basket handler user_id=%v exists=%v", userIDValueCheck, existsCheck)

	// --- получаем user_id из JWT ---
	userIDVal, exists := c.Get("user_id")
	logrus.Infof("Basket: exists=%v userIDVal=%v", exists, userIDVal)
	if !exists {
		logrus.Info("basket-4ifd; user_id отсутствует — JWT middleware не сработал")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authenticated"})
		return
	}
	userID := userIDVal.(int)
	logrus.Infof("GetRequestShipBasketAPI: user_id from context = %v", userID)

	// --- получаем или создаём черновик пользователя ---
	requestShip, err := h.Repository.GetOrCreateUserDraft(userID)
	if err != nil {
		logrus.Error("Ошибка GetOrCreateUserDraft: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	logrus.Infof("GetRequestShipBasketAPI: requestShip.ID=%v ships_len=%d", requestShip.RequestShipID, len(requestShip.Ships))
	for _, s := range requestShip.Ships {
		logrus.Infof("  ship: id=%v ships_count=%v", s.ShipID, s.ShipsCount)
	}

	// --- считаем общее количество кораблей ---
	totalShipsCount := 0
	for _, ship := range requestShip.Ships {
		totalShipsCount += ship.ShipsCount
	}

	// --- возвращаем JSON ---
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"request_ship_id": requestShip.RequestShipID,
			"ships_count":     totalShipsCount,
			"user_id":         userID,
		},
	})
}

// GetRequestShipsAPI - GET /api/request_ship - список заявок

// @Summary Получить список заявок на расчет времени погрузки
// @Description Retrieve a list of requests with optional filters
// @Tags request_ships
// @Produce json
// @Param start_date query string false "Start date filter"
// @Param end_date query string false "End date filter"
// @Param status query string false "Status filter"
// @Success 200 {object} []ds.RequestShip
// @Failure 500 {object} object "error: string"
// @Router /api/request_ship [get]
func (h *RequestShipHandler) GetRequestShipsAPI(c *gin.Context) {
	// Получаем фильтры из query-параметров
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	status := c.Query("status")

	// Вызываем репозиторий для получения списка заявок
	requestShips, err := h.Repository.GetRequestShipsFiltered(startDate, endDate, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Возвращаем JSON
	c.JSON(http.StatusOK, requestShips)
}

// GetRequestShipAPI - GET /api/request_ship/:id - одна заявка с услугами

// @Summary Одна заявка на расчет времени погрузки
// @Description Retrieve details of a specific request with its ships
// @Tags request_ships
// @Produce json
// @Param id path int true "Request ID"
// @Success 200 {object} object "request_ship_id: int, status: string, creation_date: string, containers_20ft_count: int, containers_40ft_count: int, comment: string, loading_time: int, ships: []object"
// @Failure 400 {object} object "error: string"
// @Failure 404 {object} object "error: string"
// @Router /api/request_ship/{id} [get]
func (h *RequestShipHandler) GetRequestShipAPI(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request ID",
		})
		return
	}

	requestShip, err := h.Repository.GetRequestShipExcludingDeleted(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Request not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"request_ship_id":       requestShip.RequestShipID,
		"status":                requestShip.Status,
		"creation_date":         requestShip.CreationDate,
		"formation_date":        requestShip.FormationDate,
		"completion_date":       requestShip.CompletionDate,
		"containers_20ft_count": requestShip.Containers20ftCount,
		"containers_40ft_count": requestShip.Containers40ftCount,
		"comment":               requestShip.Comment,
		"loading_time":          requestShip.LoadingTime,
		"ships": func() []gin.H {
			ships := []gin.H{}
			for _, shipInRequest := range requestShip.Ships {
				ships = append(ships, gin.H{
					"ship_id":     shipInRequest.Ship.ShipID,
					"name":        shipInRequest.Ship.Name,
					"photo_url":   shipInRequest.Ship.PhotoURL,
					"capacity":    shipInRequest.Ship.Capacity,
					"cranes":      shipInRequest.Ship.Cranes,
					"ships_count": shipInRequest.ShipsCount,
				})
			}
			return ships
		}(),
	})
}

// UpdateRequestShipAPI - PUT /api/request_ship/:id - изменения полей заявки
// @Summary Изменение полей заявки
// @Description Update fields of an existing request
// @Tags request_ships
// @Accept json
// @Produce json
// @Param id path int true "Request ID"
// @Param request body object{containers_20ft_count=int,containers_40ft_count=int,comment=string} true "Request updates"
// @Success 200 {object} object "status: string, message: string"
// @Failure 400 {object} object "error: string"
// @Failure 500 {object} object "error: string"
// @Router /api/request_ship/{id} [put]
func (h *RequestShipHandler) UpdateRequestShipAPI(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{

			"message": "Invalid request ID",
		})
		return
	}

	var updates struct {
		Containers20ftCount int    `json:"containers_20ft_count"`
		Containers40ftCount int    `json:"containers_40ft_count"`
		Comment             string `json:"comment"`
	}

	if err := c.BindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{

			"error": err.Error(),
		})
		return
	}

	// Обновляем поля без расчета времени (расчет будет при завершении)
	err = h.Repository.UpdateRequestShipFields(id, updates.Containers20ftCount, updates.Containers40ftCount, updates.Comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{

			"error": err.Error(),
		})
		return
	}

	// Проверка на запрос от формы
	if c.PostForm("_method") == "PUT" {
		c.Redirect(http.StatusFound, "/request_ship/"+strconv.Itoa(id))
		return
	}

	c.JSON(http.StatusOK, gin.H{

		"message": "Request updated successfully",
	})
}

// FormRequestShipAPI - PUT /api/request_ship/:id/formation - сформировать создателем + расчёт времени

// @Summary Сформировать заявку на расчет времени погрузки
// @Description Finalize a draft request by the creator
// @Tags request_ships
// @Produce json
// @Param id path int true "Request ID"
// @Success 200 {object} object "status: string, message: string"
// @Failure 400 {object} object "description: string"
// @Failure 404 {object} object "description: string"
// @Failure 500 {object} object "error: string"
// @Router /api/request_ship/{id}/formation [put]
func (h *RequestShipHandler) FormRequestShipAPI(c *gin.Context) {

	bodyBytes, _ := io.ReadAll(c.Request.Body)
	fmt.Println("DEBUG RAW BODY =", string(bodyBytes))
	fmt.Println("DEBUG CONTENT-TYPE =", c.GetHeader("Content-Type"))
	// восстановить body, чтобы BindJSON мог прочитать его ещё раз
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request ID"})
		return
	}

	//  1. Читаем JSON
	var body struct {
		Containers20 int    `json:"containers_20ft"`
		Containers40 int    `json:"containers_40ft"`
		Comment      string `json:"comment"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"description": "Invalid JSON payload",
		})
		return
	}

	// 2. Валидация
	if body.Containers20 <= 0 && body.Containers40 <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"description": "At least one container count must be greater than zero",
		})
		return
	}

	// --- 3. Сохраняем контейнеры ---
	err = h.Repository.UpdateRequestShipContainers(id, body.Containers20, body.Containers40)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update containers: " + err.Error(),
		})
		return
	}

	// --- 4. Меняем статус ---
	err = h.Repository.UpdateRequestShipStatus(id, "сформирован")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// --- 5. Ответ ---
	c.JSON(http.StatusOK, gin.H{
		"message": "Request formed successfully",
	})
}

// CompleteRequestShipAPI - POST /api/request_ship/:id/completion - завершить/отклонить модератором

// @Summary Завершить или отклонить заявку (модератор)
// @Description Allow port_operator to complete or reject a formed request
// @Tags request_ships
// @Produce json
// @Param id path int true "Request ID"
// @Param action formData string true "Action (complete or reject)"
// @Success 200 {object} object "status: string, message: string, loading_time: int (if completed)"
// @Failure 400 {object} object "description: string"
// @Failure 404 {object} object "description: string"
// @Failure 500 {object} object "error: string"
// @Router /api/request_ship/{id}/completion [post]
// CompleteRequestShipAPI - POST /api/request_ship/:id/completion - завершить/отклонить модератором
func (h *RequestShipHandler) CompleteRequestShipAPI(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Errorf("CompleteRequestShipAPI: Invalid request ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request ID",
		})
		return
	}

	action := c.PostForm("action")
	if action == "" {
		logrus.Errorf("CompleteRequestShipAPI: Action must be specified for request_ship_id=%d", id)
		c.JSON(http.StatusBadRequest, gin.H{
			"description": "Action must be specified",
		})
		return
	}

	requestShip, err := h.Repository.GetRequestShipExcludingDeleted(id)
	if err != nil {
		logrus.Errorf("CompleteRequestShipAPI: Request not found for request_ship_id=%d: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{
			"description": "Request not found",
		})
		return
	}

	if requestShip.Status != "сформирован" {
		c.JSON(http.StatusBadRequest, gin.H{
			"description": "Only formed requests can be completed or rejected",
		})
		return
	}

	const portOperatorID = 1

	if action == "complete" {

		go func() {
			if err := h.sendToDjangoLoadingTime(&requestShip); err != nil {
				logrus.Errorf(
					"Failed to send async loading time request for request_ship_id=%d: %v",
					requestShip.RequestShipID,
					err,
				)
			}
		}()

		// завершаем заявку БЕЗ расчёта
		err = h.Repository.CompleteRequestShip(
			id,
			portOperatorID,
			"завершен",
			0,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{
			"message": "Request completed, async calculation started",
		})
		return
	}

	if action == "reject" {
		err = h.Repository.CompleteRequestShip(id, portOperatorID, "отклонен", 0)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Request rejected successfully",
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"description": "Action must be 'complete' or 'reject'",
	})
}

// LoadingTimeCallback - POST /api/request_ship/{id}/loading-time-result

// @Summary Получить результат асинхронного расчёта времени погрузки
// @Description Callback-метод, вызываемый Django-сервисом после асинхронного расчёта
// @Tags request_ships
// @Accept json
// @Produce json
// @Param id path int true "RequestShip ID"
// @Param Authorization header string true "Bearer async token"
// @Success 200 {object} object "message: string"
// @Failure 400 {object} object "error: string"
// @Failure 401 {object} object "error: string"
// @Failure 500 {object} object "error: string"
// @Router /api/request_ship/{id}/loading-time-result [post]

// LoadingTimeCallback — callback от Django после асинхронного расчёта
func (h *RequestShipHandler) LoadingTimeCallback(c *gin.Context) {

	logrus.Info("LoadingTimeCallback has worked!! ", c.Param("id"))

	const ASYNC_TOKEN = "12345678"
	if c.GetHeader("Authorization") != "Bearer "+ASYNC_TOKEN {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid async token",
		})
		return
	}

	// Парсинг данных от Django
	var callbackData ds.DjangoLoadingTimeCallback
	if err := c.ShouldBindJSON(&callbackData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid data format",
		})
		return
	}
	logrus.Infof(
		"Async result получен: request_ship_id=%d success=%v loading_time=%f",
		callbackData.RequestShipID,
		callbackData.Success,
		callbackData.LoadingTime,
	)

	logrus.Info("Получили ответ от Django!!")
	logrus.Info("LoadingTime = ", callbackData.LoadingTime)

	// Получаем заявку для обновления
	requestShip, err := h.Repository.GetRequestShipExcludingDeleted(callbackData.RequestShipID)
	if err != nil {
		logrus.Errorf("Ошибка при получении заявки для обновления: %v", err)
		c.JSON(500, gin.H{
			"error":   err,
			"message": "не удалось получить заявку для обновления",
		})
		return
	}

	// Обновляем результат расчёта
	if callbackData.Success {
		requestShip.LoadingTime = callbackData.LoadingTime

		err = h.Repository.UpdateLoadingTime(
			callbackData.RequestShipID,
			callbackData.LoadingTime,
		)
		if err != nil {
			logrus.Errorf("Ошибка при сохранении loading_time: %v", err)
			c.JSON(500, gin.H{
				"error":   err,
				"message": "Ошибка при сохранении результата расчёта!",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Асинхронный расчёт завершён",
		"request_ship": requestShip,
	})
}

// DeleteShipFromRequestShipAPI - DELETE /api/request_ship/:id/ships/:ship_id - удаление корабля из заявки

// @Summary Удаление корабля из заявки
// @Description Remove a ship from a specific request
// @Tags request_ships
// @Produce json
// @Param id path int true "Request ID"
// @Param ship_id path int true "Ship ID"
// @Success 200 {object} object "description: string"
// @Failure 400 {object} object "status: string, description: string"
// @Failure 500 {object} object "status: string, description: string"
// @Router /api/request_ship/{id}/ships/{ship_id} [delete]
func (h *RequestShipHandler) DeleteShipFromRequestShipAPI(c *gin.Context) {
	requestShipIDStr := c.Param("id")
	shipIDStr := c.Param("ship_id")

	requestShipID, err := strconv.Atoi(requestShipIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "description": "Invalid request ship ID"})
		return
	}

	shipID, err := strconv.Atoi(shipIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "description": "Invalid ship ID"})
		return
	}

	// Удаляем корабль из заявки
	if err := h.Repository.RemoveShipFromRequestShip(requestShipID, shipID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "description": err.Error()})
		return
	}

	// Проверка на запрос от формы
	if c.PostForm("_method") == "DELETE" {
		// Получаем обновленную заявку
		updatedRequestShip, err := h.Repository.GetRequestShipExcludingDeleted(requestShipID)
		if err != nil {
			c.Redirect(http.StatusFound, "/ships")
			return
		}

		if len(updatedRequestShip.Ships) == 0 {
			c.Redirect(http.StatusFound, "/ships")
		} else {
			c.Redirect(http.StatusFound, "/request_ship/"+requestShipIDStr)
		}
		return
	}

	// Для чистых API-запросов возвращаем JSON
	c.JSON(http.StatusOK, gin.H{"description": "Ship removed from request ship"})
}

// UpdateShipInRequestAPI - PUT /api/request_ship/:id/ships/:ship_id - обновление количества кораблей в заявке

// @Summary Обновление количества кораблей в заявке
// @Description Update the number of ships in a specific request
// @Tags request_ships
// @Accept json
// @Produce json
// @Param id path int true "Request ID"
// @Param ship_id path int true "Ship ID"
// @Param request body object{ships_count=int} true "Updated ship count"
// @Success 200 {object} object "status: string, message: string"
// @Failure 400 {object} object "description: string"
// @Failure 500 {object} object "error: string"
// @Router /api/request_ship/{id}/ships/{ship_id} [put]
func (h *RequestShipHandler) UpdateShipInRequestAPI(c *gin.Context) {
	requestShipIDStr := c.Param("id")
	shipIDStr := c.Param("ship_id")

	requestShipID, err := strconv.Atoi(requestShipIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{

			"description": "Invalid request ship ID",
		})
		return
	}

	shipID, err := strconv.Atoi(shipIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{

			"description": "Invalid ship ID",
		})
		return
	}

	var input struct {
		ShipsCount int `json:"ships_count"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{

			"description": "Invalid input data",
		})
		return
	}

	if input.ShipsCount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{

			"description": "Ships count must be greater than zero",
		})
		return
	}

	// Обновляем количество кораблей в заявке
	err = h.Repository.UpdateShipCountInRequest(requestShipID, shipID, input.ShipsCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{

			"error": err.Error(),
		})
		return
	}

	// Проверка на запрос от формы
	if c.PostForm("_method") == "PUT" {
		c.Redirect(http.StatusFound, "/request_ship/"+requestShipIDStr)
		return
	}

	// Для чистых API-запросов — возвращаем JSON
	c.JSON(http.StatusOK, gin.H{

		"message": "Ship count in request updated successfully",
	})
}

// DeleteRequestShipAPI - DELETE /api/request_ship/:id - удаление всей заявки

// @Summary Удаление всей заявки
// @Description Remove an entire request from the system
// @Tags request_ships
// @Produce json
// @Param id path int true "Request ID"
// @Success 200 {object} object "status: string, message: string"
// @Failure 400 {object} object "message: string"
// @Failure 500 {object} object "error: string"
// @Router /api/request_ship/{id} [delete]
func (h *RequestShipHandler) DeleteRequestShipAPI(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Errorf("DeleteRequestShipAPI: Invalid request ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{

			"message": "Invalid request ID",
		})
		return
	}

	logrus.Infof("DeleteRequestShipAPI: Attempting to delete request_ship_id=%d", id)

	// Удаляем зависимые записи
	err = h.Repository.DB().Delete(&ds.ShipInRequest{}, "request_ship_id = ?", id).Error
	if err != nil {
		logrus.Errorf("DeleteRequestShipAPI: Failed to delete ShipInRequest for request_ship_id=%d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{

			"error": err.Error(),
		})
		return
	}

	// Удаляем заявку
	err = h.Repository.DB().Delete(&ds.RequestShip{}, id).Error
	if err != nil {
		logrus.Errorf("DeleteRequestShipAPI: Failed to delete RequestShip for request_ship_id=%d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{

			"error": err.Error(),
		})
		return
	}

	// Проверка на запрос от формы
	if c.PostForm("_method") == "DELETE" {
		logrus.Infof("DeleteRequestShipAPI: Redirecting to /ships for request_ship_id=%d", id)
		c.Redirect(http.StatusFound, "/ships")
		return
	}

	logrus.Infof("DeleteRequestShipAPI: Returning JSON for request_ship_id=%d", id)
	c.JSON(http.StatusOK, gin.H{

		"message": "Request ship deleted successfully",
	})
}
