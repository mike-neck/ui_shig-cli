package main

import (
	"fmt"
	"os"
	"strings"
)

var cacheArgumentDescription = "削除/ダウンロードするしぐれういボタンのID"

func CacheUserOrder() UserOrder {
	cache := Cache{
		Delete:    false,
		ShowLogs:  false,
		TargetIds: nil,
	}
	return UserOrder{
		Name:                "cache",
		Description:         "キャッシュの削除/ダウンロードをおこないます",
		ArgumentDescription: &cacheArgumentDescription,
		IntOptions:          []IntOption{},
		StringOptions:       []StringOption{},
		FileOptions:         []FileOption{},
		BoolOptions: []BoolOption{
			{
				Name:        "delete",
				Aliases:     []string{"d"},
				Description: "キャッシュ削除モードで起動します",
				Required:    false,
				Reference:   &cache.Delete,
			},
			{
				Name:        "show-logs",
				Aliases:     []string{"verbose", "v"},
				Description: "ログを表示します。",
				Required:    false,
				Reference:   &cache.ShowLogs,
			},
		},
		ConstructCommand: func(order UserOrder, args []string) (Command, error) {
			cache.TargetIds = args
			return &cache, nil
		},
	}
}

type Cache struct {
	Delete    bool
	ShowLogs  bool
	TargetIds []string
}

func (cache *Cache) Execute(config UiShigConfig, voices []Voice) error {
	urls := make([]VoiceURL, 0)
	if cache.TargetIds == nil || len(cache.TargetIds) == 0 {
		for _, voice := range voices {
			url := uiShigConfig.ResolvePath(voice)
			urls = append(urls, url)
		}
	} else {
		for _, id := range cache.TargetIds {
			for _, voice := range voices {
				if voice.ID != id {
					continue
				}
				url := config.ResolvePath(voice)
				urls = append(urls, url)
				break
			}
		}
	}
	if len(urls) == 0 {
		return &UiShigError{
			Message:           fmt.Sprintf("指定したしぐれういボタンが見つかりませんでした。 [%s]", cache.TargetIds),
			RecommendedAction: "ui_shig list コマンドを使ってしぐれういボタンを探してみてください。",
		}
	}
	es := make([]string, 0)
	count := 0
	if cache.Delete {
		for _, url := range urls {
			err := url.Delete()
			if err != nil {
				if cache.ShowLogs {
					_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf("error: deleting cache[%s]: %s, %v", url.ID, url.File, err))
				}
				es = append(es, url.ID)
			} else {
				if cache.ShowLogs {
					fmt.Printf("deleted cache[%s]: %s\n", url.ID, url.File)
				}
				count++
			}
		}
	} else {
		for _, url := range urls {
			_, downloaded, err := url.Load()
			if err != nil {
				if cache.ShowLogs {
					_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf("error: downloading cache[%s]: %s, %v", url.ID, url.File, err))
				}
				es = append(es, url.ID)
			} else {
				if cache.ShowLogs {
					fmt.Printf("downloaded cache[%s]: %s\n", url.ID, url.File)
				}
				if downloaded {
					count++
				}
			}
		}
	}
	if len(es) != 0 {
		ids := strings.Join(es, ",")
		return &UiShigError{
			Message:           fmt.Sprintf("以下のしぐれういボタンのキャッシュの削除/ダウンロードに失敗しました。 [%s]", ids),
			RecommendedAction: "ui_shig cache コマンドを使って再度試してみてください。",
		}
	}
	if cache.ShowLogs {
		fmt.Printf("Total %d cache files processed\n", count)
	}
	return nil
}
