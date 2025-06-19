import os
import sys
from pathlib import Path
from config.handler import loadConfig, saveConfig
from rich.console import Console
import subprocess

console = Console()

def run(args):
    config = loadConfig()
    templatesPath = Path.home() / ".cpcli" / "templates"
    templatesPath.mkdir(parents=True, exist_ok=True)

    # Handle --use
    if "--use" in args and "--filename" in args:
        try:
            aliasIndex = args.index("--use") + 1
            fileIndex = args.index("--filename") + 1
            alias = args[aliasIndex]
            newFileName = args[fileIndex]
        except IndexError:
            console.print("[red]Missing value for --use or --filename[/red]")
            return

        outputDir = config.get("template_output_path")
        if outputDir is None:
            console.print("[red]You must set 'template_output_path' in `cp-cli config` first.[/red]")
            return

        if "templates" not in config or alias not in config["templates"]:
            console.print(f"[red]Alias '{alias}' not found in templates.[/red]")
            return

        sourceTemplate = templatesPath / config["templates"][alias]
        destFile = Path(outputDir) / newFileName

        try:
            destFile.write_text(sourceTemplate.read_text())
            console.print(f"[green]Created {destFile} from template '{alias}'[/green]")
        except Exception as e:
            console.print(f"[red]Failed to copy template: {e}[/red]")
            return

        editor = config.get("preferred_editor", "code")
        try:
            subprocess.Popen([editor, str(destFile)])
        except FileNotFoundError:
            console.print(f"[red]Editor '{editor}' not found. Set it using `cp-cli config`[/red]")
        return

    # Handle --make
    elif "--make" in args and "--alias" in args:
        try:
            nameIndex = args.index("--make") + 1
            aliasIndex = args.index("--alias") + 1
            filename = args[nameIndex]
            alias = args[aliasIndex]
        except IndexError:
            console.print("[red]Missing value for --make or --alias[/red]")
            return

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

        preferredEditor = config.get("preferred_editor", "code")
        try:
            subprocess.Popen([preferredEditor, str(templateFilePath)])
        except FileNotFoundError:
            console.print(f"[red]Editor '{preferredEditor}' not found. Set it using `cp-cli config`[/red]")
        return

    # If neither
    else:
        console.print("[red]Usage:\n  cp-cli template --make name.ext --alias alias\n  cp-cli template --use alias --filename name.ext[/red]")