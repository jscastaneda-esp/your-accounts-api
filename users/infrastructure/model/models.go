package model

import (
	"time"
	"your-accounts-api/shared/infrastructure/model"
)

type CreateRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type CreateResponse struct {
	model.IDResponse
}

func NewCreateResponse(id uint) *CreateResponse {
	return &CreateResponse{
		model.NewIDResponse(id),
	}
}

type LoginRequest struct {
	CreateRequest
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
}

func NewLoginResponse(token string, expiresAt time.Time) *LoginResponse {
	return &LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt.UnixMilli(),
	}
}
