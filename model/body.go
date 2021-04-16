package model

type Body struct {
	ContentType string      `json:"contentType"`
	Type        string      `json:"type"`
	MatchType   string      `json:"matchType"`
	Json        interface{} `json:"json"`
}
