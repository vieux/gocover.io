package test

import "testing"

func TestAdd(t *testing.T) {
	res := Add(2, 2)
	if res != 4 {
		t.Errorf("Add(2, 2) = %d; want 4", res)
	}
}
