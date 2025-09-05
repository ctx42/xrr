package xrr

import (
	"testing"
	"time"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_NewMeta(t *testing.T) {
	// --- When ---
	have := Meta()

	// --- Then ---
	assert.Nil(t, have.m)
}

func Test_meta_Str(t *testing.T) {
	t.Run("not existing", func(t *testing.T) {
		// --- Given ---
		m := meta{}

		// --- When ---
		have := m.Str("A", "a")

		// --- Then ---
		assert.Equal(t, map[string]any{"A": "a"}, have.m)
	})

	t.Run("existing", func(t *testing.T) {
		// --- Given ---
		m := meta{m: map[string]any{"A": "a"}}

		// --- When ---
		have := m.Str("A", "b")

		// --- Then ---
		assert.Equal(t, map[string]any{"A": "b"}, have.m)

	})
}

func Test_meta_Int(t *testing.T) {
	t.Run("not existing", func(t *testing.T) {
		// --- Given ---
		m := meta{}

		// --- When ---
		have := m.Int("A", 1)

		// --- Then ---
		assert.Equal(t, map[string]any{"A": 1}, have.m)
	})

	t.Run("existing", func(t *testing.T) {
		// --- Given ---
		m := meta{m: map[string]any{"A": 1}}

		// --- Given ---
		have := m.Int("A", 2)

		// --- Then ---
		assert.Equal(t, map[string]any{"A": 2}, have.m)
	})
}

func Test_meta_Int64(t *testing.T) {
	t.Run("not existing", func(t *testing.T) {
		// --- Given ---
		m := meta{}

		// --- When ---
		have := m.Int64("A", 1)

		// --- Then ---
		assert.Equal(t, map[string]any{"A": int64(1)}, have.m)
	})

	t.Run("existing", func(t *testing.T) {
		// --- Given ---
		m := meta{m: map[string]any{"A": 1}}

		// --- Given ---
		have := m.Int64("A", 2)

		// --- Then ---
		assert.Equal(t, map[string]any{"A": int64(2)}, have.m)
	})
}

func Test_meta_Float64(t *testing.T) {
	t.Run("not existing", func(t *testing.T) {
		// --- Given ---
		m := meta{}

		// --- When ---
		have := m.Float64("A", 1.0)

		// --- Then ---
		assert.Equal(t, map[string]any{"A": 1.0}, have.m)
	})

	t.Run("existing", func(t *testing.T) {
		// --- Given ---
		m := meta{m: map[string]any{"A": 1.0}}

		// --- Given ---
		have := m.Float64("A", 2.0)

		// --- Then ---
		assert.Equal(t, map[string]any{"A": 2.0}, have.m)
	})
}

func Test_meta_Bool(t *testing.T) {
	t.Run("not existing", func(t *testing.T) {
		// --- Given ---
		m := meta{}

		// --- When ---
		have := m.Bool("A", true)

		// --- Then ---
		assert.Equal(t, map[string]any{"A": true}, have.m)
	})

	t.Run("existing", func(t *testing.T) {
		// --- Given ---
		m := meta{m: map[string]any{"A": false}}

		// --- Given ---
		have := m.Bool("A", true)

		// --- Then ---
		assert.Equal(t, map[string]any{"A": true}, have.m)
	})
}

func Test_meta_Time(t *testing.T) {
	t.Run("not existing", func(t *testing.T) {
		// --- Given ---
		tim := time.Now()
		m := meta{}

		// --- When ---
		have := m.Time("A", tim)

		// --- Then ---
		assert.Equal(t, map[string]any{"A": tim}, have.m)
	})

	t.Run("existing", func(t *testing.T) {
		// --- Given ---
		tim0 := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		tim1 := time.Date(2001, 1, 2, 3, 4, 5, 0, time.UTC)
		m := meta{m: map[string]any{"A": tim0}}

		// --- When ---
		have := m.Time("A", tim1)

		// --- Then ---
		assert.Equal(t, map[string]any{"A": tim1}, have.m)
	})
}

func Test_meta_set(t *testing.T) {
	t.Run("nil map", func(t *testing.T) {
		// --- Given ---
		m := meta{}

		// --- When ---
		have := m.set("A", 1)

		// --- Then ---
		assert.Equal(t, map[string]any{"A": 1}, have.m)
	})

	t.Run("existing map", func(t *testing.T) {
		// --- Given ---
		m := meta{m: map[string]any{"A": 1}}

		// --- When ---
		have := m.set("B", 2)

		// --- Then ---
		assert.Equal(t, map[string]any{"A": 1, "B": 2}, have.m)
	})
}
