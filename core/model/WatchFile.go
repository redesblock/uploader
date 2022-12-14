package model

import "time"

type WatchFile struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Path      string    `json:"path" gorm:"unique"`
	IndexExt  string    `json:"index_ext"`
	UpdatedAt time.Time `json:"-"`
	CreatedAt time.Time `json:"-"`
}
