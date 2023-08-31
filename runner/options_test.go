package runner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptions_GettersAndSetters(t *testing.T) {
	params := &CmdRunnerParams{}
	mockWatcher := &MockWatcher{}

	// Test setters
	params.SetDirToMonitor(".")
	params.SetCommand("echo 'hello'")
	params.SetExcludedDirs(map[string]bool{"test": true})
	params.SetWatcher(mockWatcher)

	// Test getters
	assert.Equal(t, ".", params.GetDirToMonitor())
	assert.Equal(t, "echo 'hello'", params.GetCommand())
	assert.Equal(t, map[string]bool{"test": true}, params.GetExcludedDirs())
	assert.Equal(t, mockWatcher, params.GetWatcher())
}

func TestCmdOpts(t *testing.T) {
	mockWatcher := &MockWatcher{}

	// Test newCmdParams and CmdOpts
	params, err := newCmdParams(
		DirToMonitor("."),
		Command("echo 'hello'"),
		ExcludedDirs(map[string]bool{"test": true}),
		RegisterWatcher(mockWatcher),
	)

	assert.NoError(t, err)
	assert.Equal(t, ".", params.GetDirToMonitor())
	assert.Equal(t, "echo 'hello'", params.GetCommand())
	assert.Equal(t, map[string]bool{"test": true}, params.GetExcludedDirs())
	assert.Equal(t, mockWatcher, params.GetWatcher())
}
