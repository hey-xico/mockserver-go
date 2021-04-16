package model

type VerifierRequest struct {
	HttpRequest HttpRequest       `json:"httpRequest,omitempty"`
	Times       VerificationTimes `json:"times"`
}
