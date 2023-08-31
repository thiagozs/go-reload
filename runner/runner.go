package runner

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type CommandRunner interface {
	RunCommand(cmd string) (*os.Process, error)
}

type CmdRunner struct {
	params  *CmdRunnerParams
	process *os.Process
}

func NewCommandRunner(opts ...CmdOpts) (*CmdRunner, error) {
	params, err := newCmdParams(opts...)
	if err != nil {
		return nil, err
	}
	return &CmdRunner{params: params}, nil
}

func (c *CmdRunner) RunCommand() (*os.Process, error) {
	return c.startProgram()
}

func (c *CmdRunner) WatchForChanges() {
	if c.params == nil || c.params.watcher == nil {
		log.Println("CmdRunner or watcher is not properly initialized")
		return
	}

	err := filepath.Walk(c.params.GetDirToMonitor(),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() && !c.params.GetExcludedDirs()[path] {
				return c.params.watcher.Add(path)
			}
			return nil
		})
	if err != nil {
		log.Printf("Error : %+v\n", err)
		return
	}

	for {
		select {
		case event, ok := <-c.params.watcher.Events():
			if !ok {
				return
			}
			c.handleEvent(event)

		case err, ok := <-c.params.watcher.Errors():
			if !ok {
				return
			}
			log.Println("Error:", err)

		case <-c.params.watcher.Exit():
			return
		}
	}
}

func (c *CmdRunner) handleEvent(event fsnotify.Event) {

	if _, excluded := c.params.GetExcludedDirs()[event.Name]; !excluded &&
		(event.Op&fsnotify.Write == fsnotify.Write ||
			event.Op&fsnotify.Create == fsnotify.Create) {
		log.Println("Modified file:", event.Name)
		if c.process != nil {
			log.Println("Killing the running process...")
			c.process.Kill()
		}
		log.Println("Starting the program...")
		process, err := c.RunCommand()
		if err != nil {
			log.Println("Error:", err)
		}

		c.process = process
	}
}

func (r *CmdRunner) startProgram() (*os.Process, error) {
	args := strings.Fields(r.params.GetCommand())

	if len(args) < 2 {
		return nil, fmt.Errorf("insufficient arguments in command")
	}

	if strings.Contains(r.params.GetCommand(), "go build") {
		goFile := args[len(args)-1]
		binaryName := args[len(args)-2]

		log.Println("Compiling the program...")
		source := r.params.GetDirToMonitor() + "/" + goFile
		target := r.params.GetDirToMonitor() + "/build/" + binaryName

		log.Printf("go build -o %s %s", target, source)
		cmdBuild := exec.Command("go", "build", "-o", target, source)
		if err := cmdBuild.Run(); err != nil {
			return nil, fmt.Errorf("failed to compile program: %v", err)
		}

		log.Println("Running the compiled binary...")
		log.Println(target)
		cmdRun := exec.Command(target)
		if err := cmdRun.Start(); err != nil {
			return nil, fmt.Errorf("failed to start program: %v", err)
		}

		return cmdRun.Process, nil
	} else {
		log.Println("Running the program directly...")
		cmd := exec.Command(args[0], args[1:]...)
		if err := cmd.Start(); err != nil {
			return nil, fmt.Errorf("failed to start program: %v", err)
		}

		return cmd.Process, nil
	}
}
