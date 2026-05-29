// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package prompt

import (
	"testing"
)

func TestProposePorts(t *testing.T) {
	orig := PortChecker
	defer func() { PortChecker = orig }()
	PortChecker = func(_ int) bool { return true }

	ports := ProposePorts([]int{80, 443, 8080})

	if len(ports) != 3 {
		t.Fatalf("expected 3 mappings, got %d", len(ports))
	}
	if ports[0].Host != 80 || ports[0].Container != 80 {
		t.Errorf("expected 80:80, got %d:%d", ports[0].Host, ports[0].Container)
	}
}

func TestProposePortsConflict(t *testing.T) {
	orig := PortChecker
	defer func() { PortChecker = orig }()

	callCount := 0
	PortChecker = func(port int) bool {
		callCount++
		return port != 3000
	}

	ports := ProposePorts([]int{3000})

	if len(ports) != 1 {
		t.Fatalf("expected 1 mapping, got %d", len(ports))
	}
	if ports[0].Host == 3000 {
		t.Error("port 3000 should be skipped because it's in use")
	}
	if ports[0].Container != 3000 {
		t.Errorf("container port should remain 3000, got %d", ports[0].Container)
	}
}

func TestProposePortsEmpty(t *testing.T) {
	ports := ProposePorts(nil)
	if len(ports) != 0 {
		t.Errorf("expected 0 mappings, got %d", len(ports))
	}
}

func TestFindFreePort(t *testing.T) {
	orig := PortChecker
	defer func() { PortChecker = orig }()
	PortChecker = func(port int) bool { return port >= 3005 }

	got := FindFreePort(3000)
	if got != 3005 {
		t.Errorf("expected 3005, got %d", got)
	}
}

func TestConfirmPortMappings(t *testing.T) {
	orig := InputFunc
	defer func() { InputFunc = orig }()
	InputFunc = func(_, _, value string) (string, error) {
		return value, nil
	}

	proposed := ProposePorts([]int{80, 443})
	origChk := PortChecker
	defer func() { PortChecker = origChk }()
	PortChecker = func(_ int) bool { return true }

	result, err := ConfirmPortMappings(proposed)
	if err != nil {
		t.Fatalf("ConfirmPortMappings error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
	if result[0].Host != 80 {
		t.Errorf("expected host port 80, got %d", result[0].Host)
	}
}

func TestConfirmPortMappingsOverride(t *testing.T) {
	orig := InputFunc
	defer func() { InputFunc = orig }()
	InputFunc = func(_, _, _ string) (string, error) {
		return "9090", nil
	}
	origChk := PortChecker
	defer func() { PortChecker = origChk }()
	PortChecker = func(_ int) bool { return true }

	proposed := ProposePorts([]int{80})
	result, err := ConfirmPortMappings(proposed)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if result[0].Host != 9090 {
		t.Errorf("expected 9090, got %d", result[0].Host)
	}
}

func TestConfirmPortMappingsInvalidInput(t *testing.T) {
	orig := InputFunc
	defer func() { InputFunc = orig }()
	InputFunc = func(_, _, value string) (string, error) {
		return "not-a-number", nil
	}
	origChk := PortChecker
	defer func() { PortChecker = origChk }()
	PortChecker = func(_ int) bool { return true }

	proposed := ProposePorts([]int{80})
	result, err := ConfirmPortMappings(proposed)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if result[0].Host != 80 {
		t.Errorf("invalid input should fall back to default, got %d", result[0].Host)
	}
}

func TestRunHuhInputType(t *testing.T) {
	fn := runHuhInput
	_ = fn
}

func TestRunHuhSelectType(t *testing.T) {
	fn := runHuhSelect
	_ = fn
}
