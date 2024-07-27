package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const PreDownloadedVoiceID = "illust"

type Voice struct {
	ID      string  `json:"id"`
	Path    string  `json:"src"`
	Volume  float64 `json:"volume"`
	First   string  `json:"a"`
	Kana    string  `json:"k"`
	Label   string  `json:"label"`
	VideoID string  `json:"videoId"`
	Time    string  `json:"time"`
}

//go:embed data/ui-shig.jsonl
var data []byte

func ReadVoices() ([]Voice, error) {
	r := bytes.NewReader(data)
	scanner := bufio.NewScanner(r)
	var voices []Voice
	for scanner.Scan() {
		var voice Voice
		err := json.Unmarshal(scanner.Bytes(), &voice)
		if err != nil {
			return nil, fmt.Errorf("おかしいなぁ？%d 件目のデータがないですぅ。 %w", len(voices)+1, err)
		}
		voices = append(voices, voice)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return voices, nil
}

//go:embed data/illust.mp3
var illust []byte

type VoiceURL struct {
	ID   string
	URL  string
	File string
}

func (vu VoiceURL) Load() ([]byte, bool, error) {
	if vu.ID == PreDownloadedVoiceID {
		return illust, false, nil
	}
	contents, err := os.ReadFile(vu.File)
	if err != nil && os.IsExist(err) {
		return nil, false, fmt.Errorf("しぐれういボタンを読み込めませんでした。 [%s] %w", vu.ID, err)
	}
	if contents != nil {
		return contents, false, nil
	}
	var httpClient http.Client
	response, err := httpClient.Get(vu.URL)
	if err != nil {
		return nil, false, fmt.Errorf("しぐれういボタンへのアクセスに失敗しました。 [%s] %w", vu.ID, err)
	}
	responseBody := response.Body
	defer func() { _ = responseBody.Close() }()
	if response.StatusCode < 200 && 300 <= response.StatusCode {
		return nil, false, fmt.Errorf("しぐれういボタンへのアクセスに失敗しました。 [%s] %s", vu.ID, response.Status)
	}
	allBytes, err := io.ReadAll(responseBody)
	if err != nil {
		return nil, false, fmt.Errorf("しぐれういボタンのデータが読み取れませんでした。 [%s] %w", vu.ID, err)
	}
	parentDirectory := vu.GetCacheDir()
	stat, err := os.Stat(parentDirectory)
	if err != nil && os.IsNotExist(err) {
		if err := os.MkdirAll(parentDirectory, 0755); err != nil {
			return nil, false, fmt.Errorf("しぐれういボタンの保存先ディレクトリーを作成できませんでした。 [%s] %w", parentDirectory, err)
		}
	} else if err != nil {
		return nil, false, fmt.Errorf("しぐれういボタンの保存先が確保できませんでした。 [%s] %w", vu.ID, err)
	} else if !stat.IsDir() {
		return nil, false, &UiShigError{
			Message:           "しぐれういボタンの保存先ディレクトリーを確保できませんでした。",
			RecommendedAction: fmt.Sprintf("%s というファイルを削除するか、別の場所に退避してください", parentDirectory),
		}
	}
	if err = os.WriteFile(vu.File, allBytes, 0644); err != nil {
		return nil, false, fmt.Errorf("しぐれういボタンを保存できませんでした。 [%s] %w", vu.ID, err)
	}
	return allBytes, true, nil
}

func (vu VoiceURL) GetCacheDir() string {
	parentDirectory, _ := filepath.Split(vu.File)
	if strings.HasSuffix(parentDirectory, string(os.PathSeparator)) {
		parentDirectory = parentDirectory[:len(parentDirectory)-1]
	}
	return parentDirectory
}

func (vu VoiceURL) Delete() error {
	if vu.ID == PreDownloadedVoiceID {
		return nil
	}
	return os.Remove(vu.File)
}
