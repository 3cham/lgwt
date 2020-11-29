package pipeline

import (
	"sync"
	"time"
)

func gen(nums ...int) <-chan int {
	out := make(chan int)

	go func() {
		for _, num := range (nums) {
			out <- num
		}
		close(out)
	}()
	return out
}

func sq(input <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		for num := range input {
			out <- sleepAndSquare(num)
		}
		close(out)
	}()
	return out
}

func merge(chans ...<-chan int) <-chan int {
	var wg = sync.WaitGroup{}

	result := make(chan int)

	getNum := func(c <-chan int){
		for n := range c{
			result <- n
		}
		wg.Done()
	}
	wg.Add(len(chans))
	for _, c := range chans {
		go getNum(c)
	}

	go func() {
		wg.Wait()
		close(result)
	}()
	return result
}

func sleepAndSquare(num int) int {
	if num % 2 == 1 {
		time.Sleep(1 * time.Second)
	}
	return num * num
}