package structs

// User 在程序里定义一个用户类型
type User struct {
	// 用户姓名
	name string
	// 用户邮箱
	email string
	// 用户退出码
	ext int
	// 是否是有特权
	privileged bool
}

// Admin 需要一个user作为管理者，并附加权限
type Admin struct {
	person User
	level  string
}

// Duration 声明一个基于int64的新类型
type Duration int64
