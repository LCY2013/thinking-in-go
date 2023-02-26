package model

import "time"

type Order struct {
	// 主键
	ID          int64         `gorm:"column:id;primary_key;not_null;auto_increment" json:"id"`
	OrderCode   string        `gorm:"column:order_code;unique_index;not_null" json:"order_code"`
	PayStatus   int32         `gorm:"column:pay_status" json:"pay_status"`
	ShipStatus  int32         `gorm:"column:ship_status" json:"ship_status"`
	Price       float64       `gorm:"column:price" json:"price"`
	OrderDetail []OrderDetail `gorm:"ForeignKey:OrderID" json:"order_detail"`
	CreateAt    time.Time     `gorm:"column:create_at" json:"create_at"`
	UpdateAt    time.Time     `gorm:"column:update_at" json:"update_at"`
}
