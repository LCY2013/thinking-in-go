package types

import "fmt"

type User struct {
	Name string `json:"name"`

	// age 同属于一个包，age就能被测试用例访问
	// 如果不属于同一个包，age就不能被测试用例访问
	age int `json:"age"`
}

func (u User) GetAge() int {
	return u.age
}

func (u *User) ChangeName(newName string) {
	u.Name = newName
}

func (u User) private() {
	fmt.Println("private")
}
