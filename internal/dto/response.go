package dto

type Response struct {
	Status int `json:"status"`
	Data   any `json:"data"`
}
