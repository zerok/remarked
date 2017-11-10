package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseLineRanges(t *testing.T) {
	result := parseLineRanges("2-5,8")
	expected := map[int]struct{}{
		2: struct{}{},
		3: struct{}{},
		4: struct{}{},
		5: struct{}{},
		8: struct{}{},
	}
	require.Equal(t, expected, result)
}
