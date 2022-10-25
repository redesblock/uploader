package model

type Voucher struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	Host    string `json:"host"`
	Voucher string `json:"voucher" gorm:"unique"`
	Usable  bool   `json:"usable"`
}
