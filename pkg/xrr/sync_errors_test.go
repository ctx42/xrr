package xrr

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_NewSyncErrors(t *testing.T) {
	// --- Given ---
	err0 := errors.New("err0")
	err1 := errors.New("err1")

	col := NewSyncErrors()

	// --- When ---
	col.Add(err0)
	col.Add(err1)

	// --- Then ---
	errs := col.Collect()
	assert.Same(t, err0, errs[0])
	assert.Same(t, err1, errs[1])
}

func Test_SyncErrors(t *testing.T) {
	t.Run("race", func(t *testing.T) {
		// --- Given ---
		col := NewSyncErrors()

		// All goroutines.
		var wgAll sync.WaitGroup
		wgAll.Add(2)

		// Start a goroutine for adding errors.
		start := make(chan struct{})
		go func() {
			<-start
			for i := 0; i < 150; i++ {
				col.Add(fmt.Errorf("err %d", i))
			}
			wgAll.Done()
		}()

		// Start another for adding errors.
		go func() {
			<-start
			for i := 0; i < 150; i++ {
				col.Add(fmt.Errorf("err %d", i))
			}
			wgAll.Done()
		}()

		// --- When ---
		close(start) // Start all goroutines.
		wgAll.Wait() // Wait for all to finish.

		// --- Then ---
		assert.Equal(t, 300, len(col.Collect()))
	})
}

func Test_SyncErrors_Add(t *testing.T) {
	t.Run("add when nil", func(t *testing.T) {
		// --- Given ---
		var se *SyncErrors

		// --- When ---
		se.Add(errors.New("msg"))

		// --- Then ---
		assert.Nil(t, se)
	})

	t.Run("nil errors not added", func(t *testing.T) {
		// --- Given ---
		se := NewSyncErrors()

		// --- When ---
		se.Add(nil, errors.New("msg"), nil)

		// --- Then ---
		assert.Len(t, 1, se.Collect())
	})

	t.Run("add joined errors", func(t *testing.T) {
		// --- Given ---
		je := errors.Join(errors.New("e0"), errors.New("e1"))
		e := errors.New("e3")
		se := NewSyncErrors()

		// --- When ---
		se.Add(je, e)

		// --- Then ---
		ers := se.Collect()
		assert.Len(t, 2, ers)
		assert.Same(t, je, ers[0])
		assert.Same(t, e, ers[1])
	})
}

func Test_SyncErrors_Collect(t *testing.T) {
	t.Run("collect when nil", func(t *testing.T) {
		// --- Given ---
		var se *SyncErrors

		// --- When ---
		ers := se.Collect()

		// --- Then ---
		assert.Nil(t, ers)
	})
}

func Test_SyncErrors_Reset(t *testing.T) {
	t.Run("nil instance", func(t *testing.T) {
		// --- Given ---
		var se *SyncErrors

		// --- When ---
		se.Reset()
	})

	t.Run("success", func(t *testing.T) {
		// --- Given ---
		se := NewSyncErrors()
		se.Add(errors.New("msg"))

		// --- When ---
		se.Reset()

		// --- Then ---
		assert.Empty(t, se.ers)
	})
}
