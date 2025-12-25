package subcommands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/CompProgTools/Kruskal/src/config"
	"github.com/charmbracelet/lipgloss"
)

var (
	greenStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	cyanStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	yellowStyle2 = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
)

func ValidateCodeforcesUser(handle string) (bool, error) {
	url := fmt.Sprintf("https://codeforces.com/api/user.rating?handle=%s", handle)

	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, err
	}

	return result["status"] == "OK", nil
}

func ValidateLeetCodeUser(handle string) (bool, error) {
	url := fmt.Sprintf("https://leetcode-api.pied.vercel.app/user/%s", handle)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)

	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, err
	}

	detail, ok := result["detail"].(string)
	return !ok || detail != "404: User not found", nil
}

func fetchRating(platform string, handle string) (int, error) {
	var url string
	
	if platform == "leetcode" {
		url = fmt.Sprintf("https://leetcode-api-pied.vercel.app/user/%s/contest", handle)
	} else if platform == "codeforces" {
		url = fmt.Sprintf("https://codeforces.com/api/user.rating?handle=%s", handle)
	} else {
		return 0, fmt.Errorf("unknown platform: %s", platform)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, err
	}

	return extractRating(platform, data)
}

func extractRating(platform string, data map[string]interface{}) (int, error) {
	if platform == "leetcode" {
		if userRank, ok := data["userContestRating"].(map[string]interface{}); ok {
			if rating, ok := userRank["rating"].(float64); ok {
				return int(rating), nil
			}
		}
	} else if platform == "codeforces" {
		if result, ok := data["result"].([]interface{}); ok && len(result) > 0 {
			lastContest := result[len(result) - 1].(map[string]interface{})
			if rating, ok := lastContest["newRating"].(float64); ok {
				return int(rating), nil
			}
		}
	}
	
	return 0, fmt.Errorf("could not extract rating")
}

func RunSync() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	updated := false
	
	platforms := map[string]string {
		"codeforces": cfg.Codeforces,
		"leetcode": cfg.LeetCode,
	}

	for platform, handle := range platforms {
		if handle == "" {
			continue
		}

		var currentRating int
		if platform == "codeforces" {
			currentRating = cfg.CodeforcesRating
		} else {
			currentRating = cfg.LeetCodeRating
		}

		newRating, err := fetchRating(platform, handle)
		if err != nil {
			continue
		}

		if currentRating == 0 {
			if platform == "codeforces" {
				cfg.CodeforcesRating = newRating
			} else {
				cfg.LeetCodeRating = newRating
			}
			fmt.Println(greenStyle.Render(fmt.Sprintf("%s (%s) rating set to %d", platform, handle, newRating)))
			updated = true
		} else if newRating != currentRating {
			diff := newRating - currentRating
			changeText := "increased"
			if diff < 0 {
				changeText = "decreased"
				diff = -diff
			}

			if platform == "codeforces" {
				cfg.CodeforcesRating = newRating
			} else {
				cfg.LeetCodeRating = newRating
			}

			fmt.Println(cyanStyle.Render(fmt.Sprintf("%s (%s) rating %s by %d to %d", platform, handle, changeText, diff, newRating)))
			updated = true
		} else {
			fmt.Println(yellowStyle2.Render(fmt.Sprintf("%s (%s) rating unchanged (%d)", platform, handle, newRating)))
		}
	}

	if updated {
		return config.SaveConfig(cfg)
	}

	return nil
}