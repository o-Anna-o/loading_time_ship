package repository

import (
	"loading_time/internal/app/ds"
	"time"

	"gorm.io/gorm"
)

func (r *Repository) GetRequestShip(id int) (ds.RequestShip, error) {
	request_ship := ds.RequestShip{}
	// обязательно проверяем ошибки, и если они появились - передаем выше, то есть хендлеру
	err := r.db.Preload("Ships.Ship").Where("id = ?", id).First(&request_ship).Error
	if err != nil {
		return ds.RequestShip{}, err
	}
	return request_ship, nil
}

// GetOrCreateUserDraft - перейти или создать черновик
func (r *Repository) GetOrCreateUserDraft(dummyUserID int) (ds.RequestShip, error) {
	var requestShip ds.RequestShip

	// Ищем существующий черновик для данного пользователя
	err := r.db.Preload("Ships.Ship").Where("status = ? AND user_id = ?", "черновик", dummyUserID).First(&requestShip).Error
	if err == nil {
		return requestShip, nil // черновик найден
	}
	if err != gorm.ErrRecordNotFound {
		return ds.RequestShip{}, err // реальная ошибка
	}

	// Создаем новый черновик
	requestShip = ds.RequestShip{
		Status:       "черновик",
		UserID:       dummyUserID,
		CreationDate: time.Now(),
	}

	err = r.db.Create(&requestShip).Error
	if err != nil {
		return ds.RequestShip{}, err
	}

	return requestShip, nil
}

// AddShipToRequestShip - добавить корабль в заявку через ORM
func (r *Repository) AddShipToRequestShip(requestShipID, shipID int) error {
	// Сначала проверяем, есть ли уже такой корабль в заявке
	var existingShip ds.ShipInRequest
	err := r.db.Where("request_ship_id = ? AND ship_id = ?", requestShipID, shipID).First(&existingShip).Error

	if err == nil {
		// Корабль уже есть в заявке - увеличиваем количество
		existingShip.ShipsCount++
		return r.db.Save(&existingShip).Error
	}

	shipInRequest := ds.ShipInRequest{
		RequestShipID: requestShipID,
		ShipID:        shipID,
		ShipsCount:    1,
	}
	return r.db.Create(&shipInRequest).Error
}

// RemoveShipFromRequestShip - удалить корабль из заявки
func (r *Repository) RemoveShipFromRequestShip(requestShipID, shipID int) error {
	return r.db.Where("request_ship_id = ? AND ship_id = ?", requestShipID, shipID).Delete(&ds.ShipInRequest{}).Error
}

// логическое удаление заявки через SQL
func (r *Repository) DeleteRequestShipSQL(requestShipID int) error {
	return r.db.Exec("UPDATE request_ship SET status = 'удалён' WHERE request_ship_id = ?", requestShipID).Error
}

// GetRequestShipExcludingDeleted - получить заявку исключая удаленные (через ORM)
func (r *Repository) GetRequestShipExcludingDeleted(id int) (ds.RequestShip, error) {
	var requestShip ds.RequestShip
	err := r.db.Preload("Ships.Ship").Where("request_ship_id = ? AND status != ?", id, "удалён").First(&requestShip).Error
	if err != nil {
		return ds.RequestShip{}, err
	}
	return requestShip, nil
}
