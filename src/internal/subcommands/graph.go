package subcommands

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/CompProgTools/CPCli/src/config"
	"github.com/CompProgTools/CPCli/src/internal/models"
	"github.com/charmbracelet/lipgloss"
	"github.com/pkg/browser"
)

var (
	infoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("12"))
)

type DataPoint struct {
	Date string
	Rating int
}

func RunGraph() error {
	solvedPath := filepath.Join(config.GetConfigPath(), "solved.json")

	if _, err := os.Stat(solvedPath); os.IsNotExist(err) {
		fmt.Println(errorStyle.Render("no solved problems logged yet"))
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
		fmt.Println(errorStyle.Render("no solved problems logged yet"))
		return nil
	}

	var points []DataPoint
	for _, entry := range solved {
		if entry.Timestamp == "" {
			continue
		}

		t, err := time.Parse(time.RFC3339, entry.Timestamp)
		if err != nil {
			continue
		}

		rating := 0
		if entry.Rating != nil {
			rating = *entry.Rating
		}

		points = append(points, DataPoint{
			Date: t.Format("2006-01-02"),
			Rating: rating,
		})
	}

	if len(points) == 0 {
		fmt.Println(errorStyle.Render("no timestamp data available to plot"))
		return nil
	}

	sort.Slice(points, func(i, j int) bool {
		return points[i].Date < points[j].Date
	})

	htmlContent, err := generateGraphHTML(points)
	if err != nil {
		return nil
	}

	htmlPath := filepath.Join(config.GetConfigPath(), "graph.html")
	if err := os.WriteFile(htmlPath, []byte(htmlContent), 0644); err != nil {
		return err
	}

	fmt.Println(greenStyle.Render(fmt.Sprintf("graph generated: %s", htmlPath)))
	fmt.Println(infoStyle.Render("opening in browser..."))

	return browser.OpenFile(htmlPath)
}

func generateGraphHTML(points []DataPoint) (string, error) {
	const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>CPCli - Solved Problems Graph</title>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            flex-direction: column;
            align-items: center;
            padding: 2rem;
        }
        
        .container {
            background: white;
            border-radius: 20px;
            box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
            padding: 2rem;
            max-width: 1200px;
            width: 100%;
        }
        
        h1 {
            color: #667eea;
            margin-bottom: 1rem;
            font-size: 2rem;
            text-align: center;
        }
        
        .stats {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 1rem;
            margin-bottom: 2rem;
        }
        
        .stat-card {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 1.5rem;
            border-radius: 10px;
            text-align: center;
        }
        
        .stat-card h3 {
            font-size: 0.9rem;
            opacity: 0.9;
            margin-bottom: 0.5rem;
        }
        
        .stat-card p {
            font-size: 1.8rem;
            font-weight: bold;
        }
        
        .chart-container {
            position: relative;
            height: 400px;
            margin-top: 2rem;
        }
        
        @media (max-width: 768px) {
            .container {
                padding: 1rem;
            }
            
            h1 {
                font-size: 1.5rem;
            }
            
            .chart-container {
                height: 300px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Kruskal Solved Problems Statistics</h1>
        
        <div class="stats">
            <div class="stat-card">
                <h3>Total Problems</h3>
                <p>{{.TotalProblems}}</p>
            </div>
            <div class="stat-card">
                <h3>Average Rating</h3>
                <p>{{.AverageRating}}</p>
            </div>
            <div class="stat-card">
                <h3>Highest Rating</h3>
                <p>{{.HighestRating}}</p>
            </div>
            <div class="stat-card">
                <h3>Lowest Rating</h3>
                <p>{{.LowestRating}}</p>
            </div>
        </div>
        
        <div class="chart-container">
            <canvas id="problemsChart"></canvas>
        </div>
    </div>

    <script>
        const dates = {{.Dates}};
        const ratings = {{.Ratings}};

        const ctx = document.getElementById('problemsChart').getContext('2d');
        const gradient = ctx.createLinearGradient(0, 0, 0, 400);
        gradient.addColorStop(0, 'rgba(102, 126, 234, 0.8)');
        gradient.addColorStop(1, 'rgba(118, 75, 162, 0.2)');

        new Chart(ctx, {
            type: 'scatter',
            data: {
                datasets: [{
                    label: 'Problem Rating',
                    data: dates.map((date, i) => ({ x: date, y: ratings[i] })),
                    backgroundColor: gradient,
                    borderColor: '#667eea',
                    borderWidth: 2,
                    pointRadius: 6,
                    pointHoverRadius: 8,
                    pointBackgroundColor: '#667eea',
                    pointBorderColor: '#fff',
                    pointBorderWidth: 2,
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        display: true,
                        position: 'top',
                    },
                    tooltip: {
                        callbacks: {
                            label: function(context) {
                                return 'Rating: ' + context.parsed.y + ' on ' + context.parsed.x;
                            }
                        }
                    }
                },
                scales: {
                    x: {
                        type: 'time',
                        time: {
                            unit: 'day',
                            displayFormats: {
                                day: 'MMM DD'
                            }
                        },
                        title: {
                            display: true,
                            text: 'Date',
                            font: {
                                size: 14,
                                weight: 'bold'
                            }
                        },
                        grid: {
                            color: 'rgba(0, 0, 0, 0.05)'
                        }
                    },
                    y: {
                        title: {
                            display: true,
                            text: 'Problem Rating',
                            font: {
                                size: 14,
                                weight: 'bold'
                            }
                        },
                        grid: {
                            color: 'rgba(0, 0, 0, 0.05)'
                        }
                    }
                }
            }
        });
    </script>
    <script src="https://cdn.jsdelivr.net/npm/chartjs-adapter-date-fns/dist/chartjs-adapter-date-fns.bundle.min.js"></script>
</body>
</html>`
	
	totalProblems := len(points)
	totalRating := 0
	minRating := points[0].Rating
	maxRating := points[0].Rating

	for _, p := range points {
		totalRating += p.Rating
		if p.Rating < minRating {
			minRating = p.Rating
		}

		if p.Rating > maxRating {
			maxRating = p.Rating
		}
	}

	avgRating := 0
	if totalProblems > 0 {
		avgRating = totalRating / totalProblems
	}

	
	dates := make([]string, len(points))
	ratings := make([]int, len(points))
	for i, p := range points {
		dates[i] = p.Date
		ratings[i] = p.Rating
	}

	datesJSON, _ := json.Marshal(dates)
	ratingsJSON, _ := json.Marshal(ratings)

	data := struct {
		TotalProblems int
		AverageRating int
		HighestRating int
		LowestRating int
		Dates template.JS
		Ratings template.JS
	}{
		TotalProblems: totalProblems,
		AverageRating: avgRating,
		HighestRating: maxRating,
		LowestRating: minRating,
		Dates: template.JS(datesJSON),
		Ratings: template.JS(ratingsJSON),
	}

	tmpl, err := template.New("graph").Parse(htmlTemplate)
	if err != nil {
		return "", err
	}

	var buf []byte
	writer := &byteWriter{buf: &buf}
	if err := tmpl.Execute(writer, data); err != nil {
		return "", err
	}

	return string(buf), nil
}

type byteWriter struct {
	buf *[]byte
}

func (w *byteWriter) Write(p []byte) (n int, err error) {
	*w.buf = append(*w.buf, p...)
	return len(p), nil
}