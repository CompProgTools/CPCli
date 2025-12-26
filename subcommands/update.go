package subcommands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func RunUpdate(args []string) error {
	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	root := filepath.Dir(execPath)

	gitDir := filepath.Join(root, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		fmt.Println(yellowStyle2.Render("this isn't a git repository. a manual update is needed"))
		return nil
	}

	cmd := exec.Command("git", "pull")
	cmd.Dir = root

	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println(errorStyle.Render("failed to update CPCli"))
		fmt.Println(string(output))
		return err
	}

	fmt.Println(greenStyle.Render("CPCli successfully updated!"))
	fmt.Println(string(output))

	return nil
}