package pipeline

func In(nums ...int) <-chan int {
	out := make(chan int)

	go func() {
		for _, num := range (nums) {
			out <- num
		}
		close(out)
	}()
	return out
}