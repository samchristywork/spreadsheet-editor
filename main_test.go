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

func TestIsCellIdentifier(t *testing.T) {
	assertEqual(t, isCellIdentifier("A0"), true)
	assertEqual(t, isCellIdentifier("B5"), true)
	assertEqual(t, isCellIdentifier("AA10"), true)
	assertEqual(t, isCellIdentifier("foo"), false)
	assertEqual(t, isCellIdentifier("123"), false)
	assertEqual(t, isCellIdentifier(""), false)
}

// Tests that handleIncrement uses the row/col parameters rather than currentCell.
func TestHandleIncrementUsesParams(t *testing.T) {
	setCellContent(0, 0, "=A0+B0")
	currentCell[0] = 5
	currentCell[1] = 5

	handleIncrement("=A0+B0", 0, 0, 1, 0)

	content, _ := getCellContent(1, 0)
	assertEqual(t, content, "=A1+B1")

	// Should NOT have written to currentCell + delta
	content, _ = getCellContent(6, 5)
	assertEqual(t, content, "")
}

func TestHandleIncrementPlainValue(t *testing.T) {
	handleIncrement("hello", 0, 0, 0, 1)
	content, _ := getCellContent(0, 1)
	assertEqual(t, content, "hello")
}

func TestStrlenWithNumericCell(t *testing.T) {
	// A0 is numeric, so collectParameters stores it as float64.
	// strlen should return an error rather than panic.
	setCellContent(0, 0, "42")
	assertEqual(t, eval("strlen(A0)"), "Error applying strlen: argument must be a string")
}

func TestSumMalformedRange(t *testing.T) {
	// sum with a single argument that has no colon should return an error, not panic.
	assertEqual(t, eval("sum(\"A0\")"), "Error applying sum: argument must be a range in the form A0:B0")
}

func TestSumNonNumericArgs(t *testing.T) {
	// sum(a, b) where arguments are not float64 should return an error.
	setCellContent(0, 0, "hello")
	setCellContent(0, 1, "world")
	assertEqual(t, eval("sum(A0,B0)"), "Error applying sum: arguments must be numeric")
}

func TestRange2D(t *testing.T) {
	// A 2D range (different row and column) should return an error.
	_, err := getCellRange(0, 0, 1, 1)
	if err == nil {
		t.Errorf("expected error for 2D range, got nil")
	}
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
