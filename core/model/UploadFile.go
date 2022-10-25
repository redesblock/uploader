package model

import (
	"time"
)

type UploadFile struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Path        string    `json:"path" gorm:"unique"`
	IndexName   string    `json:"index_name"`
	RelPath     string    `json:"rel_path"`
	Hash        string    `json:"hash"`
	ModifyAt    time.Time `json:"modify_at"`
	UpdatedAt   time.Time `json:"-"`
	CreatedAt   time.Time `json:"-"`
	WatchFileID uint      `json:"-"`
}
