package model

type ProductImage struct {
	// 主键
	ID             int64  `gorm:"column:id;primary_key;not_null;auto_increment" json:"id"`
	ImageName      string `gorm:"column:image_name" json:"image_name"`
	ImageCode      string `gorm:"column:image_code;unique_index;not_null" json:"image_code"`
	ImageUrl       string `gorm:"column:image_url" json:"image_url"`
	ImageProductId int64  `gorm:"column:image_product_id" json:"image_product_id"`
}
