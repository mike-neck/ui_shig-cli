package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
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

func TestVoiceURL_Load(t *testing.T) {
	voices, err := ReadVoices()
	if err != nil {
		t.Errorf("failed to read voices: %v\n", err)
		return
	}
	if len(voices) == 0 {
		t.Error("no voices found")
		return
	}
	src := rand.NewSource(time.Now().UnixMicro())
	r := rand.New(src)
	index := r.Intn(len(voices))
	lastVoice := voices[index]
	config := UiShigConfig{
		UiShigCacheDir: t.TempDir(),
		UiShigURL:      DefaultUiShigURL,
		UiShigReferer:  DefaultUiShigReferer,
	}
	url := config.ResolvePath(lastVoice)
	voice, b, err := url.Load(config)
	if err != nil {
		t.Fatalf("failed to load voice[%v]: %v\n", lastVoice, err)
		return
	}
	if !b {
		t.Fail()
		t.Logf("invalid return value for VoiceURL#Load.succeeded[%v], voice: %v\n", b, lastVoice)
		return
	}
	if len(voice) == 0 {
		t.Fail()
		t.Logf("invalid return value for VoiceURL#Load.voice, voice: %v\n", lastVoice)
		return
	}
	t.Logf("loaded voice[%d, %s, %s]: %d\n", index, lastVoice.ID, url.URL, len(voice))
}
