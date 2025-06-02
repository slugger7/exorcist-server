package dto

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

type PageRequestDTO struct {
	Skip  int  `form:"skip" json:"skip"`
	Limit int  `form:"limit" json:"limit"`
	Asc   bool `form:"asc" json:"asc"`
}
