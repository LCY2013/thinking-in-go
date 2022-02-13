package main

func main() {

}

type Resource string

func Poller(in, out chan *Resource) {
	for r := range in {
		// do poll -> poll url
		out <- r
	}
}
