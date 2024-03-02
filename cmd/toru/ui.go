package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/sweetbbak/toru/pkg/nyaa"
	"github.com/sweetbbak/toru/pkg/search"
)

var cutePrint = lipgloss.NewStyle().Width(40).BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("63")).Padding(1, 2)

// white text, purple background
var style = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4"))
	// PaddingTop(2).
	// PaddingLeft(4).
	// Width(22)

// open old search from json cache
func FromCache(jsonFile string) ([]nyaa.Media, error) {
	fmt.Println(cutePrint.Align(lipgloss.Center).Render("Toru anime streaming"))
	res, err := search.ReadCache(jsonFile)
	if err != nil {
		log.Fatal(err)
		return res, err
	}

	return res, nil
}

func Prompt(prompt string) (string, error) {
	var val string
	ui := huh.NewInput().Value(&val).Description(prompt).Prompt("> ")
	err := ui.Run()
	if err != nil {
		return "", err
	}

	fmt.Println(style.Render(val))
	return val, nil
}

func prettyPrint(str string) {
	fmt.Println(cutePrint.Render(str))
}
