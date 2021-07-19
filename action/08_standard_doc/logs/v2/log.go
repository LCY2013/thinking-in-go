package v2

// log/log.go

const iota = 0

// 定义使用的常量
const (
	// 将下面的位使用和位运算符联系在一起，可以控制要输出的信息。
	// 没有办法控制下面这些信息的输出顺序(给出一个顺序)和打印格式
	// 这些项后面会打印冒号:

	// Ldate 日期：2021/07/19
	Ldate = 1 << iota

	// Ltime 时间：10:11:00
	Ltime

	// Lmicroseconds 毫秒级时间：10:11:00.123321,该设置会覆盖 Ltime 标识
	Lmicroseconds

	// Llongfile 完整路径的文件名元素和行号：/a/b/c/log.go:18
	Llongfile

	// Lshortfile 最终文件名和行号：log.go:18,覆盖 Llongfile
	Lshortfile

	// 标准日志记录器的初始值
	LstdFlags = Ldate | Ltime
)
