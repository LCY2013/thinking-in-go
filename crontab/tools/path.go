package tools

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Pwd() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/master/configs", pwd), nil
}

func PwdSlave() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/slave/configs", pwd), nil
}

func GetCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return "", errors.New(`error: Can't find "/" or "\"`)
	}
	return string(path[0 : i+1]), nil
}

func GetAllCurrentDir() []string {
	var dirs []string
	pwd, _ := os.Getwd()

	//获取当前目录下的所有文件或目录信息
	filepath.Walk(pwd, func(path string, info os.FileInfo, err error) error {
		//fmt.Println(path) //打印path信息
		//fmt.Println(info.Name()) //打印文件或目录名
		dirs = append(dirs, path)
		return nil
	})
	return dirs
}

func GetCurrentDir() []string {
	var dirs []string
	pwd, _ := os.Getwd()

	//获取当前目录下的文件或目录名(包含路径)
	filepathNames, err := filepath.Glob(filepath.Join(pwd, "*"))
	if err != nil {
		log.Fatal(err)
	}

	for i := range filepathNames {
		//fmt.Println(filepathNames[i]) //打印path
		dirs = append(dirs, filepathNames[i])
	}

	return dirs
}
