package model

type ProductSeo struct {
	// 主键
	ID             int64  `gorm:"column:id;primary_key;not_null;auto_increment" json:"id"`
	SeoTitle       string `gorm:"column:seo_title" json:"seo_title"`
	SeoKeywords    string `gorm:"column:seo_keywords" json:"seo_keywords"`
	SeoDescription string `gorm:"column:seo_description" json:"seo_description"`
	SeoCode        string `gorm:"column:seo_code" json:"seo_code"`
	SeoProductID   int64  `gorm:"column:seo_product_id" json:"seo_product_id"`
}
