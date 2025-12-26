package subcommands

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	passStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")).
		Bold(true)

	failStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")).
		Bold(true)

	testHeaderStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("14")).
		Bold(true).
		Underline(true)

	testInfoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("12"))
)

type TestCase struct {
	Input string
	ExpectedOutput string
}

func RunTest(args []string) error {
	if len(args) == 0 {
		fmt.Println(errorStyle.Render("usage: kruskal test <filename.extension>"))
		fmt.Println(testInfoStyle.Render("\nexample: kruskal test solution.cpp"))
		return nil
	}

	filename := args[0]
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", filename)
	}

	ext := strings.ToLower(filepath.Ext(filename))

	validExtensions := map[string]bool {
		".cpp": true, ".cc": true, ".c": true,
		".py": true, ".java": true, ".go": true,
	}

	if !validExtensions[ext] {
		return fmt.Errorf("unsupported file extension: %s", ext)
	}

	testCases, err := collectTestCases()
	if err != nil {
		return err
	}

	if len(testCases) == 0 {
		fmt.Println(yellowStyle2.Render("no test cases provided"))
		return nil
	}

	fmt.Println()
	fmt.Println(testHeaderStyle.Render("running tests..."))
	fmt.Println()

	runCmd, cleanup, err := prepareExecution(filename, ext)
	if err != nil {
		return err
	}
	defer cleanup()

	passed := 0
	failed := 0

	for i, tc := range testCases {
		fmt.Printf("\n%s\n", testHeaderStyle.Render(fmt.Sprintf("TEST CASE #%d", i + 1)))

		result, err := executeTest(runCmd, tc, i + 1)
		if err != nil {
			fmt.Println(failStyle.Render(fmt.Sprintf("Error: %v", err)))
			failed++
			continue
		}

		if result {
			fmt.Println(passStyle.Render("PASSED"))
			passed++
		} else {
			failed++
		}
	}

	fmt.Println()
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("%s %d/%d passed\n", testHeaderStyle.Render("summary:"), passed, len(testCases))

	if failed > 0 {
		fmt.Printf("%s %d/%d failed\n", failStyle.Render("failed:"), failed, len(testCases))
	} else {
		fmt.Println(passStyle.Render("all tests passed!"))
	}
	fmt.Println(strings.Repeat("=", 50))

	return nil
}

func collectTestCases() ([]TestCase, error) {
	reader := bufio.NewReader(os.Stdin)
	
	fmt.Print("Enter number of test cases: ")
	numTestsStr, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	
	numTestsStr = strings.TrimSpace(numTestsStr)
	numTests, err := strconv.Atoi(numTestsStr)
	if err != nil || numTests <= 0 {
		return nil, fmt.Errorf("invalid number of test cases")
	}

	testCases := make([]TestCase, 0, numTests)

	for i := 0; i < numTests; i++ {
		fmt.Printf("\n%s\n", testHeaderStyle.Render(fmt.Sprintf("Test Case #%d", i+1)))
		
		fmt.Println("Enter input (end with empty line):")
		inputLines := []string{}
		emptyCount := 0
		
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			
			if strings.TrimSpace(line) == "" {
				emptyCount++
				if emptyCount >= 1 && len(inputLines) > 0 {
					break
				}
			} else {
				emptyCount = 0
				inputLines = append(inputLines, line)
			}
		}
		input := strings.Join(inputLines, "")

		fmt.Println("Enter expected output (end with empty line):")
		outputLines := []string{}
		emptyCount = 0
		
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			
			if strings.TrimSpace(line) == "" {
				emptyCount++
				if emptyCount >= 1 && len(outputLines) > 0 {
					break
				}
			} else {
				emptyCount = 0
				outputLines = append(outputLines, line)
			}
		}
		expectedOutput := strings.TrimSpace(strings.Join(outputLines, ""))

		testCases = append(testCases, TestCase{
			Input:          input,
			ExpectedOutput: expectedOutput,
		})
	}

	return testCases, nil
}

func prepareExecution(filename string, ext string) ([]string, func(), error) {
	var runCmd []string
	var cleanup func()
	exeName := "test_executable"

	switch ext {
	case ".cpp", ".cc":
		fmt.Println(testInfoStyle.Render("compiling c++ code..."))
		cmd := exec.Command("g++", "-std=c++17", "-O2", filename, "-o", exeName)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			return nil, nil, fmt.Errorf("compilation failed:\n%s", stderr.String())
		}

		fmt.Println(greenStyle.Render("compilation successful"))
		runCmd = []string{"./" + exeName}
		cleanup = func() {
			os.Remove(exeName)
		}

	case ".c":
		fmt.Println(testInfoStyle.Render("compiling c code..."))
		cmd := exec.Command("gcc", "-std=c11", "-O2", filename, "-o", exeName)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			return nil, nil, fmt.Errorf("compilation failed:\n%s", stderr.String())
		}

		fmt.Println(greenStyle.Render("compilation successful"))
		runCmd = []string{"./" + exeName}
		cleanup = func() {
			os.Remove(exeName)
		}

	case ".py":
		runCmd = []string{"python3", filename}
		cleanup = func() {}

	case ".java":
		fmt.Println(testInfoStyle.Render("compiling java code..."))
		cmd := exec.Command("javac", filename)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			return nil, nil, fmt.Errorf("compilation failed:\n%s", stderr.String())
		}

		fmt.Println(greenStyle.Render("compilation successful"))
		className := strings.TrimSuffix(filepath.Base(filename), ".java")
		runCmd = []string{"java", className}
		cleanup = func() {
			os.Remove(className + ".class")
		}

	case ".go":
		fmt.Println(testInfoStyle.Render("building go code"))
		cmd := exec.Command("go", "build", "-o", exeName, filename)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			return nil, nil, fmt.Errorf("build failed:\n%s", stderr.String())
		}

		fmt.Println(greenStyle.Render("build successful"))
		runCmd = []string{"./" + exeName}
		cleanup = func ()  {
			os.Remove(exeName)
		}

	default:
		return nil, nil, fmt.Errorf("unsupported file extension: %s", ext)
	}

	return runCmd, cleanup, nil
}

func executeTest(runCmd []string, tc TestCase, testNum int) (bool, error) {
	cmd := exec.Command(runCmd[0], runCmd[1:]...)
	cmd.Stdin = strings.NewReader(tc.Input)
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()

	select {
	case <-time.After(5 * time.Second):
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return false, fmt.Errorf("timeout (5 seconds exceeded)")
	case err := <-done:
		if err != nil {
			if stderr.Len() > 0 {
				return false, fmt.Errorf("runtime error:\n%s", stderr.String())
			}
			return false, err
		}
	}

	actualOutput := strings.TrimSpace(stdout.String())
	expectedOutput := strings.TrimSpace(tc.ExpectedOutput)

	if actualOutput == expectedOutput {
		return true, nil
	}

	fmt.Println(failStyle.Render("FAILED"))
	fmt.Println()
	fmt.Println(yellowStyle2.Render("Expected Output:"))
	fmt.Println(expectedOutput)
	fmt.Println()
	fmt.Println(errorStyle.Render("Actual Output:"))
	fmt.Println(actualOutput)
	fmt.Println()
	
	return false, nil
}