// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package recipe

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/amrelsagaei/updock/internal/prompt"
)

func wpRecipe() *Recipe {
	return &Recipe{
		Meta: Meta{Name: "wordpress", Description: "WP + MySQL"},
		Prompts: []Prompt{
			{Name: "DB_PASSWORD", Required: true, Secret: true, Generate: "password"},
			{Name: "DB_NAME", Default: "wordpress"},
			{Name: "WP_PORT", Default: "8080", Type: "port"},
		},
		Services: map[string]Service{
			"wordpress": {
				Image:      "wordpress",
				DefaultTag: "latest",
				Ports:      []string{"${WP_PORT}:80"},
				Environment: map[string]string{
					"WORDPRESS_DB_PASSWORD": "${DB_PASSWORD}",
					"WORDPRESS_DB_NAME":     "${DB_NAME}",
				},
				DependsOn: []string{"db"},
			},
			"db": {
				Image:      "mysql",
				DefaultTag: "8.0",
				Environment: map[string]string{
					"MYSQL_ROOT_PASSWORD": "${DB_PASSWORD}",
					"MYSQL_DATABASE":      "${DB_NAME}",
				},
				Volumes: []string{"db_data:/var/lib/mysql"},
			},
		},
		Volumes: map[string]struct{}{"db_data": {}},
	}
}

func TestCollectValuesAutoGenerate(t *testing.T) {
	r := wpRecipe()
	values, err := CollectValues(r, true)
	if err != nil {
		t.Fatalf("CollectValues error: %v", err)
	}

	if values["DB_PASSWORD"] == "" {
		t.Error("password should be auto-generated")
	}
	if len(values["DB_PASSWORD"]) < 32 {
		t.Errorf("generated password should be 32+ chars, got %d", len(values["DB_PASSWORD"]))
	}
	if values["DB_NAME"] != "wordpress" {
		t.Errorf("DB_NAME should default to 'wordpress', got %q", values["DB_NAME"])
	}
	if values["WP_PORT"] != "8080" {
		t.Errorf("WP_PORT should default to '8080', got %q", values["WP_PORT"])
	}
}

func TestCollectValuesManualInput(t *testing.T) {
	orig := prompt.InputFunc
	defer func() { prompt.InputFunc = orig }()
	prompt.InputFunc = func(_, _, _ string) (string, error) {
		return "manual-password", nil
	}

	r := &Recipe{
		Meta: Meta{Name: "test"},
		Prompts: []Prompt{
			{Name: "SECRET", Required: true, Secret: true},
		},
		Services: map[string]Service{
			"app": {Image: "nginx"},
		},
	}

	values, err := CollectValues(r, false)
	if err != nil {
		t.Fatal(err)
	}
	if values["SECRET"] != "manual-password" {
		t.Errorf("expected 'manual-password', got %q", values["SECRET"])
	}
}

func TestCollectValuesRequiredEmpty(t *testing.T) {
	orig := prompt.InputFunc
	defer func() { prompt.InputFunc = orig }()
	prompt.InputFunc = func(_, _, _ string) (string, error) {
		return "", nil
	}

	r := &Recipe{
		Meta: Meta{Name: "test"},
		Prompts: []Prompt{
			{Name: "REQUIRED_VAR", Required: true},
		},
		Services: map[string]Service{
			"app": {Image: "nginx"},
		},
	}

	_, err := CollectValues(r, false)
	if err == nil {
		t.Error("expected error for empty required value")
	}
}

func TestRenderWritesFiles(t *testing.T) {
	dir := t.TempDir()
	r := wpRecipe()
	values := map[string]string{
		"DB_PASSWORD": "test-pw",
		"DB_NAME":     "wpdb",
		"WP_PORT":     "9090",
	}

	if err := Render(r, values, dir, r.Meta.Name); err != nil {
		t.Fatalf("Render error: %v", err)
	}

	expectedFiles := []string{"docker-compose.yml", ".env", ".gitignore", "updock.json"}
	for _, f := range expectedFiles {
		if _, err := os.Stat(filepath.Join(dir, f)); os.IsNotExist(err) {
			t.Errorf("expected file %q", f)
		}
	}
}

func TestRenderComposeContent(t *testing.T) {
	dir := t.TempDir()
	r := wpRecipe()
	values := map[string]string{
		"DB_PASSWORD": "secret123",
		"DB_NAME":     "mydb",
		"WP_PORT":     "8080",
	}

	if err := Render(r, values, dir, r.Meta.Name); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "docker-compose.yml"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)

	if !strings.Contains(content, "wordpress:latest") {
		t.Error("compose should contain wordpress:latest")
	}
	if !strings.Contains(content, "mysql:8.0") {
		t.Error("compose should contain mysql:8.0")
	}
	if !strings.Contains(content, "8080:80") {
		t.Error("compose should have expanded port")
	}
	if !strings.Contains(content, "depends_on") {
		t.Error("compose should have depends_on")
	}
	if !strings.Contains(content, "db_data") {
		t.Error("compose should reference volumes")
	}
}

func TestRenderSharedVariable(t *testing.T) {
	dir := t.TempDir()
	r := wpRecipe()
	values := map[string]string{
		"DB_PASSWORD": "shared-pw",
		"DB_NAME":     "shared-db",
		"WP_PORT":     "8080",
	}

	if err := Render(r, values, dir, r.Meta.Name); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "docker-compose.yml"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)

	count := strings.Count(content, "shared-pw")
	if count < 2 {
		t.Errorf("DB_PASSWORD should appear in both services, found %d times", count)
	}

	dbCount := strings.Count(content, "shared-db")
	if dbCount < 2 {
		t.Errorf("DB_NAME should appear in both services, found %d times", dbCount)
	}
}

func TestRenderEnvFilePermissions(t *testing.T) {
	dir := t.TempDir()
	r := wpRecipe()
	values := map[string]string{"DB_PASSWORD": "x", "DB_NAME": "y", "WP_PORT": "8080"}

	if err := Render(r, values, dir, r.Meta.Name); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(filepath.Join(dir, ".env"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("expected .env permissions 0600, got %04o", info.Mode().Perm())
	}
}

func TestRenderEnvContainsValues(t *testing.T) {
	dir := t.TempDir()
	r := wpRecipe()
	values := map[string]string{"DB_PASSWORD": "pw123", "DB_NAME": "mydb", "WP_PORT": "8080"}

	if err := Render(r, values, dir, r.Meta.Name); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(dir, ".env"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)

	if !strings.Contains(content, "DB_PASSWORD=pw123") {
		t.Error(".env should contain DB_PASSWORD")
	}
	if !strings.Contains(content, "DB_NAME=mydb") {
		t.Error(".env should contain DB_NAME")
	}
}

func TestRenderRecordsProjectName(t *testing.T) {
	dir := t.TempDir()
	r := wpRecipe()
	values := map[string]string{"DB_PASSWORD": "x", "DB_NAME": "y", "WP_PORT": "8080"}

	if err := Render(r, values, dir, "wordpress-2"); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "updock.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "wordpress-2") {
		t.Error("updock.json should record the on-disk project name")
	}
}

func TestParsePortMapping(t *testing.T) {
	tests := []struct {
		name          string
		in            string
		wantOK        bool
		wantHost      int
		wantContainer int
	}{
		{"host and container", "8080:80", true, 8080, 80},
		{"with ip", "127.0.0.1:8080:80", true, 8080, 80},
		{"with protocol", "8080:80/tcp", true, 8080, 80},
		{"container only", "80", false, 0, 0},
		{"non-numeric", "${WP_PORT}:80", false, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, ok := parsePortMapping(tt.in)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if ok && (m.Host != tt.wantHost || m.Container != tt.wantContainer) {
				t.Errorf("got %d:%d, want %d:%d", m.Host, m.Container, tt.wantHost, tt.wantContainer)
			}
		})
	}
}

func TestRecipeToConfigRecordsPorts(t *testing.T) {
	r := &Recipe{
		Meta:     Meta{Name: "wp"},
		Services: map[string]Service{"web": {Image: "wordpress", Ports: []string{"${WP_PORT}:80"}}},
	}
	cfg := recipeToConfig(r, map[string]string{"WP_PORT": "8080"})
	if len(cfg.Ports) != 1 {
		t.Fatalf("expected 1 port mapping, got %d", len(cfg.Ports))
	}
	if cfg.Ports[0].Host != 8080 || cfg.Ports[0].Container != 80 {
		t.Errorf("got %d:%d, want 8080:80", cfg.Ports[0].Host, cfg.Ports[0].Container)
	}
}
