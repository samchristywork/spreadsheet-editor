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

func TestGetValue(t *testing.T) {
	setCellContent(0, 0, "1")

	setCellContent(0, 1, "2")

	setCellContent(0, 2, "=A0+B0")

	content, err := getCellValue(0, 0)
	if err != nil {
		t.Errorf("%v", err)
	}

	assertEqual(t, content, "1")

	content, err = getCellValue(0, 2)
	if err != nil {
		t.Errorf("%v", err)
	}

	assertEqual(t, content, "3")
}

func TestEval(t *testing.T) {
	setCellContent(0, 0, "1")

	setCellContent(0, 1, "2")

	setCellContent(0, 2, "=A0+B0")

	assertEqual(t, eval("1+1"), "2")

	assertEqual(t, eval("sum(1,1)"), "2")

	assertEqual(t, eval("sum(1,1+1)"), "3")

	assertEqual(t, eval("strlen(\"foo\")"), "3")

	assertEqual(t, eval("A0"), "1")

	assertEqual(t, eval("B0"), "2")

	assertEqual(t, eval("A0+B0"), "3")

	assertEqual(t, eval("C0"), "3")

	assertEqual(t, eval("sum(\"A0:C0\")"), "6")

	assertEqual(t, eval("sum('A0:C0')"), "6")

	assertEqual(t, eval("A0+D0"), "Error applying function: Cell D0 is empty")

	assertEqual(t, eval("asdf"), "Error applying function: asdf is not a valid cell identifier")
}

func assertRange(t *testing.T, a []string, b []string) {
	if len(a) != len(b) {
		t.Errorf("len(a) = %d != len(b) = %d", len(a), len(b))
	}

	for i := 0; i < len(a); i++ {
		assertEqual(t, a[i], b[i])
	}
}

func TestRange(t *testing.T) {
	r, err := getCellRange(0, 0, 3, 0)
	if err != nil {
		t.Errorf("%v", err)
	}

	assertRange(t, r, []string{"A0", "B0", "C0", "D0"})

	r, err = getCellRange(1, 2, 1, 2)
	if err != nil {
		t.Errorf("%v", err)
	}

	assertRange(t, r, []string{"B2"})

	r, err = getCellRange(1, 2, 1, 5)
	if err != nil {
		t.Errorf("%v", err)
	}

	assertRange(t, r, []string{"B2", "B3", "B4", "B5"})
}
