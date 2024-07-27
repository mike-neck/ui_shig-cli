package main

import (
	"os"
	"path/filepath"
	"strings"
)

type UiShigConfig struct {
	UiShigURL      string
	UiShigCacheDir string
	IssueURL       string
}

func (uc UiShigConfig) ResolvePath(voice Voice) VoiceURL {
	baseURL := uc.UiShigURL
	if strings.HasSuffix(baseURL, "/") {
		index := strings.LastIndex(baseURL, "/")
		baseURL = baseURL[:index]
	}
	voicePath := voice.Path
	if strings.HasPrefix(voicePath, "./") || strings.HasPrefix(voicePath, "/") {
		index := strings.Index(voicePath, "/")
		voicePath = voicePath[index+1:]
	}
	fragments := []string{
		baseURL,
		voicePath,
	}
	url := strings.Join(fragments, "/")

	dir := uc.UiShigCacheDir
	voidFilePath := strings.ReplaceAll(voicePath, "/", string(os.PathSeparator))
	file := filepath.Join(dir, voidFilePath)

	return VoiceURL{
		ID:   voice.ID,
		URL:  url,
		File: file,
	}
}
