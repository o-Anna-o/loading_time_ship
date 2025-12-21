package ds

// DjangoLoadingTimeRequest — отправка данных в Django
type DjangoLoadingTimeRequest struct {
	RequestShipID  int            `json:"request_ship_id"`
	Containers20ft int            `json:"containers_20ft"`
	Containers40ft int            `json:"containers_40ft"`
	Ships          []ShipCalcData `json:"ships"`
}

// ShipCalcData — данные по кораблям
type ShipCalcData struct {
	Cranes     int `json:"cranes"`
	ShipsCount int `json:"ships_count"`
}

// DjangoLoadingTimeCallback — callback от Django
type DjangoLoadingTimeCallback struct {
	RequestShipID int     `json:"request_ship_id"`
	Success       bool    `json:"success"`
	LoadingTime   float64 `json:"loading_time,omitempty"`
	ErrorMessage  string  `json:"error_message,omitempty"`
}
