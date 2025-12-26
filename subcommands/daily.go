package subcommands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/CompProgTools/Kruskal/internal/models"
	"github.com/charmbracelet/lipgloss"
)

var (
	panelStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("14")).
		Padding(1, 2)

	titleDaily = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("14"))
)

func RunDaily() error {
	url := "https://leetcode-api-pied.vercel.app/daily"

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var data models.LeetCodeDaily
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	title := data.Question.Title
	frontendID := data.Question.QuestionFrontentID
	date := data.Date
	difficulty := data.Question.Difficulty
	link := "https://leetcode/com" + data.Link

	fmt.Println(panelStyle.Render(titleDaily.Render(fmt.Sprintf("LeetCode Daily Challenge - %s", date))))
	fmt.Println()
	fmt.Printf("%s %s. %s (%s)\n", lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true).Render("Problem:"), frontendID, title, difficulty)
	fmt.Printf("%s %s\n", lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Render("Link:"), link)
	fmt.Println()
	fmt.Println(data.Question.Content)

	return nil
}