package services

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/dchlong/billing-be/internal/config"
	"github.com/dchlong/billing-be/internal/models"
	"github.com/dchlong/billing-be/internal/repository"
)

type BillServiceTestSuite struct {
	suite.Suite
	billRepo    *repository.MockIBillRepo
	billService *billService
}

func (s *BillServiceTestSuite) SetupSuite() {
	mockCtrl := gomock.NewController(s.T())
	s.billRepo = repository.NewMockIBillRepo(mockCtrl)
	repo := repository.NewMockIRepo(mockCtrl)
	repo.EXPECT().Bill().Return(s.billRepo).AnyTimes()
	s.billService = &billService{
		repo: repo,
		cfg: &config.AppConfig{
			NumberOfSecondsInABlock: 60,
		},
	}
}

func TestBillServiceTestSuite(t *testing.T) {
	suite.Run(t, new(BillServiceTestSuite))
}

func (s *BillServiceTestSuite) TestCreateCallHistory() {
	ctx := context.Background()
	const (
		testUserName     = "hlong"
		testCallDuration = int64(64000)

		expectedID = int64(1)
	)

	s.billRepo.EXPECT().Create(ctx, &models.CallHistory{
		UserName: testUserName,
		Duration: testCallDuration,
	}).DoAndReturn(func(ctx context.Context, records ...*models.CallHistory) error {
		s.Require().Len(records, 1)
		records[0].ID = expectedID
		return nil
	})
	resp, err := s.billService.CreateCallHistory(ctx, testUserName, testCallDuration)
	s.Require().NoError(err)
	s.Require().Equal(expectedID, resp.ID)
}

func (s *BillServiceTestSuite) TestCreateCallHistory_DatabaseError() {
	ctx := context.Background()
	const (
		testUserName     = "hlong"
		testCallDuration = int64(64000)

		expectedID = int64(1)
	)

	s.billRepo.EXPECT().Create(ctx, &models.CallHistory{
		UserName: testUserName,
		Duration: testCallDuration,
	}).DoAndReturn(func(ctx context.Context, records ...*models.CallHistory) error {
		s.Require().Len(records, 1)
		records[0].ID = expectedID
		return gorm.ErrInvalidDB
	})

	resp, err := s.billService.CreateCallHistory(ctx, testUserName, testCallDuration)
	s.Require().ErrorIs(err, gorm.ErrInvalidDB)
	s.Require().Nil(resp)
}

func (s *BillServiceTestSuite) TestGetBill() {
	ctx := context.Background()
	var testUserName = "hlong"
	s.billRepo.EXPECT().FindBy(ctx, &models.CallHistoryFilter{
		UserName: &testUserName,
	}).Return([]*models.CallHistory{
		{
			ID:       1,
			UserName: testUserName,
			Duration: 60000,
		},
		{
			ID:       2,
			UserName: testUserName,
			Duration: 120000,
		},
	}, nil)
	resp, err := s.billService.GetBill(ctx, testUserName)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Equal(int64(3), resp.BlockCount)
	s.Require().Equal(int64(2), resp.CallCount)
}

func (s *BillServiceTestSuite) TestGetBill_Round() {
	ctx := context.Background()
	var testUserName = "hlong"
	s.billRepo.EXPECT().FindBy(ctx, &models.CallHistoryFilter{
		UserName: &testUserName,
	}).Return([]*models.CallHistory{
		{
			ID:       1,
			UserName: testUserName,
			Duration: 64000,
		},
		{
			ID:       2,
			UserName: testUserName,
			Duration: 96000,
		},
	}, nil)
	resp, err := s.billService.GetBill(ctx, testUserName)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Equal(int64(3), resp.BlockCount)
	s.Require().Equal(int64(2), resp.CallCount)
}

func (s *BillServiceTestSuite) TestGetBill_FindError() {
	ctx := context.Background()
	var testUserName = "hlong"
	s.billRepo.EXPECT().FindBy(ctx, &models.CallHistoryFilter{
		UserName: &testUserName,
	}).Return(nil, gorm.ErrInvalidDB)
	resp, err := s.billService.GetBill(ctx, testUserName)
	s.Require().ErrorIs(err, gorm.ErrInvalidDB)
	s.Require().Nil(resp)
}
