package chalk

import "sync"

// Handler is the definition of the function that
// will be sent to the pool of go routines
type Handler func() error

type Chalk struct {
	*sync.WaitGroup
	Tasks  chan Handler
	Errors chan error
}

// New creates a new pool N go routines. It will
// return two channels. The first channel is used
// to send `Handler` functions to the pool to be
// executed. The second channel can be used to
// listen for any errors that may occur during the
// processing of a `Handler`.
func New(size int) *Chalk {
	c := Chalk{
		WaitGroup: &sync.WaitGroup{},
		Tasks:     make(chan Handler),
		Errors:    make(chan error),
	}
	go func() {
		for i := 0; i < size; i++ {
			c.Add(1)
			go func() {
				for f := range c.Tasks {
					err := f()
					if err != nil {
						c.Errors <- err
					}
				}
				c.Done()
			}()
		}
	}()
	return &c
}
