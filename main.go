package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path"
)

func main() {
	voices, err := ReadVoices()
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	userOrders := []UserOrder{
		SayUserCommand,
		ListUserOrder,
	}
	commands := make([]*cli.Command, len(userOrders))
	for i, order := range userOrders {
		commands[i] = order.toCLICommand(uiShigConfig, voices)
	}
	app := &cli.App{
		Name:        "しぐれういCLI",
		Description: "しぐれういの音声が聞けるコマンドライン・アプリケーションです。しぐれういボタンをスクレイピングしいているだけです。",
		Commands:    commands,
	}
	err = app.Run(os.Args)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
}

var uiShigConfig = UiShigConfig{
	UiShigURL:      DefaultUiShigURL,
	UiShigCacheDir: path.Join(os.Getenv("HOME"), DefaultDirectoryName, DefaultCacheDirectoryName),
}

//goland:noinspection HttpUrlsUsage
const (
	DefaultUiShigURL          = "http://cbtm.html.xdomain.jp/usbtn"
	DefaultDirectoryName      = ".ui_shig"
	DefaultCacheDirectoryName = "caches"
	DefaultIssueURL           = "https://github.com/mike-neck/ui_shig-cli/issues/new"
)
