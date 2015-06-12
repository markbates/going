package chalk

import "sync"

// Handler is the definition of the function that
// will be sent to the pool of go routines
type Handler func() error

// New creates a new pool N go routines. It will
// return two channels. The first channel is used
// to send `Handler` functions to the pool to be
// executed. The second channel can be used to
// listen for any errors that may occur during the
// processing of a `Handler`.
func New(size int) (chan Handler, chan error) {
	tasks := make(chan Handler)
	errors := make(chan error)
	go func() {
		wg := sync.WaitGroup{}
		for i := 0; i < size; i++ {
			wg.Add(1)
			go func() {
				for f := range tasks {
					err := f()
					if err != nil {
						errors <- err
					}
				}
				wg.Done()
			}()
		}
	}()
	return tasks, errors
}
