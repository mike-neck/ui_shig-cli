package main

import (
	"fmt"
	"strconv"
	"strings"
)

var ListUserOrder = UserOrder{
	Name:                "list",
	Description:         "しぐれういの音声IDなどの情報を確認します。",
	ArgumentDescription: &sayUserCommandArgumentsDescription,
	IntOptions:          []IntOption{},
	StringOptions:       []StringOption{},
	FileOptions:         []FileOption{},
	ConstructCommand: func(order UserOrder, args []string) (Command, error) {
		if len(args) == 0 {
			return List{ListItemID}, nil
		}
		all := List{
			ListItemID,
			ListItemSrc,
			ListItemKana,
			ListItemLabel,
			ListItemYouTube,
		}
		list := make(List, 0)
		for _, arg := range args {
			listItem := ListItem(arg)
			index := all.Index(listItem)
			if index != -1 {
				list = append(list, listItem)
			}
		}
		if len(list) == 0 {
			return nil, listItemNotFoundError()
		}
		return list, nil
	},
}

func listItemNotFoundError() error {
	candidates := make([]string, 0)
	for _, itemMap := range ListItemMaps {
		candidates = append(candidates, string(itemMap.ListItem))
	}
	candidateItems := strings.Join(candidates, ", ")
	return &UiShigError{
		Message:           "出力できる項目がありませんでした。",
		RecommendedAction: fmt.Sprintf("[%s] の中から選んでください。", candidateItems),
	}
}

type ListItem string

const (
	ListItemID      ListItem = "id"
	ListItemSrc     ListItem = "src"
	ListItemKana    ListItem = "kana"
	ListItemLabel   ListItem = "label"
	ListItemYouTube ListItem = "youtube"
)

type ListItemFieldMap struct {
	ListItem
	GetValue func(voice Voice) string
}

var ListItemMaps = []ListItemFieldMap{
	{
		ListItem: ListItemID,
		GetValue: func(voice Voice) string {
			return voice.ID
		},
	},
	{
		ListItem: ListItemSrc,
		GetValue: func(voice Voice) string {
			return voice.Path
		},
	},
	{
		ListItem: ListItemKana,
		GetValue: func(voice Voice) string {
			return voice.Kana
		},
	},
	{
		ListItem: ListItemLabel,
		GetValue: func(voice Voice) string {
			return voice.Label
		},
	},
	{
		ListItem: ListItemYouTube,
		GetValue: func(voice Voice) string {
			if voice.Time == "" {
				return ""
			}
			if voice.VideoID == "" {
				return ""
			}
			sec := YouTubeTimeTextToSeconds(voice.Time)
			return fmt.Sprintf("%s/%s?t=%d", YouTubeBaseURL, voice.VideoID, sec)
		},
	},
}

const YouTubeBaseURL = "https://youtu.be"

type List []ListItem

func (list List) Uniq() List {
	newList := make(List, 0)
	m := make(map[ListItem]bool, 0)
	for _, item := range list {
		if _, hasValue := m[item]; !hasValue {
			m[item] = true
			newList = append(newList, item)
		}
	}
	return newList
}

func (list List) Index(item ListItem) int {
	for i, listItem := range list {
		if listItem == item {
			return i
		}
	}
	return -1
}

func (list List) Execute(config UiShigConfig, voices []Voice) error {
	uniq := list.Uniq()
	lines := make([][]string, 0)
	for _, voice := range voices {
		line := make([]string, 0)
		for _, mapping := range ListItemMaps {
			index := uniq.Index(mapping.ListItem)
			if index == -1 {
				continue
			}
			value := mapping.GetValue(voice)
			if value != "" {
				line = append(line, value)
			}
		}
		if len(uniq) == len(line) {
			lines = append(lines, line)
		}
	}
	if len(lines) == 0 {
		return listItemNotFoundError()
	}
	widths := make([]int, len(uniq))
	for i := range uniq {
		widths[i] = 0
		for _, line := range lines {
			item := line[i]
			width := getCharacterWidth(item)
			if width > widths[i] {
				widths[i] = width
			}
		}
	}
	for _, line := range lines {
		for i := range line {
			item := line[i]
			width := widths[i]
			nonASCII := getNonASCIICharacterCount(item)
			fmt.Printf("%-*s ", width-nonASCII, item)
		}
		fmt.Println()
	}
	return nil
}

func YouTubeTimeTextToSeconds(text string) int {
	start := 0
	//2h8m3s 形式の時間を秒(変数=start)に変換する
	t := text
	for strings.Contains(t, "h") || strings.Contains(t, "m") || strings.Contains(t, "s") {
		if strings.Contains(t, "h") {
			index := strings.Index(t, "h")
			hours := t[:index]
			t = t[index+1:]
			h, _ := strconv.Atoi(hours)
			start += h * 3600
		}
		if strings.Contains(t, "m") {
			index := strings.Index(t, "m")
			minutes := t[:index]
			t = t[index+1:]
			m, _ := strconv.Atoi(minutes)
			start += m * 60
		}
		if strings.Contains(t, "s") {
			index := strings.Index(t, "s")
			seconds := t[:index]
			t = t[index+1:]
			s, _ := strconv.Atoi(seconds)
			start += s
		}
	}
	return start
}

func getCharacterWidth(s string) int {
	width := 0
	for _, r := range s {
		if r <= 127 { // ASCII
			width += 1
		} else {
			width += 2
		}
	}
	return width
}

func getNonASCIICharacterCount(s string) int {
	count := 0
	for _, r := range s {
		if r > 127 { // ASCII
			count += 1
		}
	}
	return count
}
