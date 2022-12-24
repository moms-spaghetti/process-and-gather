package main

import (
	"log"
	"sync"
	"time"
)

func processAndGather(in <-chan int, processor func(int) int, n int) []int {
	out := make(chan int, n)

	wg := new(sync.WaitGroup)
	wg.Add(n)

	for v := range in {
		go func(v int) {
			defer wg.Done()
			out <- processor(v)
		}(v)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	res := make([]int, n)
	i := 0
	for v := range out {
		res[i] = v
		i++
	}

	return res
}

func processor(n int) int {
	time.Sleep(1 * time.Second)
	return n * 2
}

func newIn(n int) <-chan int {
	ch := make(chan int, n)
	for i := 0; i < n; i++ {
		ch <- i
	}

	close(ch)

	return ch
}

func main() {
	n := 5
	in := newIn(n)
	r := processAndGather(in, processor, n)
	log.Print("r: ", r)

}
