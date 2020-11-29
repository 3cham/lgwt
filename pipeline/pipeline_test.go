package pipeline

import (
	"reflect"
	"sort"
	"testing"
)

func makeSortedArrFromChan(c <-chan int) []int {
	result := []int{}
	for num := range c {
		result = append(result, num)
	}
	sort.Ints(result)
	return result
}

func TestGen(t *testing.T) {
	t.Run("gen() should produce all the number in the parameters", func(t *testing.T) {
		out := gen(2, 3, 4, 5)
		count := 0
		for _ = range out {
			count++
		}
		if count != 4 {
			t.Fatalf("gen() does not return as expected, expected 4 numbers, got %d", count)
		}
	})
}

func TestSq(t *testing.T) {
	t.Run("sq() should return squared numbers", func(t *testing.T) {
		in := gen(1, 2, 3, 4, 5, 6)
		out1 := sq(in)
		out2 := sq(in)
		out3 := sq(in)
		out4 := sq(in)

		out := merge(out1, out2, out3, out4)
		sortedOut := makeSortedArrFromChan(out)
		var exected = makeSortedArrFromChan(gen(1, 4, 9, 16, 25, 36))
		if !reflect.DeepEqual(sortedOut, exected) {
			t.Fatalf("gen() does not return expected results, expected %v, got %v", exected, sortedOut)
		}
	})
}
