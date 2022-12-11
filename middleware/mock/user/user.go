package user

import "github.com/LCY2013/thinking-in-go/middleware/mock/person"

type User struct {
	Person person.Male
}

func NewUser(p person.Male) *User {
	return &User{Person: p}
}

func (u *User) Get(id int64) error {
	return u.Person.Get(id)
}
