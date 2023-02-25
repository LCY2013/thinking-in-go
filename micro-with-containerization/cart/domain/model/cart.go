package model

type Cart struct {
	// 主键
	ID        int64 `gorm:"column:id;primary_key;not_null;auto_increment" json:"id"`
	ProductID int64 `gorm:"column:product_id;not_null" json:"product_id"`
	Num       int64 `gorm:"column:num;not_null" json:"num"`
	SizeID    int64 `gorm:"column:size_id;not_null" json:"size_id"`
	UserID    int64 `gorm:"column:user_id;not_null" json:"user_id"`
}
