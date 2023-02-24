package model

type Product struct {
	// 主键
	ID           int64          `gorm:"column:id;primary_key;not_null;auto_increment" json:"id"`
	ProductName  string         `gorm:"column:product_name" json:"product_name"`
	ProductSku   string         `gorm:"column:product_sku" json:"product_sku"`
	ProductPrice float64        `gorm:"column:product_price" json:"product_price"`
	ProductImage []ProductImage `gorm:"ForeignKey:ImageProductID" json:"product_image"`
	ProductSize  []ProductSize  `gorm:"ForeignKey:SizeProductID" json:"product_size"`
	ProductSeo   ProductSeo     `gorm:"ForeignKey:SeoProductID" json:"product_seo"`
}
