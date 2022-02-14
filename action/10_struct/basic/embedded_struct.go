package basic

// Person 基础用户
type Person struct {
	Name  string
	Phone string
	Addr  string
}

type Book struct {
	Title  string
	Author Person
}

type Inform struct {
	// Embedded Field
	Person
	Desc string
}
