package errors

import (
	"io"
	"log"
	"os"
	"runtime"
	"syscall"
)

func syscallErr() {
	err := syscall.Chmod(":invalid path:", 0666)
	if err != nil {
		log.Fatal(err.(syscall.Errno))
	}
}

func copyFile(src, dst string) (int64, error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer func(srcFile *os.File) {
		err = srcFile.Close()
		if err != nil {
			log.Println(err)
		}
	}(srcFile)

	dstFile, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer func(dstFile *os.File) {
		err = dstFile.Close()
		if err != nil {
			log.Println(err)
		}
	}(dstFile)

	return io.Copy(dstFile, srcFile)
}

func handleErr() {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case runtime.Error:
				// 这是运行时错误类型异常
				log.Println(x)
			case error:
				// 普通错误类型异常
			default:
				// 其他类型异常
			}
		}
	}()

}
