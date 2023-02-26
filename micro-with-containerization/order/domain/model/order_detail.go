package model

type OrderDetail struct {
	ID            int64 `gorm:"column:id;primary_key;not_null;auto_increment" json:"id"`
	ProductID     int64 `gorm:"column:product_id" json:"product_id"`
	ProductNum    int64 `gorm:"column:product_num" json:"product_num"`
	ProductSizeID int64 `gorm:"column:product_size_id" json:"product_size_id"`
	ProductPrice  int64 `gorm:"column:product_price" json:"product_price"`
	OrderID       int64 `gorm:"column:order_id" json:"order_id"`
}
