package model

import "time"

type CreateRequest struct {
	UUID  string `json:"uuid" validate:"required,len=32"`
	Email string `json:"email" validate:"required,email"`
}

type CreateResponse struct {
	ID        uint      `json:"id"`
	UUID      string    `json:"uuid"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginRequest struct {
	CreateRequest
}
