package main

//接口的静态特性与动态特性
//接口的静态特性体现在接口类型变量具有静态类型，比如var err error中变量 err 的静态类型为 error。拥有静态类型，那就意味着编译器会在编译阶段对所有接口类型变量的赋值操作进行类型检查，编译器会检查右值的类型是否实现了该接口方法集合中的所有方法。如果不满足，就会报错。
//接口的动态特性，就体现在接口类型变量在运行时还存储了右值的真实类型信息，这个右值的真实类型被称为接口类型变量的动态类型。
//
//与动态语言不同的是，Go 接口还可以保证“动态特性”使用时的安全性。
//比如，编译器在编译期就可以捕捉到将 int 类型变量传给 QuackableAnimal 接口类型变量这样的明显错误，决不会让这样的错误遗漏到运行时才被发现。

type QuackableAnimal interface {
	Quack()
}
type Duck struct{}

func (Duck) Quack() {
	println("duck quack!")
}

type Dog struct{}

func (Dog) Quack() {
	println("dog quack!")
}

type Bird struct{}

func (Bird) Quack() {
	println("bird quack!")
}

func AnimalQuackInForest(a QuackableAnimal) {
	a.Quack()
}

func main() {
	animals := []QuackableAnimal{new(Duck), new(Dog), new(Bird)}
	for _, animal := range animals {
		AnimalQuackInForest(animal)
	}
}
