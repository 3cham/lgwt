package pipeline

import (
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

func sleepAndSquare(num int) int {
	if num % 2 == 1 {
		time.Sleep(1 * time.Second)
	}
	return num * num
}