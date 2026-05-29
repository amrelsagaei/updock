// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package security_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/amrelsagaei/updock/internal/project"
	"github.com/amrelsagaei/updock/internal/scaffold"
	"github.com/amrelsagaei/updock/internal/security"
)

const testSecret = "SUPER_SECRET_VALUE_12345"

func scaffoldTestProject(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	projectDir := filepath.Join(dir, "test-project")
	dataDir := filepath.Join(projectDir, "data")
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		t.Fatal(err)
	}

	cfg := &project.Config{
		Image:       "postgres",
		Tag:         "16",
		ProjectName: "test-project",
		Ports:       []project.PortMapping{{Host: 5432, Container: 5432, Protocol: "tcp"}},
		Env: []project.EnvVar{
			{Key: "POSTGRES_PASSWORD", Value: testSecret, Secret: true, Required: true},
			{Key: "POSTGRES_USER", Value: "admin", Secret: false},
			{Key: "API_TOKEN", Value: "token_xyz_789", Secret: true},
		},
	}

	if err := scaffold.WriteAll(projectDir, cfg); err != nil {
		t.Fatalf("WriteAll error: %v", err)
	}
	return projectDir
}

func TestEnvFileHas0600Permissions(t *testing.T) {
	projectDir := scaffoldTestProject(t)
	info, err := os.Stat(filepath.Join(projectDir, ".env"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("SECURITY: .env should be 0600, got %04o", info.Mode().Perm())
	}
}

func TestGitignoreContainsEnv(t *testing.T) {
	projectDir := scaffoldTestProject(t)
	data, err := os.ReadFile(filepath.Join(projectDir, ".gitignore"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), ".env") {
		t.Error("SECURITY: .gitignore must contain .env")
	}
}

func TestComposeNeverContainsSecretValues(t *testing.T) {
	projectDir := scaffoldTestProject(t)
	data, err := os.ReadFile(filepath.Join(projectDir, "docker-compose.yml"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)

	secrets := []string{testSecret, "token_xyz_789"}
	for _, s := range secrets {
		if strings.Contains(content, s) {
			t.Errorf("SECURITY: docker-compose.yml contains secret value %q", s)
		}
	}
}

func TestUpdockJsonNeverContainsSecretValues(t *testing.T) {
	projectDir := scaffoldTestProject(t)
	data, err := os.ReadFile(filepath.Join(projectDir, "updock.json"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)

	secrets := []string{testSecret, "token_xyz_789"}
	for _, s := range secrets {
		if strings.Contains(content, s) {
			t.Errorf("SECURITY: updock.json contains secret value %q", s)
		}
	}
}

func TestEnvFileContainsSecrets(t *testing.T) {
	projectDir := scaffoldTestProject(t)
	data, err := os.ReadFile(filepath.Join(projectDir, ".env"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)

	if !strings.Contains(content, testSecret) {
		t.Error("SECURITY: .env should contain the secret value (that's where secrets live)")
	}
}

func TestMaskNeverLeaksSecrets(t *testing.T) {
	secrets := []string{testSecret, "token_xyz_789", "p@ssw0rd!"}
	for _, s := range secrets {
		masked := security.MaskSecret(s, true)
		if strings.Contains(masked, s) {
			t.Errorf("SECURITY: MaskSecret leaked %q", s)
		}
		if !strings.Contains(masked, "•") {
			t.Errorf("SECURITY: MaskSecret should use mask chars, got %q", masked)
		}
	}
}
