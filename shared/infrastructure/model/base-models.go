package model

type IDResponse struct {
	ID uint `json:"id"`
}

func NewIDResponse(id uint) IDResponse {
	return IDResponse{
		ID: id,
	}
}
