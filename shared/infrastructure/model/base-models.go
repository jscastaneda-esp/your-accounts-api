package model

type IDResponse struct {
	ID uint `json:"id"`
}

func NewIDResponse(id uint) IDResponse {
	return IDResponse{
		ID: id,
	}
}

type NameResponse struct {
	Name string `json:"name"`
}

func NewNameResponse(name string) NameResponse {
	return NameResponse{
		Name: name,
	}
}

type AmountResponse struct {
	Amount float64 `json:"amount"`
}

func NewAmountResponse(amount float64) AmountResponse {
	return AmountResponse{
		Amount: amount,
	}
}
