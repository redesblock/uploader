package model

import (
	"time"
)

type UploadFile struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Path      string    `json:"path" gorm:"unique"`
	RelPath   string    `json:"rel_path" gorm:"unique"`
	IndexName string    `json:"index_name"`
	Hash      string    `json:"hash"`
	Usable    bool      `json:"usable"`
	ModifyAt  time.Time `json:"modify_at"`
	UpdatedAt time.Time `json:"-"`
	CreatedAt time.Time `json:"-"`
}
