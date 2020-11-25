package structs

import (
	"fmt"
	"testing"
)

func checkShapePerimeterValueEqual(t *testing.T, shape Shape, expected float64) {
	t.Helper()
	got := shape.Perimeter()
	if got != expected {
		t.Errorf("Wrong perimeter value for shape %#v: got %g expected %g", shape, got, expected)
	}
}

func checkShapeAreaValueEqual(t *testing.T, shape Shape, expected float64) {
	t.Helper()
	got := shape.Area()
	if got != expected {
		t.Errorf("Wrong area value for shape %#v: got %g expected %g", shape, got, expected)
	}
}

func TestPerimeter(t *testing.T) {
	var perimeterTests = []struct {
		shape    Shape
		expected float64
	}{
		{Rectangle{10.0, 10.0}, 40.0},
		{Circle{10.0}, 31.41592653589793 * 2},
	}
	for _, tt := range perimeterTests {
		checkShapePerimeterValueEqual(t, tt.shape, tt.expected)
	}
}

func TestArea(t *testing.T) {
	var areaTests = []struct {
		shape    Shape
		expected float64
	}{
		{shape: Rectangle{10.0, 10.0}, expected: 100.0},
		{shape: Circle{10.0}, expected: 314.1592653589793},
		{shape: Triangle{12, 6}, expected: 36.0},
	}
	for _, tt := range areaTests {
		msg := fmt.Sprintf("Test Area() for %v", tt)
		t.Run(msg, func(t *testing.T) {
			checkShapeAreaValueEqual(t, tt.shape, tt.expected)
		})
	}
}
