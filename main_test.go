package main

import "testing"

func TestLoadData(t *testing.T) {
	X, Y, err := loadData()

	if err != nil {
		t.Errorf("loadData does not success, Error: %s", err)
	}

	if len(X) != 150 {
		t.Errorf("len(X) is not %d. got=%d", 150, len(X))
	}

	if Y[0] != "Setosa" {
		t.Errorf("Y[0] is not Setosa. got=%q", Y[0])
	}
}
