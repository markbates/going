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

func (c *Chalk) Wait() {
	close(c.Tasks)
	close(c.Errors)
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

	for i := 0; i < size; i++ {
		c.Add(1)
		go c.worker()
	}
	// yourLinksSlice := make([]string, 50)
	// for i := 0; i < 50; i++ {
	// 	yourLinksSlice[i] = fmt.Sprintf("%d", i+1)
	// }
	//
	// lCh := make(chan string)
	// wg := new(sync.WaitGroup)
	//
	// // Adding routines to workgroup and running then
	// for i := 0; i < 3; i++ {
	// 	wg.Add(1)
	// 	go worker(lCh, wg)
	// }
	//
	// // Processing all links by spreading them to `free` goroutines
	// for _, link := range yourLinksSlice {
	// 	lCh <- link
	// }
	//
	// // Closing channel (waiting in goroutines won't continue any more)
	// close(lCh)
	//
	// // Waiting for all goroutines to finish (otherwise they die as main routine dies)
	// wg.Wait()

	// go func() {
	// 	for i := 0; i < size; i++ {
	// 		c.Add(1)
	// 		go func() {
	// 			for f := range c.Tasks {
	// 				fmt.Printf("f: %s\n", f)
	// 				err := f()
	// 				if err != nil {
	// 					c.Errors <- err
	// 				}
	// 			}
	// 			c.Done()
	// 		}()
	// 	}
	// }()
	return &c
}
