package models

type Page[T any] struct {
	Data  []T `json:"data"`
	Total int `json:"total"`
	Limit int `json:"limit"`
	Skip  int `json:"skip"`
}
