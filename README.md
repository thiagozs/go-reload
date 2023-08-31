# Go-Reload Runner (beta)

## Overview

`Go-Reload Runner` is a lightweight utility that watches for file changes in a specified directory and automatically runs a given command when a change is detected. It's particularly useful for developers who want to automate tasks like rebuilding or testing their code upon changes.

## Features

- Watches a directory for file changes
- Excludes specified directories from being watched
- Runs a specified command when a change is detected

## Installation

Clone the repository and navigate to the project directory:

```bash
git clone https://github.com/thiagozs/go-reload.git
cd go-reload
```

Build the project:

```bash
go build -o reload-runner
```

## Usage

### Basic Usage

To monitor the current directory and run `echo 'hello'` when a change is detected:

```bash
./reload-runner -dir . -cmd "echo 'hello'"
```

### Excluding Directories

To exclude certain directories from being watched:

```bash
./reload-runner -dir . -cmd "echo 'hello'" -exclude "test,logs"
```

### Example

Let's say you have a Go project and you want to rebuild it whenever a `.go` file changes. You can use `Go-Reload Runner` as follows:

```bash
./reload-runner -dir . -cmd "go build -o my_app main.go"
```

When you save a `.go` file in the directory, `Go-Reload Runner` will automatically run `go build -o my_app main.go`, rebuilding your application.

You must need a folder called `build` inside in your development folder, `Go-Reload Runner` going detect automatic this command and going make a binary file to execute

## Contributing

Feel free to open issues or submit pull requests. Your contributions are welcome!

-----

## Versioning and license

Our version numbers follow the [semantic versioning specification](http://semver.org/). You can see the available versions by checking the [tags on this repository](https://github.com/thiagozs/go-reload/tags). For more details about our license model, please take a look at the [LICENSE](LICENSE) file.

**2023**, thiagozs.
