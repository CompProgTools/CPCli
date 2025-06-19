# Welcome to the CP-Cli Documentation Site

This is where you can find quick documentation for the commands and how to use them.

## Table of Contents

- [Installation](#installation)
- [First Time Setup](#first-time-setup)

## Installation

The installation process of [CP-Cli](https://github.com/CompProgTools/CPCli) is quite simple, as of now, its just a project with python files, but sooner of later, I will turn it into `.exe`, `.dmg`/`.app`, and installations for Linux devices.

As of now, here are the installation steps:
- Clone the repository from Github
```
bash
git clone https://github.com/CompProgTools/CPCli
```
- Change directory into the repository location
```
bash
cd CPCli
```
- Run the python file
```bash
python3 main.py
```

- Test the output, does it match the following?
```bash
Hi! This is CPCli, a command line tool for competitive programmers

? What would you like to do? 
> View Repository
  Coming soon...
  Exit
```

If it matches the output above, your installation of [CP-Cli](https://github.com/CompProgTools/CPCli) is good to go!

## First Time Setup

Now that you've got *CP-Cli* installed, it's important to set it up properly in order to have a good experience!

We will now setup the account(s) on [LeetCode](https://leetcode.com) and [Codeforces](https://codeforces.com)

In order to do the first time setup, run:
```
bash
python3 main.py config
```

The `config` subcommand allows you to set your name, preferred language, code editor, Codeforces username, Leetcode username, and the template output folder.

The only confusing one should be the tempalate output folder, so you can skip that and configure everything else to your liking.

Once thats done you can go to `.cpcli/config.json` file and make sure everything is right. 

**REMEMBER: All your config files and everything else will be stored in the `.cpcli` folder.**

Once you've setup your profiles, run `python3 main.py sync`. This command fetches the rating for the account(s), if this command works successfully that means your accounts are good to go!. 

You can move onto the next step.

