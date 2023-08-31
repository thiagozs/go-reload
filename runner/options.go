package runner

type CmdRunnerParams struct {
	dirToMonitor string
	command      string
	excludedDirs map[string]bool
	watcher      Watcher
}

type CmdOpts func(o *CmdRunnerParams) error

func newCmdParams(opts ...CmdOpts) (*CmdRunnerParams, error) {
	params := &CmdRunnerParams{}
	for _, opt := range opts {
		if err := opt(params); err != nil {
			return nil, err
		}
	}
	return params, nil
}

func DirToMonitor(dir string) CmdOpts {
	return func(o *CmdRunnerParams) error {
		o.dirToMonitor = dir
		return nil
	}
}

func Command(cmd string) CmdOpts {
	return func(o *CmdRunnerParams) error {
		o.command = cmd
		return nil
	}
}

func ExcludedDirs(dirs map[string]bool) CmdOpts {
	return func(o *CmdRunnerParams) error {
		o.excludedDirs = dirs
		return nil
	}
}

func RegisterWatcher(watcher Watcher) CmdOpts {
	return func(o *CmdRunnerParams) error {
		o.watcher = watcher
		return nil
	}
}

// getters and setters

func (r CmdRunnerParams) GetDirToMonitor() string {
	return r.dirToMonitor
}

func (r CmdRunnerParams) GetCommand() string {
	return r.command
}

func (r CmdRunnerParams) GetExcludedDirs() map[string]bool {
	return r.excludedDirs
}

func (r CmdRunnerParams) GetWatcher() Watcher {
	return r.watcher
}

func (r *CmdRunnerParams) SetDirToMonitor(dir string) {
	r.dirToMonitor = dir
}

func (r *CmdRunnerParams) SetCommand(cmd string) {
	r.command = cmd
}

func (r *CmdRunnerParams) SetExcludedDirs(dirs map[string]bool) {
	r.excludedDirs = dirs
}

func (r *CmdRunnerParams) SetWatcher(watcher Watcher) {
	r.watcher = watcher
}
