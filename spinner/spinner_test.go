package spinner_test

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/dreamsofcode-io/terminal-ui/spinner"
)

func TestSpinnerStart(t *testing.T) {
	testCases := []struct {
		name     string
		duration time.Duration
		expects  string
	}{
		{
			name:     "should write correct values after 2 frames",
			duration: time.Millisecond * 25,
			expects:  "-\b\\\b",
		},
		{
			name:     "should write correct values after 6 frames",
			duration: time.Millisecond * 105,
			expects:  "-\b\\\b|\b/\b-\b\\\b",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf := &bytes.Buffer{}

			s := spinner.New(spinner.Config{
				Writer:    buf,
				FrameRate: time.Millisecond * 20,
			})

			s.Start()
			time.Sleep(tc.duration)
			s.Stop()

			data, err := io.ReadAll(buf)
			assert.NoError(t, err)

			assert.Equal(t, tc.expects, string(data))
		})
	}
}

func TestSpinnerWorksAsync(t *testing.T) {
	s := spinner.New(spinner.Config{
		Writer:    &bytes.Buffer{},
		FrameRate: time.Millisecond * 5,
	})

	done := make(chan struct{})

	go func() {
		s.Start()
		time.Sleep(10 * time.Millisecond)
		s.Stop()
		close(done)
	}()

	select {
	case <-time.After(time.Millisecond * 200):
		assert.Fail(t, "test timed out")
	case <-done:
		// Test passed
	}
}

func TestWaitingAfterStop(t *testing.T) {
	buf := &bytes.Buffer{}

	s := spinner.New(spinner.Config{
		Writer:    buf,
		FrameRate: time.Millisecond * 20,
	})

	s.Start()
	time.Sleep(time.Millisecond * 35)
	s.Stop()
	time.Sleep(time.Millisecond * 10)

	data, err := io.ReadAll(buf)
	assert.NoError(t, err)

	assert.Equal(t, "-\b\\\b", string(data))
}

func TestStop(t *testing.T) {
	t.Run("calling stop on non started spinner should do nothing", func(t *testing.T) {
		s := spinner.New(spinner.Config{
			Writer:    &bytes.Buffer{},
			FrameRate: time.Millisecond * 5,
		})

		s.Stop()
	})
}

func TestStart(t *testing.T) {
	t.Run("calling start on a started spinner should do nothing", func(t *testing.T) {
		buf := &bytes.Buffer{}
		s := spinner.New(spinner.Config{
			Writer:    buf,
			FrameRate: time.Millisecond * 5,
		})

		s.Start()
		s.Start()

		time.Sleep(6 * time.Millisecond)
		s.Stop()

		data, err := io.ReadAll(buf)
		assert.NoError(t, err)

		assert.Equal(t, "-\b\\\b", string(data))
	})

	t.Run("calling start on a stopped spinner should restart", func(t *testing.T) {
		buf := &bytes.Buffer{}
		s := spinner.New(spinner.Config{
			Writer:    buf,
			FrameRate: time.Millisecond * 5,
		})

		s.Start()
		s.Stop()
		s.Start()
		time.Sleep(6 * time.Millisecond)
		s.Stop()

		data, err := io.ReadAll(buf)
		assert.NoError(t, err)

		assert.Equal(t, "-\b-\b\\\b", string(data))
	})
}
