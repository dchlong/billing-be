package models

import (
	"time"

	lgorm "github.com/dchlong/billing-be/pkg/gorm"
)

type CallHistory struct {
	ID        int64           `json:"id"`
	UserName  string          `json:"user_name"`
	Duration  int64           `json:"duration"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt lgorm.DeletedAt `json:"deleted_at"`
}

type CallHistoryFilter struct {
	UserName *string `json:"user_name"`
}

func NewCallHistoryFilter() *CallHistoryFilter {
	return &CallHistoryFilter{}
}

func (f *CallHistoryFilter) WithUserName(userName string) *CallHistoryFilter {
	f.UserName = &userName
	return f
}
