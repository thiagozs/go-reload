package runner

import "github.com/fsnotify/fsnotify"

type Watcher interface {
	Add(name string) error
	Remove(name string) error
	Close() error
	Events() chan fsnotify.Event
	Errors() chan error
	Exit() chan bool
}

type RealWatcher struct {
	Watcher  *fsnotify.Watcher
	ExitChan chan bool
}

func NewRealWatcher() (*RealWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &RealWatcher{
		Watcher:  watcher,
		ExitChan: make(chan bool),
	}, nil
}

func (r *RealWatcher) Events() chan fsnotify.Event {
	return r.Watcher.Events
}

func (r *RealWatcher) Errors() chan error {
	return r.Watcher.Errors
}

func (r *RealWatcher) Add(name string) error {
	return r.Watcher.Add(name)
}

func (r *RealWatcher) Remove(name string) error {
	return r.Watcher.Remove(name)
}

func (r *RealWatcher) Close() error {
	return r.Watcher.Close()
}

func (r *RealWatcher) Exit() chan bool {
	return r.ExitChan
}
