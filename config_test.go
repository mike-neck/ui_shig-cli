package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func TestUiShigConfig_ResolvePath(t *testing.T) {
	voice := Voice{
		Path: "./sound/test.mp3",
	}
	type testData struct {
		cacheDir     string
		expectedPath string
		expectedURL  string
	}
	td := testData{
		expectedURL: "https://example.com/sound/test.mp3",
	}
	if runtime.GOOS == "windows" {
		td.cacheDir = "C:\\user\\test\\.ui_shig\\cache"
		td.expectedPath = fmt.Sprintf("%s\\sound\\test.mp3", td.cacheDir)
	} else {
		td.cacheDir = "/home/test/.ui_shig/caches"
		td.expectedPath = fmt.Sprintf("%s/sound/test.mp3", td.cacheDir)
	}

	config := UiShigConfig{
		UiShigURL:      "https://example.com/",
		UiShigCacheDir: td.cacheDir,
		IssueURL:       "https://example.com/issues",
	}

	v := config.ResolvePath(voice)
	assert.Equal(t, td.expectedPath, v.File)
	assert.Equal(t, td.expectedURL, v.URL)
}
