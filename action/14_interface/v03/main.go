package main

// 尽量定义“小接口”
// 接口类型的背后，是通过把类型的行为抽象成契约，建立双方共同遵守的约定，这种契约将双方的耦合降到了最低的程度。
// 和生活工作中的契约有繁有简，签署方式多样一样，代码间的契约也有多有少，有大有小，而且达成契约的方式也有所不同。
// 而 Go 选择了去繁就简的形式，这主要体现在以下两点上：
// 1、隐式契约，无需签署，自动生效
// Go 语言中接口类型与它的实现者之间的关系是隐式的，不需要像其他语言（比如 Java）那样要求实现者显式放置“implements”进行修饰，实现者只需要实现接口方法集合中的全部方法便算是遵守了契约，并立即生效了。
// 2、更倾向于“小契约”
// Go 选择了使用“小契约”，表现在代码上就是尽量定义小接口，即方法个数在 1~3 个之间的接口。Go 语言之父 Rob Pike 曾说过的“接口越大，抽象程度越弱”，这也是 Go 社区倾向定义小接口的另外一种表述。

// 小接口有哪些优势？
// 1、第一点：接口越小，抽象程度越高

// Flyable 会飞的
type Flyable interface {
	Fly()
}

// Swimable 会游泳的
type Swimable interface {
	Swim()
}

// FlySwimable 会飞且会游泳的
type FlySwimable interface {
	Flyable
	Swimable
}

// 2、第二点：小接口易于实现和测试
// 小接口拥有比较少的方法，一般情况下只有一个方法。所以要想满足这一接口，我们只需要实现一个方法或者少数几个方法就可以了，这显然要比实现拥有较多方法的接口要容易得多。
// 尤其是在单元测试环节，构建类型去实现只有少量方法的接口要比实现拥有较多方法的接口付出的劳动要少许多。

// 3、第三点：小接口表示的“契约”职责单一，易于复用组合
// Go 推崇通过组合的方式构建程序。
// Go 开发人员一般会尝试通过嵌入其他已有接口类型的方式来构建新接口类型，就像通过嵌入 io.Reader 和 io.Writer 构建 io.ReadWriter 那样。

//定义小接口，可以遵循的几点？
//1、首先，别管接口大小，先抽象出接口。
// Go 语言还比较年轻，它的设计哲学和推崇的编程理念可能还没被广大 Gopher 100% 理解、接纳和应用于实践当中，尤其是 Go 所推崇的基于接口的组合思想。
// 尽管接口不是 Go 独有的，但专注于接口是编写强大而灵活的 Go 代码的关键。因此，在定义小接口之前，我们需要先针对问题领域进行深入理解，聚焦抽象并发现接口，就像下图所展示的那样，先针对领域对象的行为进行抽象，形成一个接口集合。
//
//2、第二，将大接口拆分为小接口。
// 有了接口后，我们就会看到接口被用在了代码的各个地方。
// 一段时间后，我们就来分析哪些场合使用了接口的哪些方法，是否可以将这些场合使用的接口的方法提取出来，放入一个新的小接口中。
//
//3、最后，我们要注意接口的单一契约职责。
//那么，上面已经被拆分成的小接口是否需要进一步拆分，直至每个接口都只有一个方法呢？
//这个依然没有标准答案，不过你依然可以考量一下现有小接口是否需要满足单一契约职责，就像 io.Reader 那样。如果需要，就可以进一步拆分，提升抽象程度。

func main() {

}
