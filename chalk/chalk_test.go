package chalk_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/markbates/going/chalk"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	a := require.New(t)

	c := chalk.New(10)

	y := 0
	m := sync.Mutex{}

	for i := 0; i < 10; i++ {
		c.Tasks <- func() error {
			m.Lock()
			y++
			m.Unlock()
			return nil
		}
	}

	c.Wait()
	a.Equal(10, y)
}

func Test_Errors(t *testing.T) {
	a := require.New(t)

	c := chalk.New(10)

	c.Tasks <- func() error {
		return errors.New("boom!")
	}

	err := <-c.Errors
	a.Error(err)
}
