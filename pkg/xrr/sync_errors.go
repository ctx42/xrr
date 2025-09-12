// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package xrr

import (
	"slices"
	"sync"
)

// SyncErrors is a thread-safe error slice.
type SyncErrors struct {
	ers []error
	mx  sync.Mutex
}

// NewSyncErrors returns a new instance of [SyncErrors].
func NewSyncErrors() *SyncErrors {
	return &SyncErrors{ers: make([]error, 0, 5)}
}

// Add adds an error to [SyncErrors] in a thread-safe way. The nil errors are
// ignored. If the [SyncErrors] is nil, the call is no-op.
func (ers *SyncErrors) Add(err ...error) {
	if ers == nil {
		return
	}
	ers.mx.Lock()
	defer ers.mx.Unlock()

	for _, e := range err {
		if e == nil {
			continue
		}
		ers.ers = append(ers.ers, e)
	}
}

// Collect retrieves all errors from [SyncErrors] and resets the internal error
// slice. If the [SyncErrors] is nil, the call is no-op and returns nil.
func (ers *SyncErrors) Collect() []error {
	if ers == nil {
		return nil
	}
	ers.mx.Lock()
	defer ers.mx.Unlock()

	clone := slices.Clone(ers.ers)
	ers.ers = ers.ers[:0]
	return clone
}

// Reset clears any errors in the slice.
func (ers *SyncErrors) Reset() {
	if ers == nil {
		return
	}
	ers.mx.Lock()
	defer ers.mx.Unlock()
	ers.ers = ers.ers[:0]
}
