package ds

import "time"

type RequestShip struct {
	RequestShipID       int             `gorm:"primaryKey;column:request_ship_id"`
	Status              string          `gorm:"column:status"`
	CreationDate        time.Time       `gorm:"column:creation_date"`
	UserID              int             `gorm:"column:user_id"`
	CompletionDate      *time.Time      `gorm:"column:completion_date"`
	Containers20ftCount int             `gorm:"column:containers_20ft_count"`
	Containers40ftCount int             `gorm:"column:containers_40ft_count"`
	Comment             string          `gorm:"column:comment"`
	LoadingTime         float64         `gorm:"column:loading_time"`
	Ships               []ShipInRequest `gorm:"foreignKey:RequestShipID"`
}

func (RequestShip) TableName() string {
	return "request_ship"
}
