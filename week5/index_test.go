package main

import (
	"testing"
)

func TestGetVal(t *testing.T) {
	Check(FindGipotenusa(3, 4), 5, t)
	Check(FindGipotenusa(1, 1), 1.4142135623730951, t)
}

func Check(ret, expected float64, t *testing.T) {
	if ret != expected {
		t.Error("Expected: ", expected, "Got: ", ret)
	} else {
		t.Log("SUCCESS")
	}
}
