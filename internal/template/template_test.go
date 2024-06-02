package template

import (
	"regexp"
	"testing"
)

type GetlastIdMatchingTest struct {
	names []string
	regex *regexp.Regexp
	want  int
}

func TestGetLastIdMatching(t *testing.T) {
	tests := []GetlastIdMatchingTest{
		{names: make([]string, 0), regex: regexp.MustCompile(`(\d)`), want: 0},
		{names: []string{"banana", "ruby"}, regex: regexp.MustCompile(`(\d)`), want: 0},
		{names: []string{"banana 1", "ruby"}, regex: regexp.MustCompile(`(\d)`), want: 1},
		{names: []string{"banana 1", "ruby 2"}, regex: regexp.MustCompile(`(\d)`), want: 2},
		{names: []string{"banana 3", "ruby 2"}, regex: regexp.MustCompile(`(\d)`), want: 3},
		{names: []string{"banana 3", "ruby 2", "goat 4"}, regex: regexp.MustCompile(`(\d)`), want: 4},
	}

	for _, tc := range tests {
		got, err := getLastIdMatching(tc.names, tc.regex)
		if got != tc.want || err != nil {
			t.Errorf(`got %v, want %v, err %v`, got, tc.want, err)
		}
	}

}
