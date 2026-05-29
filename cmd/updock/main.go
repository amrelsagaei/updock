// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package main

import (
	"os"

	"github.com/amrelsagaei/updock/internal/cli"
)

func main() {
	if err := cli.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
