from config.handler import loadConfig, saveConfig
from rich.console import Console
import requests

console = Console()

def fetchRating(platform, handle):
    platform = platform.lower()
    if platform == "leetcode":
        url = f"https://leetcode-api-pied.vercel.app/user/{handle}/contests"
    elif platform == "codeforces":
        url = f"https://codeforces.com/api/user.rating?handle={handle}"
    else:
        return None
    response = requests.get(url, timeout=5)
    data = response.json()
    return extractRating(platform, data)

def extractRating(platform, data):
    try:
        if platform == "leetcode":
            userRank = data.get("userContestRanking")
            if userRank and "rating" in userRank:
                return int(userRank["rating"])
        elif platform == "codeforces":
            return int(data["result"][-1]["newRating"])
    except Exception:
        return None

def run():
    config = loadConfig()
    updated = False

    for platform, handle in list(config.items()):
        if platform.endswith("_rating"):
            continue

        currentRatingKey = f"{platform.lower()}_rating"
        currentRating = config.get(currentRatingKey)
        newRating = fetchRating(platform, handle)

        if newRating is None:
            continue

        if currentRating is None:
            config[currentRatingKey] = newRating
            console.print(f"[green]{platform} ({handle}) rating set to {newRating}[/green]")
            updated = True
        elif newRating != currentRating:
            diff = newRating - currentRating
            changeText = (
                f"increased by {diff}" if diff > 0 else f"decreased by {-diff}"
            )
            console.print(
                f"[cyan]{platform} ({handle}) rating {changeText} to {newRating}[/cyan]"
            )
            config[currentRatingKey] = newRating
            updated = True
        else:
            console.print(
                f"[yellow]{platform} ({handle}) rating unchanged ({currentRating})[/yellow]"
            )

    if updated:
        saveConfig(config)