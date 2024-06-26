# <p align="center">Welcome to</p>

<p align="center">
    <img src="K-9.svg" alt="K-9" title="K-9 logo" xmlns="http://www.w3.org/2000/svg">
</p>

Simple file watcher inspired by Taskfile based on YAML configuration.

Watches files for changes and executes user defined commands on detection.

## Usage

> [!NOTE]
> If you run K-9 in a directory without a config file it will let you create one automatically.

Since K-9 is inspired by Taskfile the config is a YAML file with a similar structure. The k-9.[yml/yaml] file looks like this:

```yaml
# K-9 config

delay: 10 # seconds

watchers:
- file:
    name: main.go # will watch only main.go
    cmds:
    - cmd /c echo Change in main.go! # has to be an executable in system path
- file:
    name: main # will watch all files named main
    cmds:
    - task build
- file:
    name: .go # will watch all files with the extension .go
    cmds:
    - go test
- file:
    name: ./test # using a directory instead of a file with `exact: true` will watch every file in the directory
    exact: true # defaults to false, if true let's you watch any path specified in name (will error if the path does not exist)
    cmds:
    - cmd /c echo detection in ./test! 
```

> [!IMPORTANT] 
> the commands need to be either an exeutable or passed into a shell of your choosing like `cmd /c`, `pwsh -command`, `bash -c` or `sh -c`.

You can add as many files to watch as you want, each one can execute the commands defined in `cmds`.

The `delay` setting limits how often the commands can be executed within a specified interval in seconds (the delay is separate for each watcher).

To watch all files in a directory use `exact: true`

## Instalation

For instructions on how to build it yourself look in [build](bin/build.md).

Prebuild libraries for:
- Windows
- Linux
- macOS

are availible in [releases](https://github.com/kociumba/K-9/releases).

These should be added to path for command line usage.
