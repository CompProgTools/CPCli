from rich.console import Console
from rich.markdown import Markdown
from rich.panel import Panel
from rich.table import Table
from rich.text import Text
import requests

console = Console()

def run():
    url = "https://leetcode-api-pied.vercel.app/daily"
    response = requests.get(url, timeout=5)
    data = response.json()
    
    title = data["question"]["title"]
    frontendId = data["question"]["questionFrontendId"]
    date = data["date"]
    difficulty = data["question"]["difficulty"]
    url = "https://leetcode.com" + data["link"]
    
    console.print(Panel.fit(
        f"[bold cyan]LeetCode Daily Challenge - {date}[/bold cyan]",
        title="📅 Daily Challenge"
    ))
    
    console.print(f"[bold yellow]{frontendId}. {title}[/bold yellow] ({difficulty})")
    console.print(f"[blue]{url}[/blue]\n")
    
    content = data["question"]["content"]
    md = Markdown(content, code_theme="monokai")
    console.print(md)