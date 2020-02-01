package test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	lasgo "github.com/iykekings/las-go"
)

type dataRow struct {
	Dept float64 `las:"DEPT"`
	Dt   string  `las:"DT"`
	Rhob float64 `las:"RHOB"`
	Nphi string  `las:"NPHI"`
	Sflu float64 `las:"SFLU"`
	Sfla float64 `las:"SFLA"`
	Ilm  string  `las:"ILM"`
	Ild  float64 `las:"ILD"`
}

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

func TestDataStruct(t *testing.T) {
	las, err := lasgo.Las("../sample/example1.las")
	if err != nil {
		panic(err)
	}

	expected := []interface{}{
		&dataRow{
			Dept: float64(1670),
			Dt:   string("123.450"),
			Rhob: float64(2550),
			Nphi: string("0.450"),
			Sflu: float64(123.45),
			Sfla: float64(123.45),
			Ilm:  string("110.200"),
			Ild:  float64(105.6),
		},
		&dataRow{
			Dept: float64(1669.875),
			Dt:   string("123.450"),
			Rhob: float64(2550),
			Nphi: string("0.450"),
			Sflu: float64(123.45),
			Sfla: float64(123.45),
			Ilm:  string("110.200"),
			Ild:  float64(105.6),
		},
		&dataRow{
			Dept: float64(1669.75),
			Dt:   string("123.450"),
			Rhob: float64(2550),
			Nphi: string("0.450"),
			Sflu: float64(123.45),
			Sfla: float64(123.45),
			Ilm:  string("110.200"),
			Ild:  float64(105.6),
		},
		&dataRow{
			Dept: float64(1669.745),
			Dt:   string("123.450"),
			Rhob: float64(2550),
			Nphi: string("-999.25"),
			Sflu: float64(123.45),
			Sfla: float64(123.45),
			Ilm:  string("110.200"),
			Ild:  float64(105.6),
		},
	}

	opts := &lasgo.DataOptions{ConcreteStruct: dataRow{}}

	actual := las.DataStruct(opts)

	if !cmp.Equal(expected, actual) {
		t.Errorf("wrong val: expected: %T %v actual: %T %v\n", expected, expected, actual, actual)
	}

}
