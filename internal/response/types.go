package response

type JsonResponse struct {
	Status     string      `json:"status"`
	StatusCode int         `json:"statusCode"`
	Data       interface{} `json:"data"`
}
