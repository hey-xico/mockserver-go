package model

type ExpectationRequest struct {
	HttpRequest  HttpRequest  `json:"httpRequest"`
	HttpResponse HttpResponse `json:"httpResponse,omitempty"`
}
