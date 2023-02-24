package model

type ProductSize struct {
	// 主键
	ID            int64  `gorm:"column:id;primary_key;not_null;auto_increment" json:"id"`
	SizeName      string `gorm:"column:size_name" json:"size_name"`
	SizeCode      string `gorm:"column:size_code" json:"size_code"`
	SizeProductID int64  `gorm:"column:size_product_id" json:"size_product_id"`
}
