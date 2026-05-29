// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package prompt

import (
	"fmt"
	"testing"
)

func TestConfirmYes(t *testing.T) {
	orig := ConfirmFunc
	defer func() { ConfirmFunc = orig }()
	ConfirmFunc = func(_ string) (bool, error) { return true, nil }

	result, err := Confirm("Delete everything?")
	if err != nil {
		t.Fatal(err)
	}
	if !result {
		t.Error("expected true")
	}
}

func TestConfirmNo(t *testing.T) {
	orig := ConfirmFunc
	defer func() { ConfirmFunc = orig }()
	ConfirmFunc = func(_ string) (bool, error) { return false, nil }

	result, err := Confirm("Delete everything?")
	if err != nil {
		t.Fatal(err)
	}
	if result {
		t.Error("expected false")
	}
}

func TestConfirmError(t *testing.T) {
	orig := ConfirmFunc
	defer func() { ConfirmFunc = orig }()
	ConfirmFunc = func(_ string) (bool, error) { return false, fmt.Errorf("interrupted") }

	_, err := Confirm("question")
	if err == nil {
		t.Error("expected error")
	}
}
