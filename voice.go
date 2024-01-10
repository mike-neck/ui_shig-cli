package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
)

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
