package subcommands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/CompProgTools/Kruskal/config"
	"github.com/CompProgTools/Kruskal/internal/models"
)

func RunStats() error {
	solvedPath := filepath.Join(config.GetConfigPath(), "solved.json")

	if _, err := os.Stat(solvedPath); os.IsNotExist(err) {
		fmt.Println(errorStyle.Render("No solved problems logged yet."))
		return nil
	}

	data, err := os.ReadFile(solvedPath)
	if err != nil {
		return err
	}

	var solved []models.SolvedProblem
	if err := json.Unmarshal(data, &solved); err != nil {
		return err
	}

	if len(solved) == 0 {
		fmt.Println(errorStyle.Render("No solved problems logged yet."))
		return nil
	}

	totalProblems := len(solved)

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

	totalRating := 0
	ratingCount := 0
	for _, entry := range solved {
		if entry.Rating != nil && *entry.Rating > 0 {
			totalRating += *entry.Rating
			ratingCount++
		}
	}

	averageRating := 0.0
	if ratingCount > 0 {
		averageRating = float64(totalRating) / float64(ratingCount)
	}

	var highestProblem *models.SolvedProblem
	maxRating := -1

	for i := range solved {
		if solved[i].Rating != nil && *solved[i].Rating > maxRating {
			maxRating = *solved[i].Rating
			highestProblem = &solved[i]
		}
	}

	tagCounter := make(map[string]int)
	for _, entry := range solved {
		for _, tag := range entry.Tags {
			tagCounter[tag]++
		}
	}

	type tagCount struct {
		tag string
		count int
	}

	var tagCounts []tagCount
	for tag, count := range tagCounter {
		tagCounts = append(tagCounts, tagCount{tag, count})
	}

	for i := 0; i < len(tagCounts) - 1; i++ {
		for j := i + 1; j < len(tagCounts); j++ {
			if tagCounts[j].count > tagCounts[i].count {
				tagCounts[i], tagCounts[j] = tagCounts[j], tagCounts[i]
			}
		}
	}

	fmt.Println(ruleStyle.Render("═══════════════════════════════════════"))
	fmt.Println(ruleStyle.Render("           CPCli Stats"))
	fmt.Println(ruleStyle.Render("═══════════════════════════════════════"))
	fmt.Println()

	fmt.Printf("%s %d\n", boldStyle.Render("Highest Streak:"), highScore)
	fmt.Printf("%s %d\n", boldStyle.Render("Total Problems Logged:"), totalProblems)
	fmt.Printf("%s %.2f\n", boldStyle.Render("Average Rating:"), averageRating)

	if highestProblem != nil {
		fmt.Printf("%s %s (Rating: %d)\n", boldStyle.Render("Highest Rating Problem:"), highestProblem.Name, *highestProblem.Rating)
	} else {
		fmt.Printf("%s N/A\n", boldStyle.Render("Highest Rating Problem:"))
	}

	if len(tagCounts) > 0 {
		fmt.Print(boldStyle.Render("Top 3 Tags: "))
		topN := 3
		if len(tagCounts) < 3 {
			topN = len(tagCounts)
		}
		for i := 0; i < topN; i++ {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("%s (%d)", tagCounts[i].tag, tagCounts[i].count)
		}
		fmt.Println()
	} else {
		fmt.Printf("%s N/A\n", boldStyle.Render("Top 3 Tags:"))
	}

	return nil
} 