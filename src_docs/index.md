# Welcome to the Kruskal Documentation Site

This is where you can find quick documentation for the commands and how to use them.

## Table of Contents

- [Installation](#installation)
- [First Time Setup](#first-time-setup)
- [Setting Up Templates](#setting-up-templates)
- [Syncing Accounts](#syncing-accounts)
- [Testing Code](#learning-how-to-test-code)
- [LeetCode Commands](#leetcode-specific-commands)
- [Codeforces Commands](#codeforces-specific-commands)
- [Statistics](#statistics-commands)

## Installation

To install Kruskal, go to the [releases](https://github.com/CompProgTools/kruskal/releases) tab on the [Github Repository](https://github.com/CompProgTools/kruskal) and install whichever fits your operating system!

**Note: Kruskal currently has support for Mac Silicon and Windows 64-Bit, there are more supported versions coming soon.**

### Mac OS Installation

Setting up the application for MacOS is quite easy. Once you have the latest release from the [releases tab](https://github.com/CompProgTools/kruskal/releases), make sure you copy the *full path* of the installed file.

Go to the terminal and type:
```
chmod +x <path of the install>
```

Once that is done *and there are no errors*, the app should be compiled for your computer. To check, go to the folder in which the file is located and type `ls`. If the app name appears then you are good to go!

You may now type in 
```
./freeport
```
And the application should open up.

Conga rats! You just installed Kruskal :D

- Test the output, does it match the following?
```bash
Hi! This is Kruskal, a command line tool for competitive programmers

? What would you like to do? 
> View Repository
  Coming soon...
  Exit
```

If it matches the output above, your installation of [kruskal](https://github.com/CompProgTools/kruskal) is good to go!

## First Time Setup

Now that you've got *Kruskal* installed, it's important to set it up properly in order to have a good experience!

We will now setup the account(s) on [LeetCode](https://leetcode.com) and [Codeforces](https://codeforces.com)

In order to do the first time setup, run:
```bash
kruskal config
```

[Note for reviewer, you may use `kashsuks` as the username for both codeforces and leetcode if you do not have an account]

The `config` subcommand allows you to set your name, preferred language, code editor, Codeforces username, Leetcode username, and the template output folder.

The only confusing one should be the tempalate output folder, so you can skip that and configure everything else to your liking.

Once thats done you can go to `.cpcli/config.json` file and make sure everything is right. 

**REMEMBER: All your config files and everything else will be stored in the `.cpcli` folder.**

Once you've setup your profiles, run `kruskal sync`. This command fetches the rating for the account(s), if this command works successfully that means your accounts are good to go!. 

You can move onto the next step.

**KEEP IN MIND THAT THIS ONLY WORKS IF YOU HAVE A RATING ON EITHER SITE**

## Setting Up Templates

In order to setup a template, you must first understand what a template is. A template is basically pre-written code that is used to speed up the process during the contest. [Here](https://github.com/the-tourist/algo/blob/master/template/multithreaded.cpp) is Tourists template for multithreaded programming.

It's important to have a good templates that works for *you*. While using someone elses template is a good start, you write well with your own template.

In order to create and use a template that Kruskal can recognize, run the command:
```bash
kruskal template --make name.ext --alias alias
```
**BUT WAIT**

In order to run this command you must first understand what each part means.

The `template` subcommand is a collection of flags under the location `subcommands/template.go`. It contains 4 main commands:

1. `--make`: This flag allows you to create a template file under the location `.cpcli/templates` following the `--make` flag, you must put the filename of your extension in the format `filename.extension`. This flag must be followed up by the `--alias` flag.

2. `--alias`: This flag allows you to give a template an alias. In the last flag you were taught how to assign a template, in this one, you will be naming it. An alias is just a name you assign to a template for easy access.

Moving onto actually using the template.

In order to use the template during a contest, here is the command:

```bash
kruskal template --use alias --filename name.ext
```
This command uses the same concepts as the last, but with a few changes in the flags used.

1. `--use`: The use flag is to be used in order to call a template by its alias.

2. `--filename`: The filename flag is used to create a filename with the format `filename.extension`.

This simple yet powerful command opens up a file at the location defined in the config command at `kruskal config`. This is what the `Set Template Output Folder` option was meant for. Setting a location in this command means that when you use the use and filename flag command, a file is created at the location defined under `Set Template Output Folder` with the selected template alias.

You can also use the `--list` flag in order to list all your teamplates, their aliases, and their template file names. Here is the usage:

```bash
kruskal template --list
```

## Syncing Accounts

If you setup your account(s) using the config command, you can fetch their ratings using the command:
```bash
kruskal sync
```

This will show your changes in rating (if any).

## Learning How To Test Code

If you would like to quick test your code by using custom testcases (stdin/stdout) you can do so by using the following command:
```bash
kruskal test filename.extension
```

The format here is pretty simple to follow, but for the filename, make sure to put the full path or else *Kruskal* wont reccognize it as a file in the source directory.

Once you enter the command, you will be asked to ask the number of testcases, **which must be an integer**, after which you can enter the input for the testcase as well as the output.

## LeetCode Specific Commands

As of now, the only LeetCode specific command is the `daily` subcommand. Here is how to use it:

```bash
kruskal daily
```

This will display todays daily question in a table format.

## Codeforces Specific Commands

Since the Codeforces API is much more diverse when compared to the third-party LeetCode API, you can expect much more commands to show up.

In order to use any Codeforces command, here is the format:
```bash
kruskal cf --flag
```

As of now, Kruskal offers two commands:

- `--list`: This flag allows you to list upcoming contests in a nice table like format.
- `--solved <contestId>/<index>`: This flag allows you to log a problem you've solved. The format is contestId/Index. Every problem on codeforces comes from a contest, so the id of the contest **not the contest number** and the index is which problem it is in a letter format.

# Statistics Commands

In order to track your progress so far, I've made some subcommands to check number of problems solved, common tags, average rating, etc. These stats include general pieces of info that *I* find useful, but if you want to contribute and add onto, feel free to do so.

Here are the statistics commands:

- `stats`: The stats command gives a general overview of your highest streak, number of problems solved, highest rating, average rating, etc. This command should be enough for most users, but you can customize as you want.