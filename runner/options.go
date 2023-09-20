package runner

type CmdRunnerParams struct {
	dirToMonitor string
	command      string
	cmdparams    string
	excludedDirs []string
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

func CmdParams(params string) CmdOpts {
	return func(o *CmdRunnerParams) error {
		o.cmdparams = params
		return nil
	}
}

func ExcludedDirs(dirs []string) CmdOpts {
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

// getters ----

func (r CmdRunnerParams) GetDirToMonitor() string {
	return r.dirToMonitor
}

func (r CmdRunnerParams) GetCommand() string {
	return r.command
}

func (r CmdRunnerParams) GetExcludedDirs() []string {
	return r.excludedDirs
}

func (r CmdRunnerParams) GetWatcher() Watcher {
	return r.watcher
}

func (r *CmdRunnerParams) GetCmdParams() string {
	return r.cmdparams
}

// setters ----

func (r *CmdRunnerParams) SetCmdParams(params string) {
	r.cmdparams = params
}

func (r *CmdRunnerParams) SetDirToMonitor(dir string) {
	r.dirToMonitor = dir
}

func (r *CmdRunnerParams) SetCommand(cmd string) {
	r.command = cmd
}

func (r *CmdRunnerParams) SetExcludedDirs(dirs []string) {
	r.excludedDirs = dirs
}

func (r *CmdRunnerParams) SetWatcher(watcher Watcher) {
	r.watcher = watcher
}
