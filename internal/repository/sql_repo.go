package repository

import (
	"gorm.io/gorm"
)

type SQLRepo struct {
	db *gorm.DB
}

func NewSQLRepo(db *gorm.DB) *SQLRepo {
	return &SQLRepo{db: db}
}

func (r *SQLRepo) Bill() IBillRepo {
	return NewBillSQLRepo(r.db)
}
