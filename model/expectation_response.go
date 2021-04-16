package model

type ExpectationResponse struct {
	Id          string      `json:"id"`
	HttpRequest HttpRequest `json:"httpRequest"`
}
