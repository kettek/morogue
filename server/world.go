package main

// world represents an entire game world state that runs in its own goroutine.
type world struct {
	live     bool
	quitChan chan struct{}
}

func (w *world) loop() {
	w.live = true
	for w.live {
		select {
		case <-w.quitChan:
			w.live = false
		}
	}
}

func (w *world) update() error {
	// TODO: Handle playrars.
	return nil
}
