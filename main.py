from rich.console import Console
from InquirerPy import inquirer
from config.handler import loadConfig, saveConfig, setAccount, isAllLinked
from subcommands.sync import fetchRating  # import fetchRating to use when linking
import requests
import sys

console = Console()

def linkAccount():
    platform = inquirer.select(
        message="Select platform to link:",
        choices=["Codeforces", "LeetCode", "AtCoder", "Back"],
        pointer=">",
    ).execute()

    if platform == "Back":
        return

    while True:
        handle = inquirer.text(
            message=f"Enter your {platform} username:"
        ).execute()

        if platform == "Codeforces":
            url = f"https://codeforces.com/api/user.rating?handle={handle}"
            try:
                response = requests.get(url, timeout=5)
                data = response.json()
                if data["status"] == "OK":
                    break
                else:
                    console.print("[red]User not found on Codeforces. Try again.[/red]")
            except Exception:
                console.print("[red]Error connecting to Codeforces. Try again.[/red]")

        elif platform == "LeetCode":
            url = f"https://leetcode-api-pied.vercel.app/user/{handle}"
            try:
                response = requests.get(url, timeout=5)
                data = response.json()
                if data.get("detail") == "404: User not found":
                    console.print("[red]User not found on LeetCode. Try again.[/red]")
                else:
                    break
            except Exception:
                console.print("[red]Error connecting to LeetCode API. Try again.[/red]")

        elif platform == "AtCoder":
            break

    setAccount(platform, handle)

    initialRating = fetchRating(platform, handle)
    if initialRating is not None:
        config = loadConfig()
        config[f"{platform.lower()}_rating"] = initialRating
        saveConfig(config)
        console.print(f"[green]{platform} account linked and rating set to {initialRating}[/green]")
    else:
        console.print(f"[yellow]{platform} account linked, but failed to fetch rating.[/yellow]")

def menu():
    options = ["View Repository"]
    if not isAllLinked():
        options.append("Link Account")
    options += ["Coming soon...", "Exit"]

    console.print("[bold blue]Hi! This is CPCli, a command line tool for competitive programmers[/bold blue]\n")

    choice = inquirer.select(
        message="What would you like to do?",
        choices=options,
        pointer=">",
    ).execute()

    import webbrowser

    if choice == "View Repository":
        console.print("[green]Opening CPCli GitHub repository...[/green]")
        webbrowser.open("https://github.com/CompProgTools/CPCli")
    elif choice == "Link Account":
        linkAccount()
    elif choice == "Coming soon...":
        console.print("[yellow]Stay tuned![/yellow]")
    elif choice == "Exit":
        console.print("[red]Goodbye![/red]")

if __name__ == "__main__":
    if len(sys.argv) > 1:
        from subcommands.sync import run as syncRun
        subcommand = sys.argv[1]
        if subcommand == "sync":
            syncRun()
        else:
            console.print(f"[red]Unknown subcommand: {subcommand}[/red]")
    else:
        menu()