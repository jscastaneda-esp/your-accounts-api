package model

import "time"

type CreateRequest struct {
	UID   string `json:"uid,omitempty" validate:"required,len=32"`
	Email string `json:"email,omitempty" validate:"required,email"`
}

type CreateResponse struct {
	ID        uint      `json:"id,omitempty"`
	UID       string    `json:"uid,omitempty"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type AuthRequest struct {
	CreateRequest
}

type AuthResponse struct {
	Token string `json:"token,omitempty"`
}

type RefreshTokenRequest struct {
	CreateRequest
	Token string `json:"token,omitempty" validate:"required"`
}

type RefreshTokenResponse struct {
	AuthResponse
}
