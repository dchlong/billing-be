package presenter

import (
	"github.com/dchlong/billing-be/internal/models"
)

type GetBillResponse struct {
	CallCount  int64 `json:"call_count"`
	BlockCount int64 `json:"block_count"`
}

type CreateCallHistoryInput struct {
	CallDuration int64 `json:"call_duration" binding:"required,gt=0"`
}

type CreateCallHistoryResponse struct {
	ID int64 `json:"id"`
}

func NewCreateCallHistoryResponse(callHistory *models.CallHistory) *CreateCallHistoryResponse {
	return &CreateCallHistoryResponse{
		ID: callHistory.ID,
	}
}
