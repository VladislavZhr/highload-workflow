package model

type SuccessResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
