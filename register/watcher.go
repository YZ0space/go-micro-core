package registry

// Watcher is an interface that returns updates
// about services within the registry.
type Watcher interface {
	// Next is a blocking call
	Next() (*Result, error)
	Stop()
}

// Result is returned by a call to Next on
// the watcher. Actions can be created, update, delete
type Result struct {
	Action  string
	Service *Service
}
