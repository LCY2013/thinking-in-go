package model

type User struct {
	// 主键
	ID int64 `gorm:"column:id;primary_key;not_null;auto_increment"`
	// 用户名称
	UserName string `gorm:"column:user_name;unique_index;not_null"`
	// 添加需要的字段
	FirstName string `gorm:"column:first_name"`
	// 加密后密码
	HashPassword string `gorm:"column:hash_password"`
}
