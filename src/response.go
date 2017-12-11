package main

// Map list response
type ListResponse struct {
	Count   int         `json:"total"`
	List    []UrlItem   `json:"list"`
}

// Default Create Success Response
type CreateResponse struct {
	Status      string  `json:"status"`
	Created     UrlItem `json:"url"`
}

// Default Success Response
type SuccessResponse struct {
	Status  string  `json:"status"`
	Message string  `json:"message"`
}

// Default Error Response Mapping
type ErrorResponse struct {
	Status      string  `json:"status"`
	Response    string  `json:"message"`
}

// Factory for the error response
func buildErrorResponse(message string) ErrorResponse {
	var responseMessage ErrorResponse
	responseMessage.Status = "error"
	responseMessage.Response = message
	return responseMessage
}

// Factory for the success response
func buildSuccessResponse(message string) SuccessResponse {
	var responseMessage SuccessResponse
	responseMessage.Status = "ok"
	responseMessage.Message = message
	return responseMessage
}

// Factory for the create success response
func buildCreateResponse(item UrlItem) CreateResponse {
	var responseMessage CreateResponse
	responseMessage.Status = "ok"
	responseMessage.Created = item
	return responseMessage
}
