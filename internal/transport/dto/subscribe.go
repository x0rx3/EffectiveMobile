package dto

import "effective_mobile/pkg/model"

type ReadSubscribeDto struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	UserID    string     `json:"user_id"`
	Price     int        `json:"price"`
	StartDate model.Date `json:"start_date"`
	EndDate   model.Date `json:"end_date"`
}

type CreateUpdateSubscribeDto struct {
	Name      string     `json:"name"`
	UserID    string     `json:"user_id"`
	Price     int        `json:"price"`
	StartDate model.Date `json:"start_date"`
	EndDate   model.Date `json:"end_date"`
}
