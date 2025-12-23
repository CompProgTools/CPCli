package models

type Config struct {
	Name string `json:"name,omitempty"`
	Codeforces string `json:"codeforces,omitempty"`
	CodeforcesRating int `json:"codeforces_rating,omitempty"`
	LeetCode string `json:"leetcode,omitempty"`
	LeetCodeRating int `json:"leetcode_rating,omitempty"`
	PreferredLanguage string `json:"preferred_language,omitempty"`
	PreferredEditor string `json:"preferred_editor,omitempty"`
	TemplateOutputPath string `json:"template_output_path,omitempty"`
	OpenKattisUsername string `json:"openkattis_username,omitempty"`
	OpenKattisPassword string `json:"openkattis_password,omitempty"`
	Templates map[string]string `json:"templates,omitempty"`
}

type SolvedProblem struct {
	ContestID int `json:"contestId"`
	Index string `json:"index"`
	Name string `json:"name"`
	Tags []string `json:"tags"`
	Rating *int `json:"rating"`
	Timestamp string `json:"timestamp"`
}

type StreakMeta struct {
	HighScore int `json:"highscore"`
}

type CodeforcesContest struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Phase string `json:"phase"`
	RelativeTimeSeconds int `json:"relativeTimeSeconds"`
}

type CodeforcesProblem struct {
	ContestID int `json:"contestId"`
	Index string `json:"index"`
	Name string `json:"name"`
	Tags []string `json:"tags,omitempty"`
	Rating *int `json:"rating,omitempty"`
}

type LeetCodeDaily struct {
	Question struct {
		Title string `json:"title"`
		QuestionFrontentID string `json:"questionFrontentId"`
		Difficulty string `json:"difficulty"`
		Content string `json:"content"`
	} `json:"question"`

	Date string `json:"date"`
	Link string `json:"link"`
}