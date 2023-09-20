package runner

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/thiagozs/go-reload/pkg/utils"
)

type CommandRunner interface {
	RunCommand(cmd string) (*os.Process, error)
}

type CmdRunner struct {
	params             *CmdRunnerParams
	process            *os.Process
	mu                 sync.Mutex
	wg                 sync.WaitGroup
	lastEventTime      time.Time
	terminateChan      chan struct{}
	processJustStarted bool
}

func NewCommandRunner(opts ...CmdOpts) (*CmdRunner, error) {
	params, err := newCmdParams(opts...)
	if err != nil {
		return nil, err
	}
	return &CmdRunner{
		params:        params,
		terminateChan: make(chan struct{}),
	}, nil
}

func (c *CmdRunner) RunCommand() (*os.Process, error) {
	return c.startProgram()
}

func (c *CmdRunner) TerminateProcess() {
	c.wg.Add(1)
	defer c.wg.Done()

	if c.process != nil {
		log.Println("Sending interrupt signal to the process...")

		if strings.Contains(c.params.GetCommand(), "go run") {
			c.terminateGoRunProcesses()
		} else {
			c.terminateOtherProcesses()
		}
	} else {
		log.Println("No process to terminate")
	}
}

func (c *CmdRunner) terminateGoRunProcesses() {
	patterns := utils.GeneratePatterns(c.params.GetCommand())
	log.Println("Patterns:", patterns)

	var wg sync.WaitGroup

	for _, pattern := range patterns {
		wg.Add(1)
		go func(pattern string) {
			defer wg.Done()

			cmd := fmt.Sprintf("ps aux | grep '%s' | grep -v grep | grep -v reload-runner | awk '{print $2}'", pattern)
			log.Println("Executing command:", cmd)
			out, err := exec.Command("sh", "-c", cmd).Output()
			if err != nil {
				log.Printf("Error finding process by pattern: %v", err)
				return
			}

			pids := strings.Split(strings.TrimSpace(string(out)), "\n")
			for _, pidStr := range pids {
				if len(pidStr) > 0 {
					c.killProcessByPID(pidStr)
				}
			}
		}(pattern)
	}

	wg.Wait()
	log.Println("All 'go run' processes have been terminated")
}

func (c *CmdRunner) killProcessByPID(pidStr string) {
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		log.Printf("Error converting pid to integer: %v", err)
		return
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		log.Printf("Error finding process: %v", err)
		return
	}

	if err := process.Kill(); err != nil {
		log.Printf("Error killing process: %v", err)
	} else {
		log.Printf("Successfully killed process with pid: %d", pid)
	}
}

func (c *CmdRunner) terminateOtherProcesses() {
	if err := c.process.Signal(syscall.SIGTERM); err != nil {
		log.Printf("Error sending interrupt signal: %v", err)
	}

	time.Sleep(2 * time.Second)

	process, err := os.FindProcess(c.process.Pid)
	if err == nil && process != nil {
		log.Println("Process is still running, killing it...")
		if err := c.process.Kill(); err != nil {
			log.Printf("Error killing the process: %v", err)
		}
	}

	log.Println("Waiting for the process to exit...")
	if _, err := c.process.Wait(); err != nil {
		log.Printf("Error waiting for process to exit: %v", err)
	}

	log.Println("Process exited successfully")
}

func (c *CmdRunner) WatchForChanges() {
	if c.params == nil || c.params.watcher == nil {
		log.Println("CmdRunner or watcher is not properly initialized")
		return
	}

	err := filepath.Walk(c.params.GetDirToMonitor(),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Printf("Error at path %s: %+v\n", path, err)
				return err
			}

			if info.IsDir() &&
				!utils.Contains(c.params.GetExcludedDirs(), path) {
				_ = c.params.watcher.Add(path)
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
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	// TODO: make this time configurable debounce
	if now.Sub(c.lastEventTime) < time.Second {
		return
	}
	c.lastEventTime = now

	if c.processJustStarted {
		c.processJustStarted = false
		return
	}

	excluded := utils.Contains(c.params.GetExcludedDirs(), event.Name)
	if !excluded && (event.Op&fsnotify.Write == fsnotify.Write ||
		event.Op&fsnotify.Create == fsnotify.Create) {
		log.Println("Modified file:", event.Name)

		if c.process != nil {
			log.Println("Killing the running process...")
			c.TerminateProcess()
			c.wg.Wait()
		}

		// TODO: make this time configurable
		time.Sleep(2 * time.Second)

		log.Println("Starting the program...")
		process, err := c.RunCommand()
		if err != nil {
			log.Println("Error:", err)
		}

		c.process = process
		c.processJustStarted = true
	}
}

func (r *CmdRunner) startProgram() (*os.Process, error) {
	args := strings.Fields(r.params.GetCommand())
	params := strings.Split(r.params.GetCmdParams(), ",")

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
		//TODO: make this configurable
		cmdBuild.Stdout = os.Stdout
		cmdBuild.Stderr = os.Stderr
		if err := cmdBuild.Run(); err != nil {
			return nil, fmt.Errorf("failed to compile program: %v", err)
		}

		log.Println("Running the compiled binary...")
		log.Println(target)
		var cmdRun *exec.Cmd
		if len(params) > 0 {
			log.Println("[Params]", params)
			cmdRun = exec.Command(target, params...)
		} else {
			cmdRun = exec.Command(target)
		}
		// TODO: make this configurable
		cmdRun.Stdout = os.Stdout
		cmdRun.Stderr = os.Stderr

		if err := cmdRun.Start(); err != nil {
			return nil, fmt.Errorf("failed to start program: %v", err)
		}

		return cmdRun.Process, nil
	} else {
		log.Println("Running the program directly...")
		cmd := exec.Command(args[0], args[1:]...)
		// TODO: make this configurable
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			return nil, fmt.Errorf("failed to start program: %v", err)
		}

		return cmd.Process, nil
	}
}
