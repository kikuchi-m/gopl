package main

import (
	"reflect"
	"testing"
)

func TestRemoveAdjacent(t *testing.T) {
	var tests = []struct {
		input    []string
		expected []string
	}{
		{[]string{}, []string{}},
		{[]string{"a"}, []string{"a"}},
		{[]string{"a", "a"}, []string{"a"}},
		{[]string{"a", "b", "a"}, []string{"a", "b", "a"}},
		{[]string{"a", "b", "a", "a"}, []string{"a", "b", "a"}},
		{[]string{"a", "b", "b", "b"}, []string{"a", "b"}},
	}

	for _, test := range tests {
		var res = RemoveAdjacent(test.input)
		if !reflect.DeepEqual(test.expected, res) {
			t.Errorf("result is expected %v, but %v", test.expected, res)
		}
		if !reflect.DeepEqual(test.expected, test.input[:len(test.expected)]) {
			t.Errorf("intput is expected to updated %v and rest, but %v", test.expected, test.input)
		}
	}
}
