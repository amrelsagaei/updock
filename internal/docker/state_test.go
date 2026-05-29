// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package docker

import (
	"fmt"
	"testing"
)

func TestDetectStateRunning(t *testing.T) {
	m := &mockRunner{output: []byte("running\nrunning\n")}
	withMockRunner(m, func() {
		state := DetectState("/tmp/p", "test")
		if state != StateRunning {
			t.Errorf("expected 'running', got %q", state)
		}
	})
}

func TestDetectStateStopped(t *testing.T) {
	m := &mockRunner{output: []byte("exited\n")}
	withMockRunner(m, func() {
		state := DetectState("/tmp/p", "test")
		if state != StateStopped {
			t.Errorf("expected 'stopped', got %q", state)
		}
	})
}

func TestDetectStateCreated(t *testing.T) {
	m := &mockRunner{output: []byte("created\n")}
	withMockRunner(m, func() {
		state := DetectState("/tmp/p", "test")
		if state != StateCreated {
			t.Errorf("expected 'created', got %q", state)
		}
	})
}

func TestDetectStateNotFound(t *testing.T) {
	m := &mockRunner{err: fmt.Errorf("no such project")}
	withMockRunner(m, func() {
		state := DetectState("/tmp/p", "test")
		if state != StateNotFound {
			t.Errorf("expected 'not found', got %q", state)
		}
	})
}

func TestDetectStateEmpty(t *testing.T) {
	m := &mockRunner{output: []byte("")}
	withMockRunner(m, func() {
		state := DetectState("/tmp/p", "test")
		if state != StateNotFound {
			t.Errorf("expected 'not found' for empty output, got %q", state)
		}
	})
}

func TestDetectStateMixed(t *testing.T) {
	m := &mockRunner{output: []byte("running\nexited\n")}
	withMockRunner(m, func() {
		state := DetectState("/tmp/p", "test")
		if state != StateStopped {
			t.Errorf("mixed state should be 'stopped', got %q", state)
		}
	})
}
