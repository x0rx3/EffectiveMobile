package model

type Filter struct {
	Field      string `json:"field"`
	Value      string `json:"value"`
	FilterType string `json:"type"` // eq, neq, like, gt, gte, lt, lte, in
}

type Sorter struct {
	Field     string `json:"field"`
	Direction string `json:"direction"` // asc | desc
}

type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type ListData struct {
	Filters    []Filter   `json:"filters"`
	Sorters    []Sorter   `json:"sorters"`
	Pagination Pagination `json:"pagination"`
}
