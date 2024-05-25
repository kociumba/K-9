# K-9
Simple file watcher inspired by Taskfile

Watches files for changes and executes user defined commands on detection.

## Usage

> If you run K-9 in a directory without a config file it will let you create one automatically.

Since K-9 is inspired by Taskfile the config is a yml file with a very similar structure as Taskfile.yml.

```yaml
# K-9 config

delay: 10 # seconds

watchers:
- file:
    name: main.go # name of the file to watch
    cmds:
    - cmd /c echo Hello from K-9! # has to be an exeutable
```

> **NOTE** the commands need to be either an exeutable or passed into a shell of your choosing like `cmd /c` or `pwsh -command`.

You can add as many files to watch as you want, each one can execute the commands defined in cmds.

The delay setting limits how often the commands can be executed within a specified interval (in seconds).

Currently, K-9 can only watch files in the same directory as the k-9.yml configuration file.

For instructions on how to build it look in [build](bin/build.md)