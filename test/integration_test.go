package test

import (
	"fmt"
	"testing"

	lasgo "github.com/iykekings/las-go"
)

func TestRowCount(t *testing.T) {
	las, err := lasgo.Las("../sample/example1.las")
	if err != nil {
		panic(err)
	}
	row := las.RowCount()
	fmt.Println(row)
	if row != 4 {
		t.Errorf("las.RowCount() == %q, want %q", row, 4)
	}
}
func TestCoulmnCount(t *testing.T) {
	las, err := lasgo.Las("../sample/example1.las")
	if err != nil {
		panic(err)
	}
	column := las.ColumnCount()
	if column != 8 {
		t.Errorf("las.ColumnCount() == %q, want %q", column, 4)
	}
}
func TestVersion(t *testing.T) {
	las, err := lasgo.Las("../sample/example1.las")
	if err != nil {
		panic(err)
	}
	vers := las.Version()
	if vers != "2.0" {
		t.Errorf("las.ColumnCount() == %q, want %q", vers, "2.0")
	}
}
func TestWrap(t *testing.T) {
	las, err := lasgo.Las("../sample/example1.las")
	if err != nil {
		panic(err)
	}
	wrap := las.Wrap()
	if wrap != false {
		t.Errorf("las.ColumnCount() == %v, want %v", wrap, false)
	}
}
