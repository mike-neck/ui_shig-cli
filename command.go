package main

type Command interface {
	Execute(config UiShigConfig, voices []Voice) error
}
