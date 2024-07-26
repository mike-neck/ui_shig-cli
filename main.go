package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path"
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
		Description: "しぐれういの音声が聞けるコマンドライン・アプリケーションです。しぐれういボタン(http://cbtm.html.xdomain.jp/usbtn/usbtn.html)をスクレイピングしいているだけです。",
		Commands:    commands,
	}
	cli.AppHelpTemplate = fmt.Sprintf(`%s
------
しぐれうい CLI

  バージョン: %s
  ビルド番号: %s
  ビルド日時: %s

`, cli.AppHelpTemplate, UiShigVersion, UiShigCommit, UiShigBuildDate)
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
		UiShigCacheDir: path.Join(homeDirectory, DefaultDirectoryName, DefaultCacheDirectoryName),
	}
	return &uiShigConfig, nil
}

//goland:noinspection HttpUrlsUsage
const (
	DefaultUiShigURL          = "http://cbtm.html.xdomain.jp/usbtn"
	DefaultDirectoryName      = ".ui_shig"
	DefaultCacheDirectoryName = "caches"
	DefaultIssueURL           = "https://github.com/mike-neck/ui_shig-cli/issues/new"
)
