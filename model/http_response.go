package model

type HttpResponse struct {
	StatusCode int         `json:"statusCode,omitempty"`
	Body       interface{} `json:"body,omitempty"`
}
