package subcommands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/CompProgTools/Kruskal/internal/models"
	"github.com/CompProgTools/Kruskal/config"
	"github.com/charmbracelet/lipgloss"
)

var (
	tableHeaderStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("14"))
)

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func indexOf(slice []string, item string) int {
	for i, s := range slice {
		if s == item {
			return i
		}
	}
	return -1
}

func RunCF(args []string) error {
	if contains(args, "--list") {
		return listContests()
	}

	if contains(args, "--solved") {
		idx := indexOf(args, "--solved")
		if idx == -1 || idx + 1 >= len(args) {
			fmt.Println(errorStyle.Render("missing argument for --solved flag"))
			return nil
		}
		return logSolved(args[idx + 1])
	}

	fmt.Println(errorStyle.Render("nnknown flag. use --list or --solved"))
	return nil
}

func listContests() error {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("http://codeforces.com/api/contest.list")
	if err != nil {
		return err 
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println(errorStyle.Render("response error, status 404 or empty body"))
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result struct {
		Status string `json:"status"`
		Result []models.CodeforcesContest `json:"result"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	if result.Status != "OK" {
		fmt.Println(errorStyle.Render("invalid API status"))
		return nil
	}

	fmt.Println(tableHeaderStyle.Render("upcoming codeforces contests"))
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("%-10s %-50s %-15s\n", "ID", "name", "start in (sec)")
	fmt.Println(strings.Repeat("-", 80))

	for _, contest := range result.Result {
		if contest.Phase != "BEFORE" {
			break
		}
		fmt.Printf("%-10d %-50s %-15d\n", contest.ID, os.Truncate(contest.Name, 48), -contest.RelativeTimeSeconds)
	}
	
	return nil
}

func logSolved(contestAndProblem string) error {
	parts := strings.Split(contestAndProblem, "/")
	if len(parts) != 2 {
		fmt.Println(errorStyle.Render("usage: cpcli cf --solved <contestId>/<problemLetter>"))
		return nil
	}

	contestId, err := strconv.Atoi(parts[0])
	if err != nil {
		fmt.Println(errorStyle.Render("invalid contest ID"))
		return nil
	}

	problemLetter := strings.ToUpper(parts[1])

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://codeforces.com/api/problemset.problems")
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result struct {
		Status string `json:"status"`
		Result struct {
			Problems []models.CodeforcesProblem `json:"problems"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	if result.Status != "OK" {
		fmt.Println(errorStyle.Render("failed to fetch problemset from codeforces"))
		return nil
	}

	var problemInfo *models.CodeforcesProblem
	for i := range result.Result.Problems {
		p := &result.Result.Problems[i]
		if p.ContestID == contestId && p.Index == problemLetter {
			problemInfo = p
			break
		}
	}

	if problemInfo == nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("problem %d/%s not found in the problemset", contestId, problemLetter)))
		return nil
	}

	solvedPath := filepath.Join(config.GetConfigPath(), "solved.json")
	var solvedData []models.SolvedProblem

	if _, err := os.Stat(solvedPath); err == nil {
		data, err := os.ReadFile(solvedPath)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(data, &solvedData); err != nil {
			return err
		}
	}

	//new entries
	entry := models.SolvedProblem {
		ContestID: contestId,
		Index: problemLetter,
		Name: problemInfo.Name,
		Tags: problemInfo.Tags,
		Rating: problemInfo.Rating,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	solvedData = append(solvedData, entry)

	data, err := json.MarshalIndent(solvedData, "", " ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(solvedPath, data, 0644); err != nil {
		return err
	}

	fmt.Println(greenStyle.Render(fmt.Sprintf("logged solved problem: %d/%s - %s", contestId, problemLetter, problemInfo.Name)))

	return nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}

	return s[:maxLen - 3] + "..."
}