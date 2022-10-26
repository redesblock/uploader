package model

type Voucher struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	Voucher string `json:"voucher" gorm:"unique"`
	Node    string `json:"node"`
	Usable  bool   `json:"usable"`
}
