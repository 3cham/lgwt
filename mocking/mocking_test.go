package main

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

type CountDownOperationSpy struct {
	Calls []string
}

func (c *CountDownOperationSpy) Sleep() {
	c.Calls = append(c.Calls, sleep)
}

func (c *CountDownOperationSpy) Write(message []byte) (n int, err error) {
	c.Calls = append(c.Calls, write)

	return
}

type SpyTime struct {
	sleptDuration time.Duration
}

func (s *SpyTime) Sleep(duration time.Duration) {
	s.sleptDuration = duration
}

func assertCorrectString(t *testing.T, got, expected string) {
	t.Helper()
	if got != expected {
		t.Fatalf("Wrong string: got %s, expected %s", got, expected)
	}
}

func TestCountDown(t *testing.T) {
	t.Run("CountDown() should print correctly", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		sleeper := &CountDownOperationSpy{}

		CountDown(buffer, sleeper)

		got := buffer.String()
		expected := "3\n2\n1\nGo!"

		assertCorrectString(t, got, expected)
	})

	t.Run("Sleep before each print", func(t *testing.T) {
		writerAndSleeper := &CountDownOperationSpy{}

		CountDown(writerAndSleeper, writerAndSleeper)

		expected := []string{
			sleep,
			write,
			sleep,
			write,
			sleep,
			write,
			sleep,
			write,
		}

		got := writerAndSleeper.Calls

		if !reflect.DeepEqual(got, expected) {
			t.Fatalf("Wrong order of calling Sleep and Write, got %v, expected %v", got, expected)
		}
	})
}

func TestConfigurableSleeper(t *testing.T) {
	t.Run("Test configurable Sleeper", func(t *testing.T) {
		sleepTime := 5 * time.Second
		spyTime := &SpyTime{}

		sleeper := ConfigurableSleeper{sleepTime, spyTime.Sleep}
		sleeper.Sleep()

		if spyTime.sleptDuration != sleepTime {
			t.Fatalf("Sleep duration is incorrect, got %v, expected %v", spyTime.sleptDuration, sleepTime)
		}
	})
}
