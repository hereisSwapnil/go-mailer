package types

type Recipient struct {
	Name   string            `json:"name"`
	Email  string            `json:"email"`
	Extra  map[string]string `json:"extra,omitempty"`
}