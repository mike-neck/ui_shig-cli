package main

import (
	"path"
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
	file := path.Join(dir, voicePath)

	return VoiceURL{
		ID:   voice.ID,
		URL:  url,
		File: file,
	}
}
