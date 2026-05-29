// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package docker

import (
	"fmt"
	"os/exec"
	"strings"
)

// lookPath is overridable in tests to simulate available scanners.
var lookPath = exec.LookPath

// ScanResult holds the output of a vulnerability scan.
type ScanResult struct {
	Tool     string
	Summary  string
	HasIssue bool
}

// ScanImage runs a vulnerability scan using the first available tool.
func ScanImage(image, tag string) (*ScanResult, error) {
	ref := image + ":" + tag

	if tool, err := lookPath("trivy"); err == nil {
		return scanWithTrivy(tool, ref)
	}

	if _, err := lookPath("docker"); err == nil {
		out, err := CommandRunner.Output("", "docker", "scout", "cves", ref, "--only-severity", "critical,high")
		if err == nil {
			summary := strings.TrimSpace(string(out))
			hasIssue := strings.Contains(summary, "critical") || strings.Contains(summary, "high") || strings.Contains(summary, "HIGH") || strings.Contains(summary, "CRITICAL")
			return &ScanResult{
				Tool:     "docker scout",
				Summary:  truncate(summary, 500),
				HasIssue: hasIssue,
			}, nil
		}
	}

	return nil, fmt.Errorf("no vulnerability scanner found - install trivy (https://trivy.dev) or Docker Scout")
}

func scanWithTrivy(tool, ref string) (*ScanResult, error) {
	// trivy exits non-zero when vulnerabilities are found - that's expected, not an error
	out, _ := CommandRunner.Output("", tool, "image", "--severity", "CRITICAL,HIGH", "--no-progress", ref)
	summary := strings.TrimSpace(string(out))
	hasIssue := strings.Contains(summary, "CRITICAL") || strings.Contains(summary, "HIGH")
	return &ScanResult{
		Tool:     "trivy",
		Summary:  truncate(summary, 500),
		HasIssue: hasIssue,
	}, nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
