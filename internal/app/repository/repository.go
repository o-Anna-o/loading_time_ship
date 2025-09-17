package repository

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Repository struct {
	Ships    []Ship
	Requests map[int]Request // Изменил на map[int]Request
}

type Ship struct {
	ID         int
	Name       string
	Speed      int
	Capacity   float32
	Length     float32
	Width      float32
	Draft      float32
	Cranes     int
	Containers int
	Features   string
	PhotoURL   string
}

type ShipInRequest struct {
	Ship  Ship
	Count int
}

type Request struct {
	ID                  int
	Ships               []ShipInRequest
	Containers20ftCount int
	Containers40ftCount int
	Comment             string
}

func NewRepository() (*Repository, error) {
	ships := []Ship{
		{
			ID:         1,
			Name:       "Ever Ace",
			Speed:      242,
			Capacity:   23992,
			Length:     400,
			Width:      61.53,
			Draft:      17.0,
			Cranes:     6,
			Containers: 11996,
			Features:   "самый большой в мире, двигатель Wartsila 70950 кВт",
			PhotoURL:   "ever-ace.png",
		},
		{
			ID:         2,
			Name:       "FESCO Diomid",
			Speed:      105,
			Capacity:   3108,
			Length:     195,
			Width:      32.20,
			Draft:      11.0,
			Cranes:     3,
			Containers: 536,
			Features:   "построен в 2010 г., судно класса Ice1 (для Арктики), дизельный двигатель, используется на Северном морском пути",
			PhotoURL:   "fesco-diomid.png",
		},
		{
			ID:         3,
			Name:       "HMM Algeciras",
			Speed:      243,
			Capacity:   23964,
			Length:     399.9,
			Width:      61.0,
			Draft:      16.5,
			Cranes:     7,
			Containers: 11982,
			Features:   "двигатель MAN B&W 11G95ME-C9.5 мощностью 64 000 кВт, двойные двигатели, система рекуперации энергии, класс DNV GL",
			PhotoURL:   "hmm-algeciras.png",
		},
		{
			ID:         4,
			Name:       "MSC Gulsun",
			Speed:      245,
			Capacity:   23756,
			Length:     399.9,
			Width:      61.4,
			Draft:      16.0,
			Cranes:     7,
			Containers: 11878,
			Features:   "первый в мире контейнеровоз, вмещающий более 23 000 TEU, двигатель MAN B&W 11G95ME-C9.5, класс DNV GL",
			PhotoURL:   "msc-gulsun.png",
		},
	}

	if len(ships) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	requests := make(map[int]Request)
	requests[1] = Request{ID: 1, Ships: []ShipInRequest{}}

	return &Repository{Ships: ships, Requests: requests}, nil
}

func (r *Repository) GetShips() ([]Ship, error) {
	// имитируем работу с БД
	ships := []Ship{
		{
			ID:         1,
			Name:       "Ever Ace",
			Speed:      242,
			Capacity:   23992,
			Length:     400,
			Width:      61.53,
			Draft:      17.0,
			Cranes:     6,
			Containers: 11996,
			Features:   "самый большой в мире, двигатель Wartsila 70950 кВт",
			PhotoURL:   "ever-ace.png",
		},
		{
			ID:         2,
			Name:       "FESCO Diomid",
			Speed:      105,
			Capacity:   3108,
			Length:     195,
			Width:      32.20,
			Draft:      11.0,
			Cranes:     3,
			Containers: 536,
			Features:   "построен в 2010 г., судно класса Ice1 (для Арктики), дизельный двигатель, используется на Северном морском пути",
			PhotoURL:   "fesco-diomid.png",
		},
		{
			ID:         3,
			Name:       "HMM Algeciras",
			Speed:      243,
			Capacity:   23964,
			Length:     399.9,
			Width:      61.0,
			Draft:      16.5,
			Cranes:     7,
			Containers: 11982,
			Features:   "двигатель MAN B&W 11G95ME-C9.5 мощностью 64 000 кВт, двойные двигатели, система рекуперации энергии, класс DNV GL",
			PhotoURL:   "hmm-algeciras.png",
		},
		{
			ID:         4,
			Name:       "MSC Gulsun",
			Speed:      245,
			Capacity:   23756,
			Length:     399.9,
			Width:      61.4,
			Draft:      16.0,
			Cranes:     7,
			Containers: 11878,
			Features:   "первый в мире контейнеровоз, вмещающий более 23 000 TEU, двигатель MAN B&W 11G95ME-C9.5, класс DNV GL",
			PhotoURL:   "msc-gulsun.png",
		},
	}

	if len(ships) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return r.Ships, nil
}

func (r *Repository) GetShip(id int) (Ship, error) {

	ships, err := r.GetShips()
	if err != nil {
		return Ship{}, err
	}

	for _, ship := range ships {
		if ship.ID == id {
			return ship, nil
		}
	}
	return Ship{}, fmt.Errorf("Контейнеровоз не найден")
}

func (r *Repository) GetShipsByName(name string) ([]Ship, error) {
	ships, err := r.GetShips()
	if err != nil {
		return []Ship{}, err
	}

	var result []Ship
	for _, ship := range ships {
		if strings.Contains(strings.ToLower(ship.Name), strings.ToLower(name)) {
			result = append(result, ship)
		}
	}

	return result, nil
}

func (r *Repository) GetShipsByCapacity(capacity string) ([]Ship, error) {
	ships, err := r.GetShips()
	if err != nil {
		return []Ship{}, err
	}
	searchCapacity, err := strconv.ParseFloat(capacity, 32)
	if err != nil {
		return []Ship{}, nil
	}

	var result []Ship
	for _, ship := range ships {
		diff := math.Abs(float64(ship.Capacity) - searchCapacity)
		if diff <= searchCapacity*0.30 {
			result = append(result, ship)
		}
	}
	return result, nil
}

func (r *Repository) GetRequest(id int) (Request, error) {
	if request, ok := r.Requests[id]; ok {
		return request, nil
	}
	return Request{}, fmt.Errorf("заявка с id=%d не найдена", id)
}

func (r *Repository) RemoveShipFromRequest(requestID int, shipID int) error {
	request, ok := r.Requests[requestID]
	if !ok {
		return fmt.Errorf("заявка с id=%d не найдена", requestID)
	}
	for i, shipInRequest := range request.Ships {
		if shipInRequest.Ship.ID == shipID {
			request.Ships = append(request.Ships[:i], request.Ships[i+1:]...)
			r.Requests[requestID] = request
			return nil
		}
	}

	return fmt.Errorf("корабль с id=%d не найден в заявке", shipID)
}
