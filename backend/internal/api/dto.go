package api

// DTOs or Data Transfer Objects are used to define the structure
// of the data that is sent and received by the API.

type Request struct {
	Code     string `json:"code"`
	Language string `json:"language"`
}

type Response struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}
