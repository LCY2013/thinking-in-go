package types

import "fmt"

type User struct {
	Name string `json:"name"`

	// age 同属于一个包，age就能被测试用例访问
	// 如果不属于同一个包，age就不能被测试用例访问
	age int64 `json:"age"`
}

func (u User) GetAge() int64 {
	return u.age
}

func (u *User) ChangeName(newName string) {
	//if newName != "" {
	//	u.Name = newName
	//}
	u.Name = newName
}

func (u *User) GetName() string {
	return u.Name
}

func (u User) private() {
	fmt.Println("private")
}
