package model

type ReadResponse struct {
	ID               uint    `json:"id,omitempty"`
	Name             string  `json:"name,omitempty"`
	Year             uint16  `json:"year,omitempty"`
	Month            uint8   `json:"month,omitempty"`
	FixedIncome      float64 `json:"fixedIncome,omitempty"`
	AdditionalIncome float64 `json:"additionalIncome,omitempty"`
	TotalBalance     float64 `json:"totalBalance,omitempty"`
	Total            float64 `json:"total,omitempty"`
	EstimatedBalance float64 `json:"estimatedBalance,omitempty"`
	ProjectId        uint    `json:"projectId,omitempty"`
}
