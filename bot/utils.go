package main

import (
	"slices"
	"strings"
)

func GetRawCommand(msgText string) []string {
	rawCommand := strings.Split(msgText, " ")
	return slices.DeleteFunc(rawCommand, func(s string) bool {
		return s == ""
	})
}
