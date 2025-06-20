import requests
from rich.console import Console
from rich.table import Table
from datetime import datetime

console = Console()

def run(args):
    if "--list" in args:
        url = "https://codeforces.com/api/contest.list?gym=false"
        try:
            response = requests.get(url, timeout=10)
            if response.status_code != 200 or not response.text.strip():
                console.print(f"[red]Response error — status {response.status_code} or empty body.[/red]")
                return

            data = response.json()
            if data["status"] != "OK":
                console.print(f"[red]Failed to fetch contests.[/red]")
                return

            contests = data["result"]
            upcoming = []

            for contest in contests:
                if contest["phase"] == "BEFORE":
                    upcoming.append(contest)
                else:
                    break

            if not upcoming:
                console.print("[yellow]No contests found[/yellow]")
                return

            table = Table(title="Upcoming Codeforces Contests", show_lines=True)
            table.add_column("ID", style="cyan", justify="right")
            table.add_column("Name", style="bold")
            table.add_column("Start Time", style="green")
            table.add_column("Duration", style="magenta")

            for contest in reversed(upcoming):
                startTime = datetime.fromtimestamp(contest["startTimeSeconds"]).strftime("%Y-%m-%d %H:%M")
                durationHours = contest["durationSeconds"] // 3600
                durationMinutes = (contest["durationSeconds"] % 3600) // 60
                durationString = f"{durationHours}h {durationMinutes}m"
                table.add_row(str(contest["id"]), contest["name"], startTime, durationString)

            console.print(table)

        except Exception as e:
            console.print(f"[red]Error fetching or processing contests: {e}[/red]")

    else:
        console.print("[red]Usage: cp-cli cf --list[/red]")