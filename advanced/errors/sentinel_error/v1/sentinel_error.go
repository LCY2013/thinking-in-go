package v1

import "fmt"

type MyError struct {
	Msg  string
	File string
	Line int
}

func (e *MyError) Error() string {
	return fmt.Sprintf("%s:%s:%s", e.File, e.File, e.Line)
}

func Test() error {
	return &MyError{"Something happened", "server.go", 12}
}
