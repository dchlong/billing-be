package services

import (
	"context"
	"time"

	"github.com/dchlong/billing-be/internal/repository"

	"github.com/dchlong/billing-be/internal/config"
	"github.com/dchlong/billing-be/internal/models"
	"github.com/dchlong/billing-be/internal/presenter"
)

//go:generate mockgen -source=bill.go -destination=bill.mock.go -package=services

type BillService interface {
	CreateCallHistory(ctx context.Context, userName string, duration int64) (*presenter.CreateCallHistoryResponse, error)
	GetBill(ctx context.Context, userName string) (*presenter.GetBillResponse, error)
}

type billService struct {
	repo repository.IRepo
	cfg  *config.AppConfig
}

func ProvideBillService(repo repository.IRepo, cfg *config.AppConfig) BillService {
	return &billService{
		repo: repo,
		cfg:  cfg,
	}
}

func (b *billService) CreateCallHistory(
	ctx context.Context, userName string, duration int64,
) (*presenter.CreateCallHistoryResponse, error) {
	callHistory := &models.CallHistory{
		UserName: userName,
		Duration: duration,
	}

	err := b.repo.Bill().Create(ctx, callHistory)
	if err != nil {
		return nil, err
	}

	return presenter.NewCreateCallHistoryResponse(callHistory), nil
}

func (b *billService) GetBill(ctx context.Context, userName string) (*presenter.GetBillResponse, error) {
	callHistories, err := b.repo.Bill().FindBy(ctx, models.NewCallHistoryFilter().WithUserName(userName))
	if err != nil {
		return nil, err
	}

	var callCount int64
	var totalDuration int64
	for _, callHistory := range callHistories {
		callCount++
		totalDuration += callHistory.Duration
	}

	numberOfMillisecondInABlock := b.cfg.NumberOfSecondsInABlock * int64(time.Second/time.Millisecond)
	blockCount := totalDuration / numberOfMillisecondInABlock
	if totalDuration%numberOfMillisecondInABlock > 0 {
		blockCount++
	}

	return &presenter.GetBillResponse{
		CallCount:  callCount,
		BlockCount: blockCount,
	}, nil
}
