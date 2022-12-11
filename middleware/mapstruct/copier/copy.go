package copier

// Copier 复制数据
// 1、深拷贝或者浅拷贝
// 2、src 和 dest 支持map和结构体
// 3、只复制结构体的公共字段
// 4、这些带有*的设计，会存在内存逃逸现象
type Copier[Src, Dest any] interface {
	// CopyTo 将src中的数据复制到dest中
	CopyTo(src *Src, dest *Dest) error
	// Copy 将创建一个Dest的实例，并将 src中的数据复制过去
	Copy(src *Src) (*Dest, error)
}
