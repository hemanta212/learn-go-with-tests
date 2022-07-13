package main

import "testing"

func TestPerimeter(t *testing.T) {
	rectangle := Rectangle{10.0, 10.0}
	got := rectangle.Perimeter()
	want := 40.0
	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}
func TestArea(t *testing.T) {
	tests := []struct {
		name    string
		shape   Shape
		hasArea float64
	}{
		{name: "Rectangle Test", shape: Rectangle{10.0, 10.0}, hasArea: 100.0},
		{name: "Circle Test", shape: Circle{10.0}, hasArea: 314.1592653589793},
		{name: "Triangle Test", shape: Triangle{10.0, 4.0}, hasArea: 20.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.shape.Area()
			if got != tt.hasArea {
				t.Errorf("%s: %#v got %g want %g", tt.name, tt, got, tt.hasArea)
			}
		})
	}

}
