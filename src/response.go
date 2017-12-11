package main

type ListResponse struct {
	Count 			int 		`json:"total"`
	List 			[]UrlItem	`json:"list"`
}

type CreateResponse struct {
	Status 		string   `json:"status"`
	Created 	UrlItem  `json:"url"`
}

type SuccessResponse struct {
	Status 		string	`json:"status"`
	Message 	string	`json:"message"`
}

type ErrorResponse struct {
	Status 		string   `json:"status"`
	Response 	string   `json:"message"`
}

func buildErrorResponse(message string) ErrorResponse {
	var responseMessage ErrorResponse
	responseMessage.Status = "error"
	responseMessage.Response = message
	return responseMessage
}

func buildSuccessResponse(message string) SuccessResponse {
	var responseMessage SuccessResponse
	responseMessage.Status = "ok"
	responseMessage.Message = message
	return responseMessage
}

func buildCreateResponse(item UrlItem) CreateResponse {
	var responseMessage CreateResponse
	responseMessage.Status = "ok"
	responseMessage.Created = item
	return responseMessage
}
