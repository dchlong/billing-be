package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/dchlong/billing-be/internal/models"
)

type BillSQLRepo struct {
	db *gorm.DB
}

func (r *BillSQLRepo) Create(ctx context.Context, records ...*models.CallHistory) error {
	err := r.db.WithContext(ctx).Create(records).Error
	return err
}

func (r *BillSQLRepo) buildFilter(db *gorm.DB, filter *models.CallHistoryFilter) *gorm.DB {
	if filter == nil {
		return db
	}

	if filter.UserName != nil {
		db = db.Where("user_name = ?", *filter.UserName)
	}

	return db
}

func (r *BillSQLRepo) FindByID(ctx context.Context, id int64) (*models.CallHistory, error) {
	var record *models.CallHistory
	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&record).Error
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (r *BillSQLRepo) FindBy(ctx context.Context, filter *models.CallHistoryFilter) ([]*models.CallHistory, error) {
	var records []*models.CallHistory
	err := r.buildFilter(r.db.WithContext(ctx), filter).Find(&records).Error
	if err != nil {
		return nil, err
	}

	return records, nil
}

func NewBillSQLRepo(db *gorm.DB) IBillRepo {
	return &BillSQLRepo{
		db: db,
	}
}
