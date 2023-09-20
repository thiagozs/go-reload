package runner

import (
	"os"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
)

type MockWatcher struct {
	EventsChan chan fsnotify.Event
	ErrorsChan chan error
	ExitChan   chan bool
}

func (m *MockWatcher) Add(name string) error {
	return nil
}

func (m *MockWatcher) Remove(name string) error {
	return nil
}

func (m *MockWatcher) Close() error {
	return nil
}

func (m *MockWatcher) Events() chan fsnotify.Event {
	return m.EventsChan
}

func (m *MockWatcher) Errors() chan error {
	return m.ErrorsChan
}

func (m *MockWatcher) Exit() chan bool {
	return m.ExitChan
}

func TestCmdRunner_RunCommand(t *testing.T) {
	mockWatcher := &MockWatcher{
		EventsChan: make(chan fsnotify.Event),
		ErrorsChan: make(chan error),
		ExitChan:   make(chan bool),
	}

	cmdRunner, err := NewCommandRunner(
		DirToMonitor("."),
		Command("echo 'hello'"),
		ExcludedDirs([]string{"test"}),
		RegisterWatcher(mockWatcher),
	)
	assert.NoError(t, err)

	// Test RunCommand
	process, err := cmdRunner.RunCommand()
	assert.NoError(t, err)
	assert.IsType(t, &os.Process{}, process)
}

func TestCmdRunner_WatchForChanges(t *testing.T) {
	mockWatcher := &MockWatcher{
		EventsChan: make(chan fsnotify.Event, 1),
		ErrorsChan: make(chan error, 1),
		ExitChan:   make(chan bool),
	}

	cmdRunner, err := NewCommandRunner(
		DirToMonitor("."),
		Command("echo 'hello'"),
		ExcludedDirs([]string{"test"}),
		RegisterWatcher(mockWatcher),
	)
	assert.NoError(t, err)

	go func() {
		mockWatcher.EventsChan <- fsnotify.Event{Name: "main.go", Op: fsnotify.Write}
	}()

	done := make(chan bool, 1)

	go func() {
		go cmdRunner.WatchForChanges()
		<-time.After(1 * time.Second)
		mockWatcher.ExitChan <- true
		done <- true
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("Test timed out")
	}
}

func TestCmdRunnerParams_GettersAndSetters(t *testing.T) {
	params := &CmdRunnerParams{}
	watcher := &RealWatcher{Watcher: &fsnotify.Watcher{}}
	params.SetDirToMonitor(".")
	params.SetCommand("echo 'hello'")
	params.SetExcludedDirs([]string{"test"})
	params.SetWatcher(watcher)

	assert.Equal(t, ".", params.GetDirToMonitor())
	assert.Equal(t, "echo 'hello'", params.GetCommand())
	assert.Equal(t, []string{"test"}, params.GetExcludedDirs())
	assert.IsType(t, &RealWatcher{}, params.GetWatcher())
}
