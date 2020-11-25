package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

const (
	countDownStart = 3
	lastString     = "Go!"
	write          = "write"
	sleep          = "sleep"
)

type Sleeper interface {
	Sleep()
}

type ConfigurableSleeper struct {
	duration time.Duration
	sleep    func(time.Duration)
}

func (c *ConfigurableSleeper) Sleep() {
	c.sleep(c.duration)
}

func CountDown(writer io.Writer, sleeper Sleeper) {
	for i := countDownStart; i > 0; i-- {
		sleeper.Sleep()
		fmt.Fprintf(writer, "%d\n", i)
	}

	sleeper.Sleep()
	fmt.Fprintf(writer, "%s", lastString)
}

func main() {
	sleeper := ConfigurableSleeper{
		duration: 1 * time.Second,
		sleep:    time.Sleep,
	}
	CountDown(os.Stdout, &sleeper)
}
