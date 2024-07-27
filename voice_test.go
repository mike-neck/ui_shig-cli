package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func TestVoiceURL_GetCacheDir(t *testing.T) {
	type testData struct {
		isWindows        bool
		expectedCacheDir string
		VoiceURL
	}
	data := []testData{
		{
			isWindows:        true,
			expectedCacheDir: "C:\\Users\\test\\.ui_shig\\cache\\sound",
			VoiceURL: VoiceURL{
				File: "C:\\Users\\test\\.ui_shig\\cache\\sound\\test.mp3",
			},
		},
		{
			isWindows:        false,
			expectedCacheDir: "/Users/test/.ui_shig/cache/sound",
			VoiceURL: VoiceURL{
				File: "/Users/test/.ui_shig/cache/sound/test.mp3",
			},
		},
	}
	for _, d := range data {
		if (runtime.GOOS == "windows") && d.isWindows {
			t.Run(fmt.Sprintf("GetCacheDir-windows[file=%s]", d.File), func(t *testing.T) {
				dir := d.GetCacheDir()
				assert.Equal(t, d.expectedCacheDir, dir)
			})
		}
	}
}
