import os
import sys
from pathlib import Path
from config.handler import loadConfig, saveConfig
from rich.console import Console
import subprocess

console = Console()

def run(args):
    if "--make" not in args or "--alias" not in args:
        console.print("[red]Usage: cp-cli template --make name.extension --alias name[/red]")
        return
    
    try:
        nameIndex = args.index("--make") + 1
        aliasIndex = args.index("--alias") + 1
        filename = args[nameIndex]
        alias = args[aliasIndex]
    except IndexError:
        console.print("[red]Missing value for --make or --alias[/red]")
        return
    
    config = loadConfig()
    templatesPath = Path.home() /".cpcli" /"templates"
    templatesPath.mkdir(parents=True, exist_ok=True)
    
    templateFilePath = templatesPath / filename
    if templateFilePath.exists():
        console.print(f"[yellow]{filename} already exists. Opening it...[/yellow]")
    else:
        templateFilePath.touch()
        console.print(f"[green]Created template: {filename}[/green]")
        
    if "templates" not in config:
        config["templates"] = {}
        
    config["templates"][alias] = filename
    saveConfig(config)
    
    # open in editor
    preferredEditor = config.get("preferred_editor", "code")  # default to VSCode
    try:
        subprocess.Popen([preferredEditor, str(templateFilePath)])
    except FileNotFoundError:
        console.print(f"[red]Editor '{preferredEditor}' not found. Please set a valid preferred editor in config.[/red]")