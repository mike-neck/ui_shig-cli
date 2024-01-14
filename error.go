package main

import (
	"fmt"
	"strings"
)

type UiShigError struct {
	Message           string
	RecommendedAction string
}

func (v *UiShigError) Error() string {
	action := fmt.Sprintf("対処方法: %s", v.RecommendedAction)
	return strings.Join([]string{"[ERROR]", v.Message, action}, "\n")
}
