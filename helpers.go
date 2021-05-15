package main

import "strings"

func Parameterize(s string) (newString string) {
	newString = strings.ReplaceAll(s, " ", "_")
	newString = strings.ReplaceAll(newString, "/", "_OR_")

	return
}
