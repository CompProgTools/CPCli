package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/CompProgTools/Kruskal/config"
	"github.com/CompProgTools/Kruskal/subcommands"
	"github.com/charmbracelet/lipgloss"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/textinput"
	"github.com/pkg/browser"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("12")).
			MarginBottom(1)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9"))

	yellowStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("11"))
)

func linkAccount() error {
	platforms := []string{"Codeforces", "LeetCode", "Back"}

	sp := selection.New("select a platform to link: ", platforms)
	sp.Filter = nil

	platform, err := sp.RunPrompt()
	if err != nil {
		return err
	}

	if platform == "Back" {
		return nil
	}

	for {
		input := textinput.New(fmt.Sprintf("enter your %s username: ", platform))
		input.Placeholder = "username"

		handle, err := input.RunPrompt()
		if err != nil {
			return err
		}

		valid := false
		if platform == "Codeforces" {
			valid, err = subcommands.ValidateCodeforcesUser(handle)
		} else if platform == "LeetCode" {
			valid, err = subcommands.ValidateLeetCodeUser(handle)
		}

		if err != nil {
			fmt.Println(errorStyle.Render(fmt.Sprintf("error connecting to %s: %v", platform, err)))
			continue
		}

		if !valid {
			fmt.Println(errorStyle.Render(fmt.Sprintf("user not found on %s. Try again.", platform)))
			continue
		}

		platformLower := strings.ToLower(platform)
		if err := config.SetAccount(platformLower, handle); err != nil {
			return err
		}

		break
	}

	fmt.Println(successStyle.Render(fmt.Sprintf("%s account linked successfully!", platform)))

	return nil
}

func showMenu() error {
	allLinked, err := config.IsAllLinked()
	if err != nil {
		return err
	}

	options := []string{"View repository"}
	if !allLinked {
		options = append(options, "Link Account")
	}
	options = append(options, "coming soon...", "exit")

	fmt.Println(titleStyle.Render("Hi! This is Kruskal, a command line interface for competitive programmers!"))

	sp := selection.New("what would you like to do?", options)
	sp.Filter = nil

	choice, err := sp.RunPrompt()
	if err != nil {
		return err
	}

	switch choice {
	case "View repository":
		fmt.Println(successStyle.Render("opening Kruskal GitHub repository"))
		_ = browser.OpenURL("https://github.com/CompProgTools/kruskal")
	case "Link Account":
		return linkAccount()
	case "Coming soon...":
		fmt.Println(yellowStyle.Render("stay tuned!"))
	case "Exit":
		fmt.Println(errorStyle.Render("goodbye!"))
	}

	return nil
}

func main() {
	if len(os.Args) > 1 {
		subcommand := os.Args[1]
		args := os.Args[2:]

		var err error
		switch subcommand {
		case "sync":
			err = subcommands.RunSync()
		case "streak":
			err = subcommands.RunStreak()
		case "stats":
			err = subcommands.RunStats()
		case "graph":
			err = subcommands.RunGraph()
		case "test":
			err = subcommands.RunTest(args)
		case "config":
			err = subcommands.RunConfigInteractive()
		case "template":
			err = subcommands.RunTemplate(args)
		case "daily":
			err = subcommands.RunDaily()
		case "cf":
			err = subcommands.RunCF(args)
		case "update":
			err = subcommands.RunUpdate(args)
		default:
			fmt.Fprint(os.Stderr, errorStyle.Render(fmt.Sprintf("unknown command: %s\n", subcommand)))
			os.Exit(1)
		}

		if err != nil {
			fmt.Fprint(os.Stderr, errorStyle.Render(fmt.Sprintf("error: %v\n", err)))
			os.Exit(1)
		}
	} else {
		if err := showMenu(); err != nil {
			fmt.Fprintf(os.Stderr, "Error %v\n", err)
			os.Exit(1)
		}
	}
}