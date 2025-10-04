package handler

import (
	"loading_time/internal/app/ds"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetShips(ctx *gin.Context) {
	var ships []ds.Ship
	var err error

	searchQuery := ctx.Query("search")
	if searchQuery == "" {
		ships, err = h.Repository.GetShips()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		ships, err = h.Repository.GetShipsByName(searchQuery)
		if err != nil {
			logrus.Error(err)
		}
	}

	// для подсчета кораблей в заявке - используем ORM
	requestShipCount := 0
	var requestShipID int // для ID заявки

	requestShip, err := h.Repository.GetOrCreateUserDraft(1)
	if err == nil {
		// Проверяем что заявка не удалена (GetOrCreateUserDraft уже исключает удаленные по статусу "черновик")
		logrus.Infof("Найдена заявка ID=%d, количество кораблей в заявке: %d", requestShip.RequestShipID, len(requestShip.Ships))
		for i, shipInRequest := range requestShip.Ships {
			logrus.Infof("Корабль %d: ID=%d, количество: %d", i, shipInRequest.ShipID, shipInRequest.ShipsCount)
			requestShipCount += shipInRequest.ShipsCount
		}
		// СОХРАНЯЕМ ID ЗАЯВКИ ДЛЯ ШАБЛОНА
		requestShipID = requestShip.RequestShipID
	} else {
		logrus.Errorf("Ошибка получения заявки: %v", err)
		// Если ошибка, всё равно создаем/получаем заявку для получения ID
		requestShip, err = h.Repository.GetOrCreateUserDraft(1)
		if err == nil {
			requestShipID = requestShip.RequestShipID
		}
	}

	logrus.Infof("Итоговый счетчик для отображения: %d", requestShipCount)

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"ships":              ships,
		"search":             searchQuery,
		"request_ship_count": requestShipCount,
		"request_ship_id":    requestShipID,
	})
}

func (h *Handler) GetShip(ctx *gin.Context) {
	idStr := ctx.Param("id") // получаем id корабля из урла (то есть из /ship/:id)
	// через двоеточие мы указываем параметры, которые потом сможем считать через функцию выше
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
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
