import json
from pathlib import Path
import matplotlib.pyplot as plt
import matplotlib.dates as mdates
from datetime import datetime
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
    
    # collect data
    dates = []
    ratings = []
    
    for entry in solved:
        timestamp = entry.get("timestamp")
        rating = entry.get("rating", 0) or 0
        if timestamp:
            try:
                dt = datetime.fromisoformat(timestamp)
                dates.append(dt)
                ratings.append(rating)
            except Exception:
                continue
            
    if not dates:
        console.print("[red]No timestamp data available to plot.[/red]")
        return
    
    sortedPoints = sorted(zip(dates, ratings))
    sortedDates, sortedRatings = zip(*sortedPoints)
    
    # lets get plotting (YEHAHHHH)
    plt.figure(figsize=(10, 5))
    plt.scatter(sortedDates, sortedRatings, color="blue", label="Problems Solved")
    
    plt.title("CPCli Solved Problems Over Time")
    plt.xlabel("Date")
    plt.ylabel("Problem Rating")
    
    # format the date axis
    
    plt.gca().xaxis.set_major_locator(mdates.AutoDateLocator())
    plt.gca().xaxis.set_major_formatter(mdates.DateFormatter("%Y-%m-%d"))
    plt.xticks(rotation=45)
    
    plt.legend()
    plt.tight_layout()
    plt.show()