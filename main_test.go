package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsCommitMessageMinor(t *testing.T) {
	var tests = []struct {
		input string
		want  bool
	}{
		{"", false},
		{`Some message
		MINOR release
		
		footer`, true},
		{"some message", false},
		{"Minor release", false},
	}
	for _, test := range tests {
		assert.Equal(t, isCommitMessageMinor(test.input), test.want, "Must be equal")
	}
}

func TestIsCommitMessageMajor(t *testing.T) {
	var tests = []struct {
		input string
		want  bool
	}{
		{"", false},
		{`Some message
		MAJOR release
		
		footer`, true},
		{"some message", false},
		{"Minor release", false},
	}
	for _, test := range tests {
		assert.Equal(t, isCommitMessageMajor(test.input), test.want, "Must be equal")
	}
}
