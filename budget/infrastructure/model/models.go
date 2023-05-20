package model

type ReadResponse struct {
	ID               uint    `json:"id,omitempty"`
	Name             string  `json:"name,omitempty"`
	Year             uint16  `json:"year,omitempty"`
	Month            uint8   `json:"month,omitempty"`
	FixedIncome      float64 `json:"fixed_income,omitempty"`
	AdditionalIncome float64 `json:"additional_income,omitempty"`
	TotalBalance     float64 `json:"total_balance,omitempty"`
	Total            float64 `json:"total,omitempty"`
	EstimatedBalance float64 `json:"estimated_balance,omitempty"`
	ProjectId        uint    `json:"project_id,omitempty"`
}
