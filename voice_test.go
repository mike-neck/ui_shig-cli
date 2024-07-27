package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestVoiceURL_GetCacheDir(t *testing.T) {
	t.Run("DirectoryName Matches ExpectedName", func(t *testing.T) {
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
	})

	t.Run("DirectoryName Is Valid", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "voice-url-test*user")
		if err != nil {
			t.Fatal("Failed to create temp dir", err)
			return
		}
		//goland:noinspection GoUnhandledErrorResult
		defer os.RemoveAll(tempDir)

		path := filepath.Join(tempDir, ".ui_shig", "cache", "sound", "test.mp3")
		voiceURL := VoiceURL{File: path}

		err = voiceURL.CreateCacheDirIfNotExists()

		if !assert.Nil(t, err) {
			return
		}

		cacheDir := voiceURL.GetCacheDir()
		if !assert.DirExists(t, cacheDir) {
			return
		}

		f, err := os.Create(path)
		if err != nil {
			t.Fatal("Failed to open file", err)
			return
		}
		//goland:noinspection GoUnhandledErrorResult
		defer f.Close()
		_, err = f.Write([]byte("test"))
		assert.Nilf(t, err, "failed to write to file: %v\n", err)
	})
}
