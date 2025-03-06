package stats

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAll(t *testing.T) {
	pathToTestData := "./testdata/"
	result, err := TopWords(pathToTestData, 4)
	if err != nil {
		t.Fatalf("failed to get stats: %v", err)
	}
	expected := `"five": 5
"four": 4
"three": 3
"two": 2`
	diff := cmp.Diff(expected, result)
	if diff != "" {
		t.Errorf("want -, got +\n%s", diff)
	}
}
