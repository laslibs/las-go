package lasgo

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-cmp/cmp"
	"github.com/mitchellh/mapstructure"
)

func TestChunk(t *testing.T) {
	testData := []string{"1", "2", "3", "4", "5", "6", "7"}
	wanted := [][]string{{"1", "2"}, {"3", "4"}, {"5", "6"}, {"7"}}
	got := chunk(testData, 2)
	if !reflect.DeepEqual(got, wanted) {
		t.Errorf("chunk(%q) == %q, want %q", testData, got, wanted)
	}
}

func TestRemoveComment(t *testing.T) {
	testData := []string{"#ignore\nDon't Ignore\n#\n  # also ignore", "Don't ignore\n#\nBe there"}
	wanted := [][]string{{"Don't Ignore"}, {"Don't ignore", "Be there"}}
	for i := 0; i < len(testData); i++ {
		got := removeComment(testData[i])
		if !reflect.DeepEqual(got, wanted[i]) {
			t.Errorf("chunk(%q) == %q, want %q", testData[i], got, wanted[i])
		}
	}
}

func TestStructConvert(t *testing.T) {
	type list struct {
		Index int    `las:"index"`
		Item  string `las:"item"`
	}

	testData := []string{"1", "books", "2", "bicycles", "3", "cars", "4", "computers"}

	header := []string{"index", "item"}

	chnk := chunk(testData, len(header))

	expected := []interface{}{
		&list{Index: int(1), Item: "books"},
		&list{Index: int(2), Item: "bicycles"},
		&list{Index: int(3), Item: "cars"},
		&list{Index: int(4), Item: "computers"},
	}

	opts := &DataOptions{
		ConcreteStruct: list{},
		DecoderConfig: &StructorConfig{
			DecodeHook:       mapstructure.StringToTimeHookFunc(time.RFC3339),
			WeaklyTypedInput: true,
		},
	}

	// opts := &DataOptions{ConcreteStruct: list{}}

	ctx := context.Background()
	actual, err := structConvert(ctx, &chnk, header, opts)
	if err != nil {
		t.Errorf("Error encountered: %s\n", err)
	}
	spew.Dump(actual)

	if !cmp.Equal(expected, actual) {
		t.Errorf("wrong val: expected: %T %v actual: %T %v\n", expected, expected, actual, actual)
	}

}
