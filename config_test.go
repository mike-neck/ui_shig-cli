package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func TestUiShigConfig_ResolvePath(t *testing.T) {
	type expectation struct {
		cacheDir     string
		expectedPath string
		expectedURL  string
	}
	exp := expectation{
		expectedURL: "https://example.com/sound/test.mp3",
	}
	if runtime.GOOS == "windows" {
		exp.cacheDir = "C:\\user\\test\\.ui_shig\\cache"
		exp.expectedPath = fmt.Sprintf("%s\\sound\\test.mp3", exp.cacheDir)
	} else {
		exp.cacheDir = "/home/test/.ui_shig/caches"
		exp.expectedPath = fmt.Sprintf("%s/sound/test.mp3", exp.cacheDir)
	}
	type testData struct {
		Voice
		UiShigConfig
	}
	noSuffixURL := "https://example.com"
	suffixedURL := "https://example.com/"
	tests := []testData{
		{
			Voice: Voice{
				Path: "./sound/test.mp3",
			},
			UiShigConfig: UiShigConfig{
				UiShigURL:      noSuffixURL,
				UiShigCacheDir: exp.cacheDir,
				IssueURL:       "https://example.com/issues",
			},
		},
		{
			Voice: Voice{
				Path: "/sound/test.mp3",
			},
			UiShigConfig: UiShigConfig{
				UiShigURL:      noSuffixURL,
				UiShigCacheDir: exp.cacheDir,
				IssueURL:       "https://example.com/issues",
			},
		},
		{
			Voice: Voice{
				Path: "sound/test.mp3",
			},
			UiShigConfig: UiShigConfig{
				UiShigURL:      noSuffixURL,
				UiShigCacheDir: exp.cacheDir,
				IssueURL:       "https://example.com/issues",
			},
		},
		{
			Voice: Voice{
				Path: "./sound/test.mp3",
			},
			UiShigConfig: UiShigConfig{
				UiShigURL:      suffixedURL,
				UiShigCacheDir: exp.cacheDir,
				IssueURL:       "https://example.com/issues",
			},
		},
		{
			Voice: Voice{
				Path: "/sound/test.mp3",
			},
			UiShigConfig: UiShigConfig{
				UiShigURL:      suffixedURL,
				UiShigCacheDir: exp.cacheDir,
				IssueURL:       "https://example.com/issues",
			},
		},
		{
			Voice: Voice{
				Path: "sound/test.mp3",
			},
			UiShigConfig: UiShigConfig{
				UiShigURL:      suffixedURL,
				UiShigCacheDir: exp.cacheDir,
				IssueURL:       "https://example.com/issues",
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d url[%s] + path[%s]", i, test.UiShigURL, test.Path), func(t *testing.T) {
			config := test.UiShigConfig
			voice := test.Voice
			v := config.ResolvePath(voice)
			assert.Equal(t, exp.expectedPath, v.File)
			assert.Equal(t, exp.expectedURL, v.URL)
		})
	}
}
