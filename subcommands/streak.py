from rich.console import Console
from rich.table import Table
import json
from datetime import datetime, timedelta
from pathlib import Path

console = Console()

def run():
    solvedPath = Path.home() / ".cpcli" / "solved.json"
    if not solvedPath.exists():
        console.print("[yellow]You haven't solved any problems yet.[/yellow]")
        return
    
    solvedData = json.loads(solvedPath.read_text())
    if not solvedData:
        console.print("[yellow]You haven't solved any problems yet.[/yellow]")
        return
    
    # set of unique dates
    solvedDates = set()
    solvedByDate = {}
    
    for entry in solvedData:
        dateStr = entry["timestamp"].split("T")[0]
        solvedDates.add(dateStr)
        
        if dateStr not in solvedByDate:
            solvedByDate[dateStr] = []
        solvedByDate[dateStr].append(entry)
        
    
    solvedDates = sorted(solvedDates)
    
    today = datetime.now().date()
    yesterday = today - timedelta(days=1)
    
    dateObjs = [datetime.strptime(d, "%Y-%m-%d").date() for d in solvedDates]
    
    streak = 0
    checkDay = today
    
    while checkDay in dateObjs:
        streak += 1
        checkDay = timedelta(days=1)
        
    streakMetaPath = Path.home() / ".cpcli" / "streak.json"
    if streakMetaPath.exists():
        streakMeta = json.loads(streakMetaPath.read_text())
        highScore = streakMeta.get("highScore", 0)
    else:
        highScore = 0
        
    if streak > highScore:
        highScore = streak
        streakMetaPath.write_text(json.dumps({"highscore": highScore}), indent=4)
        
    console.print(f"[bold cyan]Your current streak:[/bold cyan] {streak} day(s)")
    console.print(f"[bold magenta]Your all-time highest streak:[/bold magenta] {highScore} day(s)")
    
    if streak > 0:
        yestStr = yesterday.isoformat()
        if yestStr in solvedByDate:
            table = Table(title="Problems solved yesterday")
            table.add_column("Contest ID", style="cyan")
            table.add_column("Problem", style="green")
            table.add_column("Tags", style="magenta")
            
            for prob in solvedByDate[yestStr]:
                tags = ", ".join(prob.get("tags", []))
                table.add_row(
                    f"{prob['contestId']}/{prob['index']}",
                    prob["name"],
                    tags
                )
                
            console.print(table)
        else:
            console.print("[yellow]You didn't solve any problems yesterday. Streak would reset if no problems solved today.[/yellow]")
    else:
        console.print("[yellow]No current streak. Get solving to start one![/yellow]")