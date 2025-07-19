from config.handler import loadConfig
from rich.console import Console
import sys

console = Console()

def run(args):
    if "--login" not in args:
        console.print("[red]Usage: cp-cli openKat --login[/red]")
        return
    try:
        from autokattis import OpenKattis # type: ignore
    except ImportError:
        console.print("[red]autokattis is not installed. Run:[/red]")
        console.print("  [green]pip install autokattis[/green]")
        return
    
    config = loadConfig()
    user = config.get("openkattis_username")
    pwd = config.get("openkattis_password")
    
    if not user or not pwd:
        console.print("[red]OpenKattis credentials not set. Run cp-cli config to save them first.[/red]")
        return
    
    try:
        console.print(f"[blue]Logging in as {user}...[/blue]")
        kt = OpenKattis(user, pwd)
        
        # fetch problems to check if api is working
        
        probs = kt.problems(low_detail_mode=True)
        problemNums = len(probs) if probs is not None else 0
        
        console.print(f"[green]Login successful! You have solved {problemNums} problems on OpenKattis.[/green]")
        
    except Exception as e:
        console.print(f"[red]Login failed or error fetching data: {e}[/red]")
        sys.exit(1)