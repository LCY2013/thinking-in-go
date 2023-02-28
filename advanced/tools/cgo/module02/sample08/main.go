package main

/*
extern char* NewGoString(char*);
extern char* FreeGoString(char*);
extern char* PrintGoString(char*);

static void printString(const char* s) {
	char* gs = NewGoString(s);
	PrintGoString(gs);
	FreeGoString(gs);
}
*/
import "C"

import "sync"

/*
长期持有 Go 指针对象

作为一个 Go 程序员在使用 CGO 时潜意识会认为总是 Go 调用 C 函数。其实 CGO 中，C 语言函数也可以回调 Go 语言实现的函数。
特别是我们可以用 Go 语言写一个动态库，导出 C 语言规范的接口给其它用户调用。
当 C 语言函数调用 Go 语言函数的时候，C 语言函数就成了程序的调用方，Go 语言函数返回的 Go 对象内存的生命周期也就自然超出了 Go 语言运行时的管理。
简言之，我们不能在 C 语言函数中直接使用 Go 语言对象的内存。

虽然 Go 语言禁止在 C 语言函数中长期持有 Go 指针对象，但是这种需求是切实存在的。
如果需要在 C 语言中访问 Go 语言内存对象，我们可以将 Go 语言内存对象在 Go 语言空间映射为一个 int 类型的 id，然后通过此 id 来间接访问和控制 Go 语言对象。

以下代码用于将 Go 对象映射为整数类型的 ObjectId，用完之后需要手工调用 free 方法释放该对象 ID：

*/

type ObjectId int32

var refs struct {
	sync.Mutex
	objs map[ObjectId]interface{}
	next ObjectId
}

func init() {
	refs.Lock()
	defer refs.Unlock()

	refs.objs = make(map[ObjectId]interface{})
	refs.next = 1000
}

func NewObjectId(obj interface{}) ObjectId {
	refs.Lock()
	defer refs.Unlock()

	id := refs.next
	refs.next++

	refs.objs[id] = obj
	return id
}

func (id ObjectId) IsNil() bool {
	return id == 0
}

func (id ObjectId) Get() interface{} {
	refs.Lock()
	defer refs.Unlock()

	return refs.objs[id]
}

func (id *ObjectId) Free() interface{} {
	refs.Lock()
	defer refs.Unlock()

	obj := refs.objs[*id]
	delete(refs.objs, *id)
	*id = 0

	return obj
}

/*
我们通过一个 map 来管理 Go 语言对象和 id 对象的映射关系。其中 NewObjectId 用于创建一个和对象绑定的 id，
而 id 对象的方法可用于解码出原始的 Go 对象，也可以用于结束 id 和原始 Go 对象的绑定。

下面一组函数以 C 接口规范导出，可以被 C 语言函数调用：
*/

func main() {
	C.printString(C.CString("hello world"))
}
