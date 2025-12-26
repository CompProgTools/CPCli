package subcommands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/CompProgTools/Kruskal/internal/models"
	"github.com/CompProgTools/Kruskal/config"
	"github.com/charmbracelet/lipgloss"
)

var (
	templateTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		Underline(true)
)

func RunTemplate(args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	templatesPath := filepath.Join(config.GetConfigPath(), "templates")
	if err := os.MkdirAll(templatesPath, 0755); err != nil {
		return err
	}

	if len(args) == 0 {
		return printTemplateUsage()
	}

	if contains(args, "--list") {
		return listTemplates(cfg)
	}

	if contains(args, "--delete") {
		return deleteTemplate(args, cfg, templatesPath)
	}

	if contains(args, "--use") && contains(args, "--filename") {
		return useTemplate(args, cfg, templatesPath)
	}

	if contains(args, "--make") && contains(args, "--alias") {
		return makeTemplate(args, cfg, templatesPath)
	}

	return printTemplateUsage()
}

func printTemplateUsage() error {
	usage := `
	Template Management Commands:
		
		Create a new template:
			kruskal template --make <filename.ext> --alias <alias>
		
		Use a template:
			kruskal template --use <alias> --filename <output.ext>
		
		List all templates:
			kruskal template --list
		
		Delete a template:
			kruskal template --delete <alias>
	
	Examples:
		kruskal template --make template.cpp --alias cpp
		kruskal template --use cpp --filename solution.cpp
		kruskal template --list
		kruskal template --delete cpp
`
	fmt.Println(infoStyle.Render(usage))
	return nil
}

func listTemplates(cfg *models.Config) error {
	if len(cfg.Templates) == 0 {
		fmt.Println(yellowStyle2.Render("no templates found. use `--make` to create one."))
		return nil
	}

	fmt.Println(templateTitleStyle.Render("saved templates"))
	fmt.Println()

	maxAliasLen := 0
	for alias := range cfg.Templates {
		if len(alias) > maxAliasLen {
			maxAliasLen = len(alias)
		}
	}

	for alias, filename := range cfg.Templates {
		fmt.Printf("	%s %s %s\n", 
			greenStyle.Render(fmt.Sprintf("%-*s", maxAliasLen, alias)), 
			infoStyle.Render("->"), 
			filename)
	}

	fmt.Println()
	fmt.Printf("total templates: %d\n", len(cfg.Templates))

	return nil
}

func deleteTemplate(args []string, cfg *models.Config, templatesPath string) error {
	idx := indexOf(args, "--delete")
	if idx == -1 || idx + 1 >= len(args) {
		fmt.Println(errorStyle.Render("missing alias for --delete"))
		return nil
	}

	alias := args[idx + 1]

	if cfg.Templates == nil || cfg.Templates[alias] == "" {
		fmt.Println(errorStyle.Render(fmt.Sprintf("template alias '%s' not found", alias)))
		return nil
	}

	filename := cfg.Templates[alias]
	templatePath := filepath.Join(templatesPath, filename)

	if err := os.Remove(templatePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete template file: %v", err)
	}

	delete(cfg.Templates, alias)
	if err := config.SaveConfig(cfg); err != nil {
		return err
	}

	fmt.Println(greenStyle.Render(fmt.Sprintf("template '%s' deleted successfully", alias)))
	return nil
}

func useTemplate(args []string, cfg *models.Config, templatesPath string) error {
	aliasIdx := indexOf(args, "--use")
	fileIdx := indexOf(args, "--filename")

	if aliasIdx == -1 || fileIdx == -1 || aliasIdx + 1 >= len(args) || fileIdx + 1 >= len(args) {
		fmt.Println(errorStyle.Render("missing value for --use or --filename"))
		return nil
	}

	alias := args[aliasIdx + 1]
	newFilename := args[fileIdx + 1]

	outputDir := cfg.TemplateOutputPath
	if outputDir == "" {
		fmt.Println(errorStyle.Render("template output path not set. use 'kruskal config' to set it"))
		return nil
	}

	if cfg.Templates == nil || cfg.Templates[alias] == "" {
		fmt.Println(errorStyle.Render(fmt.Sprintf("template alias '%s' not found", alias)))
		fmt.Println(infoStyle.Render("use 'kruskal template --list' to see available templates"))
		return nil
	}

	sourceTemplate := filepath.Join(templatesPath, cfg.Templates[alias])
	destFile := filepath.Join(outputDir, newFilename)

	if _, err := os.Stat(sourceTemplate); os.IsNotExist(err) {
		return fmt.Errorf("template file not found: %s", sourceTemplate)
	}

	if _, err := os.Stat(destFile); err == nil {
		fmt.Printf("file %s already exists. overwrite? (y/n): ", newFilename)
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" {
			fmt.Println(yellowStyle2.Render("operation cancelled"))
			return nil
		}
	}

	content, err := os.ReadFile(sourceTemplate)
	if err != nil {
		return fmt.Errorf("failed to read template: %v", err)
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	if err := os.WriteFile(destFile, content, 0644); err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}

	fmt.Println(greenStyle.Render(fmt.Sprintf("created %s from template '%s'", destFile, alias)))

	editor := cfg.PreferredEditor
	if editor == "" {
		editor = getDefaultEditor()
	}

	if err := openInEditor(editor, destFile); err != nil {
		fmt.Println(yellowStyle2.Render(fmt.Sprintf("created file but couldn't open editor: %v", err)))
		fmt.Println(infoStyle.Render("set your preferred editor with 'kruskal config'"))
	} else {
		fmt.Println(infoStyle.Render(fmt.Sprintf("opened in %s", editor)))
	}

	return nil
}

func makeTemplate(args []string, cfg *models.Config, templatesPath string) error {
	nameIdx := indexOf(args, "--make")
	aliasIdx := indexOf(args, "--alias")

	if nameIdx == -1 || aliasIdx == -1 || nameIdx + 1 >= len(args) || aliasIdx + 1 >= len(args) {
		fmt.Println(errorStyle.Render("missing value for --make or --alias"))
		return nil
	}

	filename := args[nameIdx + 1]
	alias := args[aliasIdx + 1]

	if !strings.Contains(filename, ".") {
		return fmt.Errorf("filename must include an extension (e.g., template.cpp)")
	}

	templateFilePath := filepath.Join(templatesPath, filename)

	fileExists := false
	if _, err := os.Stat(templateFilePath); err == nil {
		fileExists = true
		fmt.Println(yellowStyle2.Render(fmt.Sprintf("template file '%s' already exists. opening it...", filename)))
	} else {
		ext := strings.ToLower(filepath.Ext(filename))
		boilerplate := getBoilerplate(ext)

		if err := os.WriteFile(templateFilePath, []byte(boilerplate), 0644); err != nil {
			return fmt.Errorf("failed to create template: %v", err)
		}

		fmt.Println(greenStyle.Render(fmt.Sprintf("created template file: %s", filename)))
	}

	if cfg.Templates == nil {
		cfg.Templates = make(map[string]string)
	}

	if existingFile, exists := cfg.Templates[alias]; exists && existingFile != filename {
		fmt.Printf("alias '%s' already exists for '%s'. overwrite? (y/n): ", alias, existingFile)
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" {
			fmt.Println(yellowStyle2.Render("operation cancelled."))
			return nil
		}
	}

	cfg.Templates[alias] = filename
	if err := config.SaveConfig(cfg); err != nil {
		return err
	}

	if !fileExists {
		fmt.Println(greenStyle.Render(fmt.Sprintf("template alias '%s' created", alias)))
	}

	editor := cfg.PreferredEditor
	if editor == "" {
		editor = getDefaultEditor()
	}

	if err := openInEditor(editor, templateFilePath); err != nil {
		fmt.Println(yellowStyle2.Render(fmt.Sprintf("template created but couldn't open editor: %v", err)))
		fmt.Println(infoStyle.Render(fmt.Sprintf("edit manually: %s", templateFilePath)))
	} else {
		fmt.Println(infoStyle.Render(fmt.Sprintf("opened in %s", editor)))
	}

	return nil
}

func getBoilerplate(ext string) string {
	switch ext {
	case ".cpp", ".cc":
		return `#include <bits/stdc++.h>
using namespace std;

int main() {
    ios_base::sync_with_stdio(false);
    cin.tie(NULL);
    
    
    
    return 0;
}
`

	case ".c":
		return `#include <stdio.h>
#include <stdlib.h>

int main() {
    
    
    return 0;
}
`

	case ".py":
		return `#!/usr/bin/env python3

def main():
    pass

if __name__ == "__main__":
    main()
`

	case ".java":
		return `import java.util.*;
import java.io.*;

public class Solution {
    public static void main(String[] args) {
        Scanner sc = new Scanner(System.in);
        
        
        
        sc.close();
    }
}
`

	case ".go":
		return `package main

import (
	"fmt"
)

func main() {
	
}
`
	
	default:
		return ""
	}
}

func getDefaultEditor() string {
	switch runtime.GOOS {
	case "windows":
		return "notepad"
	case "darwin":
		return "open -e"
	default:
		return "nano"
	}
}

func openInEditor(editor string, filePath string) error {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		if strings.Contains(strings.ToLower(editor), "code") {
			vscPath := filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", "Microsoft VS Code", "Code.exe")
			if _, err := os.Stat(vscPath); err == nil {
				cmd = exec.Command(vscPath, filePath)
			} else {
				cmd = exec.Command("cmd", "/c", "start", editor, filePath)
			}
		} else if editor == "notepad" {
			cmd = exec.Command(editor, filePath)
		} else {
			cmd = exec.Command("cmd", "/c", "start", editor, filePath)
		}
	} else if runtime.GOOS == "darwin" && strings.HasPrefix(editor, "open") {
		parts := strings.Split(editor, " ")
		args := append(parts[1:], filePath)
		cmd = exec.Command(parts[0], args...)
	} else {
		cmd = exec.Command(editor, filePath)
	}

	return cmd.Start()
}