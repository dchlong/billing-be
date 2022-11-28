package repository

import (
	"context"

	"github.com/dchlong/billing-be/internal/models"
)

//go:generate mockgen -source=irepo.go -destination=irepo.mock.go -package=repository

type IRepo interface {
	Bill() IBillRepo
}

type IBillRepo interface {
	Create(ctx context.Context, records ...*models.CallHistory) error
	FindByID(ctx context.Context, id int64) (*models.CallHistory, error)
	FindBy(ctx context.Context, filter *models.CallHistoryFilter) ([]*models.CallHistory, error)
}
