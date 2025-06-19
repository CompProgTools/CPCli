import json
from pathlib import Path

configPath = Path.home() /".cpcli"
configFile = configPath / "config.json"

def loadConfig():
    if not configPath.exists():
        configPath.mkdir(parents=True)
        
    if not configFile.exists():
        with open(configFile, "w") as f:
            json.dump({}, f)
            
    with open(configFile, "r") as f:
        return json.load(f)
        
def saveConfig(config):
    with open(configFile, "w") as f:
        json.dump(config, f, indent=4)
        
def setAccount(platform: str, handle: str):
    config = loadConfig()
    config[platform.lower()] = handle
    saveConfig(config)
    
def isAllLinked():
    config = loadConfig()
    return all(k in config for k in ["codeforces", "leetcode"])