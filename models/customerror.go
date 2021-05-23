package models
type CustomError struct {
	Message string            `json:"message,omitempty"`
	IsPasswordError bool      `json:"isPasswordError,omitempty"`
	Errors []string           `json:"errors,omitempty" `
	Success bool           `json:"success,omitempty"`
}