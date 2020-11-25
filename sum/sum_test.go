package sum

import (
	"fmt"
	"reflect"
	"testing"
)

func assertCorrectMessageForSum(t *testing.T, got int, expected int) {
	t.Helper()
	if expected != got {
		t.Errorf("Got %d, expected %d", got, expected)
	}
}

func assertCorrectMessageForSumAll(t *testing.T, got []int, expected []int) {
	t.Helper()

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Got %v, expected %v", got, expected)
	}
}

func TestSum(t *testing.T)  {
	t.Run("Sum should return correct sum", func(t *testing.T) {
		numbers := []int{1,2,3,4,5}

		got := Sum(numbers)
		expected := 15

		assertCorrectMessageForSum(t, got, expected)
	})

	t.Run("Sum should return correct sum for another array", func(t *testing.T) {
		numbers := []int{1,2,3,4}
		got := Sum(numbers)
		expected := 10
		assertCorrectMessageForSum(t, got, expected)
	})
}

func TestSumAll(t *testing.T) {
	t.Run("SumAll returns array of sums", func(t *testing.T) {
		got := SumAll([]int{1,2}, []int{1,2,3,4})
		expected := []int{3, 10}

		assertCorrectMessageForSumAll(t, got, expected)
	})
}

func TestSumTail(t *testing.T) {
	t.Run("SumAll returns array of sums", func(t *testing.T) {
		got := SumAllTail([]int{1,2}, []int{1,2,3,4})
		expected := []int{2, 9}

		assertCorrectMessageForSumAll(t, got, expected)
	})
}

func ExampleSum() {
	sum := Sum([]int{1,2,3,3})
	fmt.Println(sum)
	// Output: 9
}

func ExampleSumAll() {
	sumAll := SumAll([]int{1,2,2}, []int{1,2,3,4,4})
	fmt.Println(sumAll)
	// Output: [5 14]
}

func ExampleSumAllTail() {
	sumAll := SumAllTail([]int{1,2,2}, []int{1,2,3,4,4}, []int{})
	fmt.Println(sumAll)
	// Output: [4 13 0]
}