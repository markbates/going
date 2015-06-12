package chalk_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/markbates/going/chalk"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	a := require.New(t)

	tasks, _ := chalk.New(10)

	y := 0
	m := sync.Mutex{}

	for i := 0; i < 10; i++ {
		tasks <- func() error {
			m.Lock()
			y++
			m.Unlock()
			return nil
		}
	}

	// give the works a second to run
	time.Sleep(10 * time.Millisecond)
	a.Equal(10, y)
}

func Test_Errors(t *testing.T) {
	a := require.New(t)

	tasks, errs := chalk.New(10)

	tasks <- func() error {
		return errors.New("boom!")
	}

	err := <-errs
	a.Error(err)
}
