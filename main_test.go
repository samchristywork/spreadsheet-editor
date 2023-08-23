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

	assertEqual(t, row("GT41"), 41)
}

func TestFixedWidth(t *testing.T) {
	assertEqual(t, fixedWidth("foo", 5), "foo  ")

	assertEqual(t, fixedWidth("foo", 2), "f…")

	assertEqual(t, fixedWidth("foo", 0), "")

	assertEqual(t, fixedWidth("", 5), "     ")

	assertEqual(t, fixedWidth("fooBarBaz", 5), "fooB…")
}

func TestGetColumnName(t *testing.T) {
	assertEqual(t, getColumnName(0), "A")

	assertEqual(t, getColumnName(1), "B")

	assertEqual(t, getColumnName(2), "C")

	assertEqual(t, getColumnName(26), "AA")

	assertEqual(t, getColumnName(54), "BC")
}

func TestCellContent(t *testing.T) {
	setCellContent(1, 1, "123")

	content, err := getCellContent(1, 1)
	if err != nil {
		t.Errorf("%v", err)
	}

	assertEqual(t, content, "123")
}

