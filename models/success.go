package models

type Success struct {
	Message string            `json:"message,omitempty"`
	Success bool           `json:"success,omitempty"`
}