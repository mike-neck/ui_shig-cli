package main

import "testing"

func TestNewUiShigConfig(t *testing.T) {
	config, err := newUiShigConfig()
	if err != nil {
		t.Fatalf("newUiShigConfig() failed: %v", err)
		return
	}
	if config == nil {
		t.Fatalf("newUiShigConfig() returned nil config")
		return
	}
}
