package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/thiagozs/go-reload/runner"
)

var (
	dirToMonitor string
	command      string
	excluded     string
)

func main() {
	flag.StringVar(&dirToMonitor, "dir", ".", "Directory to monitor for changes")
	flag.StringVar(&command, "cmd", "", "Command to run when a change is detected")
	flag.StringVar(&excluded, "exclude", "", "Comma-separated list of directories to exclude from monitoring")
	flag.Parse()

	if command == "" {
		log.Fatal("Please specify a command to run using the -cmd flag")
	}

	excludedDirs := make(map[string]bool)
	for _, dir := range strings.Split(excluded, ",") {
		excludedDirs[dir] = true
	}

	realWatcher, err := runner.NewRealWatcher()
	if err != nil {
		log.Fatal(err)
	}

	opts := []runner.CmdOpts{
		runner.DirToMonitor(dirToMonitor),
		runner.Command(command),
		runner.ExcludedDirs(excludedDirs),
		runner.RegisterWatcher(realWatcher),
	}

	r, err := runner.NewCommandRunner(opts...)
	if err != nil {
		log.Fatal(err)
	}

	go r.WatchForChanges()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	<-done

	realWatcher.Close()

	log.Println("Stopping the watcher...")
}
