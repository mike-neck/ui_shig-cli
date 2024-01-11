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
	"path"
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

type VoiceError struct {
	Message           string
	RecommendedAction string
}

func (v *VoiceError) Error() string {
	return strings.Join([]string{"[ERROR]", v.Message, v.RecommendedAction}, "\n")
}

//go:embed data/illust.mp3
var illust []byte

type VoiceURL struct {
	ID   string
	URL  string
	File string
}

func (vu VoiceURL) Load() ([]byte, error) {
	if vu.ID == PreDownloadedVoiceID {
		return illust, nil
	}
	contents, err := os.ReadFile(vu.File)
	if err != nil && os.IsExist(err) {
		return nil, fmt.Errorf("しぐれういボタンを読み込めませんでした。 [%s] %w", vu.ID, err)
	}
	if contents != nil {
		return contents, nil
	}
	var httpClient http.Client
	response, err := httpClient.Get(vu.URL)
	if err != nil {
		return nil, fmt.Errorf("しぐれういボタンへのアクセスに失敗しました。 [%s] %w", vu.ID, err)
	}
	responseBody := response.Body
	defer func() { _ = responseBody.Close() }()
	if 200 <= response.StatusCode && response.StatusCode < 300 {
		return nil, fmt.Errorf("しぐれういボタンへのアクセスに失敗しました。 [%s] %s", vu.ID, response.Status)
	}
	allBytes, err := io.ReadAll(responseBody)
	if err != nil {
		return nil, fmt.Errorf("しぐれういボタンのデータが読み取れませんでした。 [%s] %w", vu.ID, err)
	}
	parentDirectory, _ := path.Split(vu.File)
	stat, err := os.Stat(parentDirectory)
	if err != nil && os.IsNotExist(err) {
		if err := os.MkdirAll(parentDirectory, 0755); err != nil {
			return nil, fmt.Errorf("しぐれういボタンの保存先ディレクトリーを作成できませんでした。 [%s] %w", parentDirectory, err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("しぐれういボタンの保存先が確保できませんでした。 [%s] %w", vu.ID, err)
	} else if !stat.IsDir() {
		return nil, &VoiceError{
			Message:           "しぐれういボタンの保存先ディレクトリーを確保できませんでした。",
			RecommendedAction: fmt.Sprintf("%s というファイルを削除するか、別の場所に退避してください", parentDirectory),
		}
	}
	if err = os.WriteFile(vu.File, allBytes, 0644); err != nil {
		return nil, fmt.Errorf("しぐれういボタンを保存できませんでした。 [%s] %w", vu.ID, err)
	}
	return allBytes, nil
}
