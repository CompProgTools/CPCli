from rich.console import Console
from InquirerPy import inquirer
from config.handler import loadConfig, saveConfig, setAccount, isAllLinked

console = Console()

def linkAccount():
    platform = inquirer.select(
        message="Select platform to link:",
        choices=["Codeforces", "LeetCode", "AtCoder", "Back"],
        pointer=">",
    ).execute()

    if platform == "Back":
        return

    handle = inquirer.text(
        message=f"Enter your {platform} username:"
    ).execute()

    setAccount(platform, handle)
    console.print(f"[green]{platform} account linked successfully as '{handle}'[/green]")

def main():
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

    if choice == "View Repository":
        console.print("[green]Opening repository... (not implemented yet)[/green]")
    elif choice == "Link Account":
        linkAccount()
    elif choice == "Coming soon...":
        console.print("[yellow]Stay tuned![/yellow]")
    elif choice == "Exit":
        console.print("[red]Goodbye![/red]")

if __name__ == "__main__":
    main()