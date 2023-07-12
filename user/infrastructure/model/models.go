package model

import "your-accounts-api/shared/infrastructure/model"

type CreateRequest struct {
	UID   string `json:"uid,omitempty" validate:"required,min=28,max=32"`
	Email string `json:"email,omitempty" validate:"required,email"`
}

type CreateResponse struct {
	model.IDResponse
}

type LoginRequest struct {
	CreateRequest
}

type LoginResponse struct {
	Token string `json:"token,omitempty"`
}

func NewCreateResponse(id uint) *CreateResponse {
	return &CreateResponse{
		model.IDResponse{
			ID: id,
		},
	}
}

func NewLoginResponse(token string) *LoginResponse {
	return &LoginResponse{
		Token: token,
	}
}
