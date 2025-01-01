package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path/filepath"
)

var (
	UiShigVersion   = "v0.0.0"
	UiShigCommit    = "000000"
	UiShigBuildDate = "2022-02-02T02:02:02Z"
)

func main() {
	uiShigConfig, err := newUiShigConfig()
	if err != nil {
		log.Fatalf("user home directory error: %v\n", err)
	}
	voices, err := ReadVoices()
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	userOrders := []UserOrder{
		SayUserCommand,
		ListUserOrder,
		CacheUserOrder(),
	}
	commands := make([]*cli.Command, len(userOrders))
	for i, order := range userOrders {
		commands[i] = order.toCLICommand(*uiShigConfig, voices)
	}
	app := &cli.App{
		Name:        "しぐれういCLI",
		Description: "しぐれういの音声が聞けるコマンドライン・アプリケーションです。しぐれういボタン(https://leiros.cloudfree.jp/usbtn/usbtn.html)をスクレイピングしいているだけです。",
		Commands:    commands,
	}
	cli.AppHelpTemplate = fmt.Sprintf(`%s
------
しぐれうい CLI

  バージョン: %s
  ビルド番号: %s
  ビルド日時: %s

------
お困りの際はこちらから質問・不具合報告をしてください。

  GitHub: %s

`, cli.AppHelpTemplate, UiShigVersion, UiShigCommit, UiShigBuildDate, DefaultIssueURL)
	err = app.Run(os.Args)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
}

func newUiShigConfig() (*UiShigConfig, error) {
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	var uiShigConfig = UiShigConfig{
		UiShigURL:      DefaultUiShigURL,
		UiShigCacheDir: filepath.Join(homeDirectory, DefaultDirectoryName, DefaultCacheDirectoryName),
	}
	return &uiShigConfig, nil
}

//goland:noinspection HttpUrlsUsage
const (
	DefaultUiShigURL          = "https://leiros.cloudfree.jp/usbtn/usbtn.html"
	DefaultDirectoryName      = ".ui_shig"
	DefaultCacheDirectoryName = "caches"
	DefaultIssueURL           = "https://github.com/mike-neck/ui_shig-cli/issues/new"
)
