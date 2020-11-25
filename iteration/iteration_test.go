package iteration

import (
	"fmt"
	"testing"
)

func TestIteration(t *testing.T) {

	t.Run("Should return correct value for first run", func(t *testing.T) {
		repeated := Repeat("a")
		expected := "aaaaa"

		assertCorrectMessage(t, repeated, expected)
	})

	t.Run("Should return correct value for any character", func(t *testing.T) {
		repeated := Repeat("b")
		expected := "bbbbb"

		assertCorrectMessage(t, repeated, expected)
	})
}

func BenchmarkRepeat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Repeat("a")
	}
}

func assertCorrectMessage(t *testing.T, repeated string, expected string) {
	t.Helper()
	if repeated != expected {
		t.Errorf("expected %q but got %q", expected, repeated)
	}
}

func ExampleRepeat() {
	repeated := Repeat("a")
	fmt.Println(repeated)
	// Output: aaaaa
}
