package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// AddShipToRequestShip добавляет корабль в заявку
func (h *Handler) AddShipToRequestShip(c *gin.Context) {
	shipIDStr := c.Param("ship_id")
	shipID, err := strconv.Atoi(shipIDStr)
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	request_ship, err := h.Repository.GetOrCreateUserDraft(1)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	err = h.Repository.AddShipToRequestShip(request_ship.RequestShipID, shipID)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	logrus.Infof("Корабль %d добавлен в заявку %d через ORM", shipID, request_ship.RequestShipID)
	c.Redirect(http.StatusFound, fmt.Sprintf("/request_ship/%d", request_ship.RequestShipID))
}

// DeleteRequestShip - логическое удаление заявки
func (h *Handler) DeleteRequestShip(c *gin.Context) {
	requestShipID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	if err := h.Repository.DeleteRequestShipSQL(requestShipID); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusFound, "/ships")
}

// RemoveShipFromRequestShip удаляет корабль из заявки
func (h *Handler) RemoveShipFromRequestShip(c *gin.Context) {
	request_shipIDStr := c.Param("id")
	shipIDStr := c.Param("ship_id")

	request_shipID, err := strconv.Atoi(request_shipIDStr)
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	shipID, err := strconv.Atoi(shipIDStr)
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.RemoveShipFromRequestShip(request_shipID, shipID)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusFound, "/request_ship/"+request_shipIDStr)
}
