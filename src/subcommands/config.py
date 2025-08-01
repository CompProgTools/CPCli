from config.handler import loadConfig, saveConfig
from InquirerPy import inquirer
from rich.console import Console
import getpass

console = Console()

def run(_):
    config = loadConfig()
    while True:
        choice = inquirer.select(
            message="What would you like to configure?",
            choices=[
                "Set Name",
                "Set Preferred Language",
                "Set Preferred Code Editor",
                "Set Codeforces Username",
                "Set LeetCode Username",
                "Set Template Output Folder",
                "Set OpenKattis Credentials",
                "Back"
            ],
            pointer=">",
        ).execute()
        
        if choice == "Back":
            break
        
        if choice == "Set Name":
            name = inquirer.text(message="Enter your name:").execute()
            config["name"] = name
            console.print(f"[green]Name set to {name}[/green]")
            
        elif choice == "Set Preferred Language":
            lang = inquirer.select(
                message="Choose your preferred language:",
                choices=["C++", "Python", "Java", "C", "Other"],
                pointer=">",
            ).execute()
            
            config["preferred_language"] = lang
            console.print("[green]Preferred language set to {lang}[/green]")
        
        elif choice == "Set OpenKattis Credentials":
            kattisUser = inquirer.text(message="Enter OpenKattis Username: ").execute()
            kattisPass = getpass.getpass("Enter OpenKattis password (hidden): ")
            config["openkattis_username"] = kattisUser
            config["openkattis_password"] = kattisPass
            console.print("[green]OpenKattis credentials saved.[/green]")
            
        elif choice == "Set Preferred Code Editor":
            editorMap = {
                "VSCode": "code",
                "Neovim": "nvim",
                "Vim": "vim",
                "Sublime": "subl",
                "Atom": "atom",
                "Other": None
            }

            editorChoice = inquirer.select(
                message="Choose your preferred code editor:",
                choices=list(editorMap.keys()),
                pointer=">",
            ).execute()

            if editorMap[editorChoice] is None:
                custom = inquirer.text(message="Enter the terminal command for your editor:").execute()
                config["preferred_editor"] = custom
                console.print(f"[green]Preferred editor set to {custom}[/green]")
            else:
                config["preferred_editor"] = editorMap[editorChoice]
                console.print(f"[green]Preferred editor set to {editorMap[editorChoice]}[/green]")
        
        elif choice == "Set Codeforces Username":
            handle = inquirer.text(message="Enter Codeforces Handle: ").execute()
            config["codeforces"] = handle
            console.print("[green]Codeforces username updated to {handle}[/green]")
            
        elif choice == "Set LeetCode Username":
            handle = inquirer.text(message="Enter LeetCode handle:").execute()
            config["leetcode"] = handle
            console.print(f"[green]LeetCode username updated to {handle}[/green]")
        
        elif choice == "Set Template Output Folder":
            path = inquirer.text(message="Enter absolute path for template outputs to start writing code: ").execute()
            config["template_output_path"] = path
            console.print(f"[green]Template output path saved to {path}[/green]")
        
        saveConfig(config)