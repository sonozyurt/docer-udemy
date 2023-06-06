package main

import "testing"

func TestAdd(t *testing.T) {
	x := add(5, 5)
	want := 10
	if x != want {
		t.Fatalf("failed got %d ,wanted %d", x, want)
	}

}
