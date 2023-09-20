package main

import (
	"flag"
	"fmt"
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
	params       string
)

func main() {
	flag.StringVar(&dirToMonitor, "dir", ".", "Directory to monitor for changes")
	flag.StringVar(&command, "cmd", "", "Command to run when a change is detected")
	flag.StringVar(&excluded, "exclude", "", "Comma-separated list of directories to exclude from monitoring")
	flag.StringVar(&params, "params", "", "Comma-separated list of parameters to pass to the command")
	flag.Parse()

	if command == "" {
		log.Fatal("Please specify a command to run using the -cmd flag")
	}

	var dirExcluded []string
	if strings.Contains(excluded, ",") {
		norm := strings.Split(excluded, ",")
		for _, v := range norm {
			dirExcluded = append(dirExcluded,
				fmt.Sprintf("%s/%s", dirToMonitor, v))
		}
	} else {
		dirExcluded = append(dirExcluded,
			fmt.Sprintf("%s/%s", dirToMonitor, excluded))
	}

	rw, err := runner.NewRealWatcher()
	if err != nil {
		log.Fatal(err)
	}

	opts := []runner.CmdOpts{
		runner.DirToMonitor(dirToMonitor),
		runner.Command(command),
		runner.ExcludedDirs(dirExcluded),
		runner.RegisterWatcher(rw),
		runner.CmdParams(params),
	}

	r, err := runner.NewCommandRunner(opts...)
	if err != nil {
		log.Fatal(err)
	}

	go r.WatchForChanges()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	<-done

	rw.Close()

	log.Println("Stopping the watcher...")
}
