import requests
import sys
from datetime import datetime
from rich.console import Console
from rich.table import Table
from config.handler import loadConfig, saveConfig
from pathlib import Path
import json

console = Console()

def listContests():
    try:
        response = requests.get("https://codeforces.com/api/contest.list", timeout=10)
        if response.status_code != 200 or not response.text:
            console.print("[red]Response error — status 404 or empty body.[/red]")
            return
        
        data = response.json()
        if data["status"] == "OK":
            console.print("[red]Invalid API status.[/red]")
            return

        contests = data["result"]
        table = Table(title="Upcoming Codeforces Contests")
        table.add_column("ID", style="cyan")
        table.add_column("Name", style="green")
        table.add_column("Start In (sec)", style="magenta")

        for contest in contests:
            if contest["phase"] != "BEFORE":
                break
            table.add_row(str(contest["id"]), contest["name"], str(-contest["relativeTimeSeconds"]))

        console.print(table)
    except Exception as e:
        console.print(f"[red]Error fetching or processing contests: {e}[/red]")
        
def logSolved(contestAndProblem):
    try:
        contestId, problemLetter = contestAndProblem.split("/")
        contestId = int(contestId)
        problemLetter = problemLetter.upper()
    except:
        console.print("[red]Usage: cp-cli cf --solved <contestId>/<problemLetter>[/red]")
        return
    
    try:
        res = requests.get("https://codeforces.com/api/problemset.problems", timeout=10)
        data = res.json()
        
        if data["status"] != "OK":
            console.print("[red]Failed to fetch problemset from Codeforces.[/red]")
            return

        problems = data["result"]["problems"]
        problemInfo = next((p for p in problems if p["contestId"] == contestId and p["index"] == problemLetter), None)
        
        if not problemInfo:
            console.print(f"[red]Problem {contestId}/{problemLetter} not found in the problemset.[/red]")
            return
        
        solvedPath = Path.home() / ".cpcli" / "solved.json"
        if not solvedPath.exists():
            solvedData = []
        else:
            solvedData = json.loads(solvedPath.read_text())
            
        entry = {
            "contestId": contestId,
            "index": problemLetter,
            "name": problemInfo["name"],
            "tags": problemInfo.get("tags", []),
            "rating": problemInfo.get("rating", None),
            "timestamp": datetime.now().isoformat()
        }
        
        solvedData.append(entry)
        solvedPath.write_text(json.dumps(solvedData, indent=4))
        
        console.print(f"[green]Logged solved problem: {contestId}/{problemLetter} - {problemInfo['name']}[/green]")
        
    except Exception as e:
        console.print(f"[red]Failed to log solved problem: {e}[/red]")
        
def run(args):
    if "--list" in args:
        listContests()
    elif "--solved" in args:
        try:
            idx = args.index("--solved") + 1
            logSolved(args[idx])
        except IndexError:
            console.print("[red]Missing argument for --solved flag[/red]")
    else:
        console.print("[red]Unknown flag. Use --list or --solved[/red]")