// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package prompt

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

// ConfirmFunc runs a yes/no confirmation. Replaceable in tests.
var ConfirmFunc = runHuhConfirm

func runHuhConfirm(title string) (bool, error) {
	var confirmed bool
	err := huh.NewConfirm().
		Title(title).
		Affirmative("Yes").
		Negative("No").
		Value(&confirmed).
		Run()
	return confirmed, err
}

// Confirm asks the user a yes/no question.
func Confirm(question string) (bool, error) {
	result, err := ConfirmFunc(question)
	if err != nil {
		return false, fmt.Errorf("confirmation: %w", err)
	}
	return result, nil
}
