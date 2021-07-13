package structs

import "fmt"

// 声明用户类型的变量
var bill User

func UserFunc() {
	// 声明user类型的变量，并初始化所有字段
	lisa := User{
		name:       "Lisa",
		email:      "lisa@email.com",
		ext:        123,
		privileged: true,
	}
	fmt.Println(lisa)

	// 不使用字段名称，创建结构体类型的值
	bob := User{"Bob", "bob@email.com", 123, true}
	fmt.Println(bob)

	// 声明Admin类型的bianl
	fred := Admin{
		person: User{
			name:       "Lisa",
			email:      "lisa@email.com",
			ext:        123,
			privileged: true,
		},
		level: "super",
	}
	fmt.Println(fred)
}

func DurationFun() {
	var dur Duration
	// # fufeng.org/sample05/structs
	//structs/options.go:37:6: cannot use int64(1000) (type int64) as type Duration in assignment
	// dur = int64(1000)
	fmt.Println(dur)
}

// NewUser 创建一个用户信息
func NewUser(name, email string) User {
	return User{
		name:  name,
		email: email,
	}
}

// CreateUser 创建一个用户信息
func CreateUser(name, email string) *User {
	return &User{
		name:  name,
		email: email,
	}
}

// Notify 使用值接收者实现了一个方法,这里会使用该值接收者的一个副本来执行
func (user User) Notify() {
	fmt.Printf("Sending User Email To %s<%s>\n",
		user.name,
		user.email)
}

// ChangeEmail 使用指针接收者实现的一个方法
func (user *User) ChangeEmail(email string) {
	user.email = email
}
