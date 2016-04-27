package bulldozer

import (
	"fmt"
)

var MaxOutstanding = 1000

var sem = make(chan int, MaxOutstanding)

func process(r int) {
	fmt.Println(r)
}

func handle(r int) {
	sem <- 1   // Wait for active queue to drain.
	process(r) // May take a long time.
	<-sem      // Done; enable next request to run.
}

func Serve(queue chan int) {
	for {
		req := <-queue
		go handle(req) // Don't wait for handle to finish.
	}
}

func main() {
	queue := make(chan int)
	go Serve(queue)
	for i := 0; i < 10000000; i++ {
		queue <- i
	}
}
