package lasgo

import (
	"reflect"
	"testing"
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
