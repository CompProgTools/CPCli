# Welcome to the CP-Cli Documentation Site

This is where you can find quick documentation for the commands and how to use them.

## Table of Contents

- [Installation](#installation)
- [Setting Up The Command](#setting-up-the-command)
- [First Time Setup](#first-time-setup)
- [Setting Up Templates](#setting-up-templates)
- [Syncing Accounts](#syncing-accounts)
- [Testing Code](#learning-how-to-test-code)
- [LeetCode Commands](#leetcode-specific-commands)
- [Codeforces Commands](#codeforces-specific-commands)
- [Statistics](#statistics-commands)

## Installation

**If all you want to do is update the codebase, run `python3 main.py update`**

The installation process of [CP-Cli](https://github.com/CompProgTools/CPCli) is quite simple, as of now, its just a project with python files, but sooner of later, I will turn it into `.exe`, `.dmg`/`.app`, and installations for Linux devices.

As of now, here are the installation steps:
- Clone the repository from Github
```
git clone https://github.com/CompProgTools/CPCli
```
- Change directory into the repository location
```
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

# Setting Up The Command

If you dont need to use the `cp-cli` command specifically, you can do so by just using `python3 main.py`, but for those who need to use the command you can do so by running the following command.

```bash
pip install .
```

This command uses the `setup.py` file in order to setup the `cp-cli` command in order for usage.

You can now use the `cp-cli` command followed by any of the subcommand listed.

## First Time Setup

Now that you've got *CP-Cli* installed, it's important to set it up properly in order to have a good experience!

We will now setup the account(s) on [LeetCode](https://leetcode.com) and [Codeforces](https://codeforces.com)

In order to do the first time setup, run:
```
python3 main.py config
```

The `config` subcommand allows you to set your name, preferred language, code editor, Codeforces username, Leetcode username, and the template output folder.

The only confusing one should be the tempalate output folder, so you can skip that and configure everything else to your liking.

Once thats done you can go to `.cpcli/config.json` file and make sure everything is right. 

**REMEMBER: All your config files and everything else will be stored in the `.cpcli` folder.**

Once you've setup your profiles, run `python3 main.py sync`. This command fetches the rating for the account(s), if this command works successfully that means your accounts are good to go!. 

You can move onto the next step.

## Setting Up Templates

In order to setup a template, you must first understand what a template is. A template is basically pre-written code that is used to speed up the process during the contest. [Here](https://github.com/the-tourist/algo/blob/master/template/multithreaded.cpp) is Tourists template for multithreaded programming.

It's important to have a good templates that works for *you*. While using someone elses template is a good start, you write well with your own template.

In order to create and use a template that CP-Cli can recognize, run the command:
```
cp-cli template --make name.ext --alias alias
```
**BUT WAIT**

In order to run this command you must first understand what each part means.

The `template` subcommand is a collection of flags under the location `subcommands/template.py`. It contains 4 main commands:

1. `--make`: This flag allows you to create a template file under the location `.cpcli/templates` following the `--make` flag, you must put the filename of your extension in the format `filename.extension`. This flag must be followed up by the `--alias` flag.

2. `--alias`: This flag allows you to give a template an alias. In the last flag you were taught how to assign a template, in this one, you will be naming it. An alias is just a name you assign to a template for easy access.

Moving onto actually using the template.

In order to use the template during a contest, here is the command:

```
cp-cli template --use alias --filename name.ext
```
This command uses the same concepts as the last, but with a few changes in the flags used.

1. `--use`: The use flag is to be used in order to call a template by its alias.

2. `--filename`: The filename flag is used to create a filename with the format `filename.extension`.

This simple yet powerful command opens up a file at the location defined in the config command at `python3 main.py config`. This is what the `Set Template Output Folder` option was meant for. Setting a location in this command means that when you use the use and filename flag command, a file is created at the location defined under `Set Template Output Folder` with the selected template alias.

You can also use the `--list` flag in order to list all your teamplates, their aliases, and their template file names. Here is the usage:

```
python3 main.py template --list
```

## Syncing Accounts

If you setup your account(s) using the config command, you can fetch their ratings using the command:
```bash
python3 main.py sync
```

This will show your changes in rating (if any).

## Learning How To Test Code

If you would like to quick test your code by using custom testcases (stdin/stdout) you can do so by using the following command:
```bash
python3 main.py test filename.extension
```

The format here is pretty simple to follow, but for the filename, make sure to put the full path or else *CP-Cli* wont reccognize it as a file in the source directory.

Once you enter the command, you will be asked to ask the number of testcases, **which must be an integer**, after which you can enter the input for the testcase as well as the output.

## LeetCode Specific Commands

As of now, the only LeetCode specific command is the `daily` subcommand. Here is how to use it:

```bash
python3 main.py daily
```

This will display todays daily question in a table format.

## Codeforces Specific Commands

Since the Codeforces API is much more diverse when compared to the third-party LeetCode API, you can expect much more commands to show up.

In order to use any Codeforces command, here is the format:
```bash
python3 main.py cf --flag
```

As of now, CP-Cli offers two commands:

- `--list`: This flag allows you to list upcoming contests in a nice table like format.
- `--solved <contestId>/<index>`: This flag allows you to log a problem you've solved. The format is contestId/Index. Every problem on codeforces comes from a contest, so the id of the contest **not the contest number** and the index is which problem it is in a letter format.

# Statistics Commands

In order to track your progress so far, I've made some subcommands to check number of problems solved, common tags, average rating, etc. These stats include general pieces of info that *I* find useful, but if you want to contribute and add onto, feel free to do so.

Here are the statistics commands:

- `stats`: The stats command gives a general overview of your highest streak, number of problems solved, highest rating, average rating, etc. This command should be enough for most users, but you can customize as you want.