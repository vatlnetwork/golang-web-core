package util

import (
	"fmt"
	"log"
)

func PrintColor(color, format string, parts ...any) {
	strng := fmt.Sprintf(format, parts...)

	_, ok := colors[color]
	if ok {
		color = colors[color]
	}

	fmt.Printf("\033[38;2;%vm%v\033[0m", color, strng)
}

func LogColor(color, format string, parts ...any) {
	strng := fmt.Sprintf(format, parts...)

	_, ok := colors[color]
	if ok {
		color = colors[color]
	}

	log.Printf("\033[38;2;%vm%v\033[0m", color, strng)
}

func LogFatal(err error) {
	log.Fatalf("\033[38;2;%vm%v\033[0m", colors["red"], err)
}

var colors map[string]string = map[string]string{
	"green":      "0;255;0",
	"lightgreen": "100;255;100",
	"red":        "255;0;0",
	"blue":       "0;0;255",
}
