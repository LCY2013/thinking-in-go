package limit

import (
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/gatefs"
	"testing"
)

func TestVfsLimit(t *testing.T) {
	fs := gatefs.New(vfs.OS("/path"), make(chan bool, 8))

	/*
		其中 vfs.OS("/path") 基于本地文件系统构造一个虚拟的文件系统，
		然后 gatefs.New 基于现有的虚拟文件系统构造一个并发受控的虚拟文件系统。
		并发数控制的原理在前面一节已经讲过，就是通过带缓存管道的发送和接收规则来实现最大并发阻塞：

		gatefs 对此做一个抽象类型 gate，增加了 enter 和 leave 方法分别对应并发代码的进入和离开。
		当超出并发数目限制的时候，enter 方法会阻塞直到并发数降下来为止。

		gatefs 包装的新的虚拟文件系统就是将需要控制并发的方法增加了 enter 和 leave 调用而已。
	*/

	fs.Lstat("")
}
