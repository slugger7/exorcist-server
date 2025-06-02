package models

type PageDTO[T any] struct {
	Data  []T `json:"data"`
	Total int `json:"total"`
	Limit int `json:"limit"`
	Skip  int `json:"skip"`
}

func DataToPage[T any, S any](data []T, o PageDTO[S]) PageDTO[T] {
	return PageDTO[T]{
		Limit: o.Limit,
		Skip:  o.Skip,
		Total: o.Total,
		Data:  data,
	}
}
