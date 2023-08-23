package main

import (
	"testing"
)

func assertEqual(t *testing.T, a, b interface{}) {
	if a != b {
		t.Errorf("\"%v\" != \"%v\"", a, b)
	}
}

func TestCol(t *testing.T) {
	assertEqual(t, column("A0"), 0)

	c, _ := splitColumnRow("A0")
	assertEqual(t, c, "A")

	c, _ = splitColumnRow("G4")
	assertEqual(t, c, "G")

	c, _ = splitColumnRow("GT41")
	assertEqual(t, c, "GT")

	assertEqual(t, column("C9"), 2)

	assertEqual(t, column("BC9"), 54)
}

func TestRow(t *testing.T) {
	assertEqual(t, row("A0"), 0)

	_, r := splitColumnRow("A0")
	assertEqual(t, r, "0")

	_, r = splitColumnRow("G4")
	assertEqual(t, r, "4")

	_, r = splitColumnRow("GT41")
	assertEqual(t, r, "41")

	assertEqual(t, row("C9"), 9)
}
