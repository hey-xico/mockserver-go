package model

type HttpRequest struct {
	Method          string                 `json:"method,omitempty"`
	Path            string                 `json:"path"`
	PathParameters  map[string][]string    `json:"pathParameters,omitempty"`
	QueryParameters map[string]interface{} `json:"queryStringParameters,omitempty"`
	Cookies         map[string]string      `json:"cookies,omitempty"`
	Body            *Body                  `json:"body,omitempty"`
}
