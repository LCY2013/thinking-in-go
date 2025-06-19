package main

import (
	"fmt"
	"regexp"
)

func main() {
	/*ch := make(chan struct{})
	go func() {
		for {
			select {
			case <-time.After(time.Second):
				fmt.Println("second")
			case <-ch:
				break
			}
		}
	}()

	time.Sleep(5 * time.Second)
	close(ch)*/

	compile := regexp.MustCompile(`<([^)]+)>`)

	fmt.Println(compile.FindString(
		`body:"*" custom:<kind:"*" path:"/echo/service/v1/example/echo/c">`))

	text := "My email is (example@example.com)"
	re := regexp.MustCompile(`\(([^)]+)\)`)
	matches := re.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		fmt.Println(match[1])
	}
}
