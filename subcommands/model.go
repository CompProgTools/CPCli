package subcommands

import (
	"fmt"
	"strings"
	"syscall"

	"github.com/CompProgTools/Kruskal/config"
	"github.com/CompProgTools/Kruskal/internal/models"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/textinput"
	"golang.org/x/term"
)

func RunConfig() error {
	return RunConfigInteractive()
}

func RunConfigInteractive() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	for {
		choices := []string{
			"set name",
			"set preferred language",
			"set preferred code editor",
			"set codeforces username",
			"set leetcode username",
			"set template output folder",
			"set openkattis credentials",
			"view current configuration",
			"exit",
		}

		sp := selection.New("configuration menu - what would you like to configure?", choices)
		sp.Filter = nil
		sp.PageSize = 10

		choice, err := sp.RunPrompt()
		if err != nil {
			return err
		}

		if choice == "exit" {
			fmt.Println(greenStyle.Render("configuration saved"))
			break
		}

		switch choice {
		case "set name":
			if err := setName(cfg); err != nil {
				return err
			}

		case "set preferred language":
			if err := setPreferredLanguage(cfg); err != nil {
				return err
			}

		case "set preferred code editor":
			if err := setPreferredEditor(cfg); err != nil {
				return err
			}

		case "set codeforces username":
			if err := setCodeforcesUsername(cfg); err != nil {
				return err
			}

		case "set leetcode username":
			if err := setLeetCodeUsername(cfg); err != nil {
				return err
			}

		case "set template output folder":
			if err := setTemplateOutputPath(cfg); err != nil {
				return err
			}

		case "set openkattis credentials":
			if err := setOpenKattisCredentials(cfg); err != nil {
				return err
			}

		case "view current configuration":
			displayCurrentConfig(cfg)
		}

		if choice != "view current configuration" && choice != "exit" {
			if err := config.SaveConfig(cfg); err != nil {
				return err
			}
		}

		fmt.Println()
	}

	return nil
}

func setName(cfg *models.Config) error {
	currentValue := ""
	if cfg.Name != "" {
		currentValue = fmt.Sprintf(" (current: %s)", cfg.Name)
	}

	input := textinput.New(fmt.Sprintf("enter your name%s:", currentValue))
	input.Placeholder = "your name"
	if cfg.Name != "" {
		input.InitialValue = cfg.Name
	}

	name, err := input.RunPrompt()
	if err != nil {
		return err
	}

	name = strings.TrimSpace(name)
	if name == "" {
		fmt.Println(yellowStyle2.Render("name cannot be empty. keeping previous value."))
		return nil
	}

	cfg.Name = name
	fmt.Println(greenStyle.Render(fmt.Sprintf("name set to: %s", name)))
	return nil
}

func setPreferredLanguage(cfg *models.Config) error {
	langChoices := []string{"c++", "c", "python", "java", "go", "other"}

	langSp := selection.New("choose your preferred programming language:", langChoices)
	langSp.Filter = nil

	lang, err := langSp.RunPrompt()
	if err != nil {
		return err
	}

	if lang == "other" {
		input := textinput.New("enter your preferred language:")
		customLang, err := input.RunPrompt()
		if err != nil {
			return err
		}
		lang = strings.TrimSpace(customLang)
	}

	cfg.PreferredLanguage = lang
	fmt.Println(greenStyle.Render(fmt.Sprintf("preferred language set to: %s", lang)))
	return nil
}

func setPreferredEditor(cfg *models.Config) error {
	editorMap := map[string]string{
		"vscode":         "code",
		"neovim":         "nvim",
		"vim":            "vim",
		"sublime text":   "subl",
		"atom":           "atom",
		"emacs":          "emacs",
		"nano":           "nano",
		"custom command": "",
	}

	editorChoices := []string{
		"vscode", "neovim", "vim", "sublime text",
		"atom", "emacs", "nano", "custom command",
	}

	editorSp := selection.New("choose your preferred code editor:", editorChoices)
	editorSp.Filter = nil

	editorChoice, err := editorSp.RunPrompt()
	if err != nil {
		return err
	}

	if editorChoice == "custom command" {
		input := textinput.New("enter the terminal command for your editor:")
		input.Placeholder = "e.g., code, vim, emacs"

		custom, err := input.RunPrompt()
		if err != nil {
			return err
		}

		custom = strings.TrimSpace(custom)
		cfg.PreferredEditor = custom
		fmt.Println(greenStyle.Render(fmt.Sprintf("preferred editor set to: %s", custom)))
	} else {
		cfg.PreferredEditor = editorMap[editorChoice]
		fmt.Println(greenStyle.Render(fmt.Sprintf("preferred editor set to: %s (%s)", editorChoice, editorMap[editorChoice])))
	}

	return nil
}

func setCodeforcesUsername(cfg *models.Config) error {
	currentValue := ""
	if cfg.Codeforces != "" {
		currentValue = fmt.Sprintf(" (current: %s)", cfg.Codeforces)
	}

	input := textinput.New(fmt.Sprintf("enter codeforces username%s:", currentValue))
	input.Placeholder = "username"
	if cfg.Codeforces != "" {
		input.InitialValue = cfg.Codeforces
	}

	handle, err := input.RunPrompt()
	if err != nil {
		return err
	}

	handle = strings.TrimSpace(handle)
	if handle == "" {
		fmt.Println(yellowStyle2.Render("username cannot be empty. keeping previous value."))
		return nil
	}

	cfg.Codeforces = handle
	fmt.Println(greenStyle.Render(fmt.Sprintf("codeforces username set to: %s", handle)))
	return nil
}

func setLeetCodeUsername(cfg *models.Config) error {
	currentValue := ""
	if cfg.LeetCode != "" {
		currentValue = fmt.Sprintf(" (current: %s)", cfg.LeetCode)
	}

	input := textinput.New(fmt.Sprintf("enter leetcode username%s:", currentValue))
	input.Placeholder = "username"
	if cfg.LeetCode != "" {
		input.InitialValue = cfg.LeetCode
	}

	handle, err := input.RunPrompt()
	if err != nil {
		return err
	}

	handle = strings.TrimSpace(handle)
	if handle == "" {
		fmt.Println(yellowStyle2.Render("username cannot be empty. keeping previous value."))
		return nil
	}

	cfg.LeetCode = handle
	fmt.Println(greenStyle.Render(fmt.Sprintf("leetcode username set to: %s", handle)))
	return nil
}

func setTemplateOutputPath(cfg *models.Config) error {
	currentValue := ""
	if cfg.TemplateOutputPath != "" {
		currentValue = fmt.Sprintf(" (current: %s)", cfg.TemplateOutputPath)
	}

	input := textinput.New(fmt.Sprintf("enter absolute path for template outputs%s:", currentValue))
	input.Placeholder = "/path/to/folder"
	if cfg.TemplateOutputPath != "" {
		input.InitialValue = cfg.TemplateOutputPath
	}

	path, err := input.RunPrompt()
	if err != nil {
		return err
	}

	path = strings.TrimSpace(path)
	if path == "" {
		fmt.Println(yellowStyle2.Render("path cannot be empty. keeping previous value."))
		return nil
	}

	cfg.TemplateOutputPath = path
	fmt.Println(greenStyle.Render(fmt.Sprintf("template output path set to: %s", path)))
	return nil
}

func setOpenKattisCredentials(cfg *models.Config) error {
	userInput := textinput.New("enter openkattis username:")
	username, err := userInput.RunPrompt()
	if err != nil {
		return err
	}

	username = strings.TrimSpace(username)
	if username == "" {
		fmt.Println(yellowStyle2.Render("username cannot be empty."))
		return nil
	}

	fmt.Print("enter openkattis password (hidden): ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return err
	}
	fmt.Println()
	password := string(bytePassword)

	if password == "" {
		fmt.Println(yellowStyle2.Render("password cannot be empty."))
		return nil
	}

	cfg.OpenKattisUsername = username
	cfg.OpenKattisPassword = password
	fmt.Println(greenStyle.Render("openkattis credentials saved."))
	return nil
}

func displayCurrentConfig(cfg *models.Config) {
	fmt.Println()
	fmt.Println(configTitleStyle.Render("current configuration"))
	fmt.Println(strings.Repeat("=", 50))

	displayConfigField("name", cfg.Name)
	displayConfigField("preferred language", cfg.PreferredLanguage)
	displayConfigField("preferred editor", cfg.PreferredEditor)
	displayConfigField("codeforces username", cfg.Codeforces)
	displayConfigField("leetcode username", cfg.LeetCode)
	displayConfigField("template output path", cfg.TemplateOutputPath)
	displayConfigField("openkattis username", cfg.OpenKattisUsername)

	if cfg.OpenKattisPassword != "" {
		fmt.Printf("%-25s: %s\n",
			configItemStyle.Render("openkattis password"),
			configValueStyle.Render("***hidden***"))
	} else {
		fmt.Printf("%-25s: %s\n",
			configItemStyle.Render("openkattis password"),
			yellowStyle2.Render("not set"))
	}

	if len(cfg.Templates) > 0 {
		fmt.Printf("%-25s: %s\n",
			configItemStyle.Render("templates"),
			configValueStyle.Render(fmt.Sprintf("%d template(s)", len(cfg.Templates))))
	} else {
		fmt.Printf("%-25s: %s\n",
			configItemStyle.Render("templates"),
			yellowStyle2.Render("none"))
	}

	fmt.Println(strings.Repeat("=", 50))
	fmt.Println()
}

func displayConfigField(fieldName string, value string) {
	if value != "" {
		fmt.Printf("%-25s: %s\n",
			configItemStyle.Render(fieldName),
			configValueStyle.Render(value))
	} else {
		fmt.Printf("%-25s: %s\n",
			configItemStyle.Render(fieldName),
			yellowStyle2.Render("not set"))
	}
}