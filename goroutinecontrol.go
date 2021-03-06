package main

import (
	"fmt"
	"time"
)

type WorkerChannel struct {
	Receiver chan int
	Finisher chan int
}

func NewWorkerChannel(goRoutineCount int) *WorkerChannel {

	return &WorkerChannel{
		Receiver: make(chan int, goRoutineCount),
		Finisher: make(chan int),
	}
}

// Send the channel to resp to.
// id is there just for making the logs better of.
func worker(myChan *WorkerChannel, freeChan chan *WorkerChannel, id int) {
	go func() {
		for {
			select {
			case data := <-myChan.Receiver:
				// Processing task
				fmt.Println("Channel id ", id, " ", data)
				// Done let me ask for more work.
				freeChan <- myChan
			case <-myChan.Finisher:
				close(myChan.Receiver)
				close(myChan.Finisher)
				return
			}
		}

	}()
}

func initializeWorkers(workerCount int) chan *WorkerChannel {

	freeWorkerChan := make(chan *WorkerChannel, workerCount)
	func() {
		for i := 0; i < workerCount; i++ {

			workerChan := NewWorkerChannel(workerCount)
			worker(workerChan, freeWorkerChan, i)
			// Everyone is free right now. Ask for some work please !
			freeWorkerChan <- workerChan
		}
	}()
	return freeWorkerChan
}

// Scheduler returns the pipe send data on.
// @args - the workers that are free.
func scheduler(freeWorkerChan chan *WorkerChannel, exitChan chan int, workerCount int) (pipe chan int, finish chan int) {

	pipe = make(chan int, workerCount)
	finish = make(chan int)
	go func() {
		for {
			// pickData only if someone is free.
			freeChan := <-freeWorkerChan
			select {
			case data := <-pipe:
				// Assigning the args for work.
				freeChan.Receiver <- data
			case <-finish:
				// Make sure all the workerChan's are done
				freeChan.Finisher <- 1
				fmt.Println("Closed i")
				for i := 0; i < workerCount-1; i++ {
					freeChan := <-freeWorkerChan
					freeChan.Finisher <- 1
					fmt.Println("Closed ", i)
				}
				close(pipe)
				close(finish)
				exitChan <- 1
				return
			}
		}
	}()
	return
}

func main() {
	exitChan := make(chan int)
	workerCount := 100
	freeWorkerChan := initializeWorkers(workerCount)

	pipe, finish := scheduler(freeWorkerChan, exitChan, workerCount)

	// Pumping data in for the worker
	for i := 0; i < 100000; i++ {
		time.Sleep(1 * time.Millisecond)
		pipe <- i
	}

	time.Sleep(5 * time.Second)

	for i := 0; i < 100000; i++ {
		pipe <- i
	}

	finish <- 1
	//	finish <- 1
	<-exitChan
}
