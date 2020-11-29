package pipeline

import "testing"

func TestIn(t *testing.T) {
	t.Run("In() should produce all the number in the parameters", func(t *testing.T) {
		out := In(2,3,4,5)
		count := 0
		for _ = range out {
			count ++
		}
		if count != 4 {
			t.Fatalf("In() does not return as expected, expected 4 numbers, got %d", count)
		}
	})
}
