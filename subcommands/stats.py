import json
from collections import Counter
from pathlib import Path
from rich.console import Console

console = Console()

def run():
    solvedPath = Path.home() / ".cpcli" / "solved.json"
    if not solvedPath.exists():
        console.print("[red]No solved problems logged yet.[/red]")
        return

    with solvedPath.open() as f:
        solved = json.load(f)

    if not solved:
        console.print("[red]No solved problems logged yet.[/red]")
        return

    totalProblems = len(solved)

    streakMetaPath = Path.home() / ".cpcli" / "streak_meta.json"
    highScore = 0
    if streakMetaPath.exists():
        try:
            meta = json.loads(streakMetaPath.read_text())
            highScore = meta.get("highscore", 0)
        except Exception:
            highScore = 0
            
    ratings = [entry.get("rating", 0) or 0 for entry in solved]
    averageRating = sum(ratings) / totalProblems if totalProblems else 0

    highestProblem = None
    maxRating = -1
    for entry in solved:
        rating = entry.get("rating", 0) or 0
        if rating > maxRating:
            maxRating = rating
            highestProblem = entry

    tagCounter = Counter()
    for entry in solved:
        tags = entry.get("tags", [])
        tagCounter.update(tags)

    topTags = tagCounter.most_common(3)

    console.rule("[bold cyan]CPCli Stats")

    console.print(f"[bold]Highest Streak:[/bold] {highScore}")
    console.print(f"[bold]Total Problems Logged:[/bold] {totalProblems}")
    console.print(f"[bold]Average Rating:[/bold] {averageRating:.2f}")

    console.print(f"[bold]Highest Rating Problem:[/bold]", end=" ")
    if highestProblem:
        problemName = highestProblem.get("problemName", "N/A")
        problemRating = highestProblem.get("rating", "N/A")
        console.print(f"[bold]{problemName}[/bold] (Rating: {problemRating})")
    else:
        console.print("N/A")

    if topTags:
        tagsStr = ", ".join(f"{tag} ({count})" for tag, count in topTags)
        console.print(f"[bold]Top 3 Tags:[/bold] {tagsStr}")
    else:
        console.print("[bold]Top 3 Tags:[/bold] N/A")