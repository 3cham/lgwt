package pipeline

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
			out <- num * num
		}
		close(out)
	}()
	return out
}