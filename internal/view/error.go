package view

type ErrorDetails struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	ErrorMessage ErrorDetails `json:"error"`
}
