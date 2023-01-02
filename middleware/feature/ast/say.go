package ast

import "fmt"

func Say(name string) string {
	return fmt.Sprintf("hello, %s", name)
}

func SayAll(firstName, lastName string, age int) string {
	return fmt.Sprintf("hello, %s %s, age %d", firstName, lastName, age)
}
