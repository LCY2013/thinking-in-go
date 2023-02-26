package model

type Payment struct {
	// 主键
	ID            int64  `gorm:"column:id;primary_key;not_null;auto_increment" json:"id"`
	PaymentName   string `gorm:"column:payment_name" json:"payment_name"`     // 支付名称
	PaymentSID    string `gorm:"column:payment_sid" json:"payment_sid"`       // 支付SID
	PaymentStatus bool   `gorm:"column:payment_status" json:"payment_status"` // 支付通道状态 true 为生产
	PaymentImage  string `gorm:"column:payment_image" json:"payment_image"`   // 支付图片或者logo
}
