package test

import "testing"

func TestSub(t *testing.T) {
	res := Sub(3, 2)
	if res != 1 {
		t.Errorf("Sub(3, 2) = %d; want 1", res)
	}
}
