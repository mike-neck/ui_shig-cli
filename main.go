package main

import (
	"log"
	"os"
	"path"
)

func main() {
	say := Say{ID: "illust"}
	voices, err := ReadVoices()
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	err = say.Execute(uiShigConfig, voices)
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
