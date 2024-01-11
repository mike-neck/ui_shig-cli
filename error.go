package main

import "strings"

type UiShigError struct {
	Message           string
	RecommendedAction string
}

func (v *UiShigError) Error() string {
	return strings.Join([]string{"[ERROR]", v.Message, v.RecommendedAction}, "\n")
}
