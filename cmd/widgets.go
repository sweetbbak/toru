package main

import (
	"github.com/charmbracelet/huh"
)

func userInput(title string) string {
	var val string
	huh.NewInput().Title(title).Value(&val).Prompt("> ").Run()
	return val
}

func Confirm() bool {
	var confirm bool
	huh.NewConfirm().
		Title("Are you sure?").
		Affirmative("Yes!").
		Negative("No.").
		Value(&confirm)
	return confirm
}
