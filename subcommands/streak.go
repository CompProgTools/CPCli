package subcommands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/CompProgTools/Kruskal/config"
	"github.com/CompProgTools/Kruskal/internal/models"
	"github.com/charmbracelet/lipgloss"
)

var (
	cyanBold = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("14"))
	magentaBold = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("13"))
)

func RunStreak() error {
	solvedPath := filepath.Join(config.GetConfigPath(), "solved.json")

	if _, err := os.Stat(solvedPath); os.IsNotExist(err) {
		fmt.Println(yellowStyle2.Render("You haven't solved any problems yet."))
		return nil
	}

	data, err := os.ReadFile(solvedPath)
	if err != nil {
		return err
	}

	var solvedData []models.SolvedProblem
	if err := json.Unmarshal(data, &solvedData); err != nil {
		return err
	}

	if len(solvedData) == 0 {
		fmt.Println(yellowStyle2.Render("You haven't solved any problems yet."))
		return nil
	}

	dateSet := make(map[string]bool)
	solvedByDate := make(map[string][]models.SolvedProblem)

	for _, entry := range solvedData {
		dateStr := entry.Timestamp[:10]
		dateSet[dateStr] = true
		solvedByDate[dateStr] = append(solvedByDate[dateStr], entry)
	}

	//sort the dates
	var dates []string
	for date := range dateSet {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	var dateObjs []time.Time
	for _, dateStr := range dates {
		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}
		dateObjs = append(dateObjs, t)
	}

	today := time.Now().Truncate(24 * time.Hour)
	streak := 0
	checkDay := today

	for {
		found := false
		for _, d := range dateObjs {
			if d.Equal(checkDay) {
				found = true
				break
			}
		}
		
		if !found {
			break
		}

		streak++
		checkDay = checkDay.AddDate(0, 0, -1)
	}

	streakMetaPath := filepath.Join(config.GetConfigPath(), "streak_meta.json")
	highScore := 0

	if _, err := os.Stat(streakMetaPath); err == nil {
		data, err := os.ReadFile(streakMetaPath)
		if err == nil {
			var meta models.StreakMeta
			if err := json.Unmarshal(data, &meta); err == nil {
				highScore = meta.HighScore
			}
		}
	}

	if streak > highScore {
		highScore = streak
		meta := models.StreakMeta{HighScore: highScore}
		data, _ := json.MarshalIndent(meta, "", " ")
		_ = os.WriteFile(streakMetaPath, data, 0644)
	}

	fmt.Println(cyanBold.Render(fmt.Sprintf("Your current streak: %d day(s)", streak)))
	fmt.Println(magentaBold.Render(fmt.Sprintf("Your all-time highest streak: %d day(s)", highScore)))

	//show yesterdays problem if a streak exists
	if streak > 0 {
		yesterday := today.AddDate(0, 0, -1).Format("2006-01-02")
		if probs, exists := solvedByDate[yesterday]; exists {
			fmt.Println("\nProblems solved yesterday:")
			for _, prob := range probs {
				tags := ""
				if len(prob.Tags) > 0 {
					tags = " [" + prob.Tags[0]
					for i := 1; i < len(prob.Tags); i++ {
						tags += ", " + prob.Tags[i]
					}
					tags += "]"
				}
				fmt.Printf(" . %d/%s - %s%s\n", prob.ContestID, prob.Index, prob.Name, tags)
			}
		} else {
			fmt.Println(yellowStyle2.Render("You didn't solve any problems yesterday"))
		}
	} else {
		fmt.Println(yellowStyle2.Render("No current streak. Get solving to start one!"))
	}

	return nil
}