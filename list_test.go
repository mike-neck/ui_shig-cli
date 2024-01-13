package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestYouTubeTimeTextToSeconds(t *testing.T) {
	tests := []struct {
		Input  string
		Expect int
	}{
		{
			Input:  "30s",
			Expect: 30,
		},
		{
			Input:  "2m12s",
			Expect: 132,
		},
		{
			Input:  "1h1s",
			Expect: 60*60 + 1,
		},
		{
			Input:  "1h17m",
			Expect: 60*60 + 17*60,
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("Test %d: input[%s] -> %ds", i, test.Input, test.Expect), func(t *testing.T) {
			result := YouTubeTimeTextToSeconds(test.Input)
			assert.Equal(t, test.Expect, result)
		})
	}
}
