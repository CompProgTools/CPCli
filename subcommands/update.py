import subprocess
from rich.console import Console
from pathlib import Path

console = Console()

def run(args):
    root = Path(__file__).resolve().parents[1]
    if not (root / ".git").exists():
        console.print("[yellow]This isn't a git repository. Manual update needed.[/yellow]")
        return
    
    try:
        result = subprocess.run(["git", "pull"], cwd=str(root), capture_output=True, text=True)
        if result.returncode == 0:
            console.print("[green]CP-Cli successfully updated![/green]")
            console.print(result.stdout)
        else:
            console.print("[red]Failed to update CPCli.[/red]")
            console.print(result.stderr)
    
    except FileNotFoundError:
        console.print("[red]Git not found. Make sure it's installed and in your PATH.[/red]")