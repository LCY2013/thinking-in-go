package main

import v3 "fufeng.org/standard/logs/v3"

func main() {
	v3.Trace.Println("I have something standard to say")
	v3.Info.Println("Special Information")
	v3.Warning.Println("There is something you need to known about")
	v3.Error.Println("Something has failed")
}
