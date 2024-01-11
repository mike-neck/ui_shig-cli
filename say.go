package main

import (
	"bytes"
	"fmt"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"io"
	"time"
)

type Say struct {
	ID string
}

func (say Say) Execute(config UiShigConfig, voices []Voice) error {
	for _, voice := range voices {
		if voice.ID != say.ID {
			continue
		}
		voiceURL := config.ResolvePath(voice)
		contents, err := voiceURL.Load()
		if err != nil {
			return err
		}
		return say.PlayBytes(config, voiceURL, contents)
	}
	return &UiShigError{
		Message:           fmt.Sprintf("指定したしぐれういボタンが見つかりませんでした。 [%s]", say.ID),
		RecommendedAction: "ui_shig list コマンドを使ってしぐれういボタンを探してみてください。",
	}
}

func (say Say) PlayBytes(config UiShigConfig, v VoiceURL, voiceBytes []byte) error {
	reader := bytes.NewReader(voiceBytes)
	rc := io.NopCloser(reader)
	stream, format, err := mp3.Decode(rc)
	if err != nil {
		return &UiShigError{
			Message:           fmt.Sprintf("指定したしぐれういボタンが壊れていました。 [%s] %w", say.ID, err),
			RecommendedAction: fmt.Sprintf("ダウンロードしたファイル(%s)を一度ゴミ箱に捨てて再度試してください。それでも駄目な場合はWEBのしぐれういボタンが変更された可能性があるので、こちら(%s)に連絡してみてください。", v.File, config.IssueURL),
		}
	}
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Millisecond*100))
	if err != nil {
		_ = stream.Close()
		//TODO UiShigError を使う
		return err
	}
	buffer := beep.NewBuffer(format)
	buffer.Append(stream)
	_ = stream.Close()
	newStream := buffer.Streamer(0, buffer.Len())
	ch := make(chan interface{})
	speaker.Play(beep.Seq(newStream, beep.Callback(func() {
		ch <- "finished"
	})))
	<-ch
	return nil
}
