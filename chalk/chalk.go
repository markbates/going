package chalk

import "sync"

// Task is the definition of the function that
// will be sent to the pool of go routines
type Task func() error

// Chalk wraps the `sync.WaitGroup` and the channels
// used to add tasks to the pool and return errors
// for the tasks in the pool.
type Chalk struct {
	*sync.WaitGroup
	// Tasks channel receives `Task` functions and
	// adds them to the pool.
	Tasks chan Task
	// Errors handles communication of errors that may
	// occur when running a `Task` in the pool.
	Errors chan error
}

// Wait will close the `c.Tasks` channel and then
// wait for all of the tasks in the pool to finish.
// This is a *blocking* operation.
func (c *Chalk) Wait() {
	close(c.Tasks)
	c.WaitGroup.Wait()
}

func (c *Chalk) worker() {
	defer c.Done()

	for f := range c.Tasks {
		err := f()
		if err != nil {
			c.Errors <- err
		}
	}
}

// New creates a new pool N go routines. It will
// return two channels. The first channel is used
// to send `Task` functions to the pool to be
// executed. The second channel can be used to
// listen for any errors that may occur during the
// processing of a `Task`.
func New(size int) *Chalk {
	c := Chalk{
		WaitGroup: &sync.WaitGroup{},
		Tasks:     make(chan Task),
		Errors:    make(chan error),
	}

	for i := 0; i < size; i++ {
		c.Add(1)
		go c.worker()
	}
	return &c
}
