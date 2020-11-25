package dependency

import (
	"bytes"
	"fmt"
	"testing"
)

func TestGreet(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "", args: args{name: "world"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := bytes.Buffer{}
			Greet(&writer, tt.args.name)

			got := writer.String()
			expected := "Hello, world"

			assertCorrectMessage(t, got, expected)
		})
	}
}

func assertCorrectMessage(t *testing.T, got string, expected string) {
	t.Helper()
	if got != expected {
		t.Errorf("expected %q but got %q", expected, got)
	}
}

func ExampleGreet() {
	writer := bytes.Buffer{}
	Greet(&writer, "world")
	fmt.Println(writer.String())
	// Output: Hello, world
}
