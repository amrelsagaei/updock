// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/amrelsagaei/updock/internal/hub"
	"github.com/amrelsagaei/updock/internal/prompt"
	"github.com/amrelsagaei/updock/internal/recipe"
	"github.com/charmbracelet/huh"
)

// fakeSource implements imageSource for tests - no network.
type fakeSource struct {
	searchResults []hub.SearchResult
	searchErr     error
	tags          []hub.Tag
	tagsErr       error
	imgCfg        *hub.ImageConfig
	inspectErr    error
}

func (f *fakeSource) Search(_ string, _ int) ([]hub.SearchResult, error) {
	return f.searchResults, f.searchErr
}

func (f *fakeSource) Tags(_ string, _ int) ([]hub.Tag, error) {
	return f.tags, f.tagsErr
}

func (f *fakeSource) Inspect(_, _ string) (*hub.ImageConfig, error) {
	return f.imgCfg, f.inspectErr
}

// stubPrompts replaces all interactive prompt funcs with deterministic ones.
func stubPrompts(t *testing.T, selectVal, inputVal string) {
	t.Helper()
	origSel, origIn, origPort := prompt.SelectFunc, prompt.InputFunc, prompt.PortChecker
	t.Cleanup(func() {
		prompt.SelectFunc = origSel
		prompt.InputFunc = origIn
		prompt.PortChecker = origPort
	})
	prompt.SelectFunc = func(_ string, _ []huh.Option[string]) (string, error) { return selectVal, nil }
	prompt.InputFunc = func(_, _, value string) (string, error) {
		if inputVal != "" {
			return inputVal, nil
		}
		return value, nil
	}
	prompt.PortChecker = func(_ int) bool { return true }
}

func TestCreateFromImageSingle(t *testing.T) {
	stubPrompts(t, "nginx", "")
	root := t.TempDir()

	src := &fakeSource{
		searchResults: []hub.SearchResult{
			{RepoName: "nginx", IsOfficial: true, PullCount: 1000000},
		},
		imgCfg: &hub.ImageConfig{
			ExposedPorts: []int{80},
			EnvDefaults:  map[string]string{},
		},
	}

	var out bytes.Buffer
	projectPath, err := createProject(&createOptions{
		name:         "nginx",
		projectsRoot: root,
		autoGenPass:  true,
		source:       src,
		out:          &out,
	})
	if err != nil {
		t.Fatalf("createProject error: %v", err)
	}

	for _, f := range []string{"docker-compose.yml", ".env", ".gitignore", "updock.json"} {
		if _, err := os.Stat(filepath.Join(projectPath, f)); os.IsNotExist(err) {
			t.Errorf("expected %q to be created", f)
		}
	}

	if filepath.Base(projectPath) != "nginx" {
		t.Errorf("expected project dir 'nginx', got %q", filepath.Base(projectPath))
	}
}

func TestCreateFromImageChooseVersion(t *testing.T) {
	origSel, origIn, origPort := prompt.SelectFunc, prompt.InputFunc, prompt.PortChecker
	t.Cleanup(func() {
		prompt.SelectFunc = origSel
		prompt.InputFunc = origIn
		prompt.PortChecker = origPort
	})
	prompt.PortChecker = func(_ int) bool { return true }
	prompt.InputFunc = func(_, _, value string) (string, error) { return value, nil }

	// First select returns the "choose version" sentinel, second returns a tag.
	calls := 0
	prompt.SelectFunc = func(title string, _ []huh.Option[string]) (string, error) {
		calls++
		if strings.Contains(title, "version") {
			return "1.25.0", nil
		}
		return "__choose_version__", nil
	}

	root := t.TempDir()
	src := &fakeSource{
		searchResults: []hub.SearchResult{{RepoName: "nginx", PullCount: 100}},
		tags:          []hub.Tag{{Name: "1.25.0"}, {Name: "latest"}},
		imgCfg:        &hub.ImageConfig{EnvDefaults: map[string]string{}},
	}

	var out bytes.Buffer
	projectPath, err := createProject(&createOptions{
		name: "nginx", projectsRoot: root, source: src, out: &out,
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	data, _ := os.ReadFile(filepath.Join(projectPath, "updock.json"))
	if !strings.Contains(string(data), "1.25.0") {
		t.Error("chosen version 1.25.0 should be recorded in metadata")
	}
}

func TestCreateFromImageNoResults(t *testing.T) {
	stubPrompts(t, "", "")
	src := &fakeSource{searchResults: nil}

	var out bytes.Buffer
	_, err := createProject(&createOptions{
		name: "doesnotexist", projectsRoot: t.TempDir(), source: src, out: &out,
	})
	if err == nil {
		t.Error("expected error for no search results")
	}
}

func TestCreateFromImageSearchError(t *testing.T) {
	stubPrompts(t, "", "")
	src := &fakeSource{searchErr: fmt.Errorf("network down")}

	var out bytes.Buffer
	_, err := createProject(&createOptions{
		name: "nginx", projectsRoot: t.TempDir(), source: src, out: &out,
	})
	if err == nil {
		t.Error("expected error when search fails")
	}
}

func TestCreateFromImageInspectFallback(t *testing.T) {
	stubPrompts(t, "nginx", "")
	root := t.TempDir()

	src := &fakeSource{
		searchResults: []hub.SearchResult{{RepoName: "nginx", PullCount: 100}},
		inspectErr:    fmt.Errorf("registry unreachable"),
	}

	var out bytes.Buffer
	projectPath, err := createProject(&createOptions{
		name: "nginx", projectsRoot: root, autoGenPass: true, source: src, out: &out,
	})
	if err != nil {
		t.Fatalf("inspect failure should be non-fatal, got: %v", err)
	}
	if _, err := os.Stat(filepath.Join(projectPath, "docker-compose.yml")); os.IsNotExist(err) {
		t.Error("project should still be created after inspect fallback")
	}
	if !strings.Contains(out.String(), "minimal configuration") {
		t.Error("should warn about minimal configuration")
	}
}

func TestCreateFromImageInvalidName(t *testing.T) {
	stubPrompts(t, "", "")
	src := &fakeSource{}

	var out bytes.Buffer
	_, err := createProject(&createOptions{
		name: "BAD;NAME", projectsRoot: t.TempDir(), source: src, out: &out,
	})
	if err == nil {
		t.Error("expected validation error for invalid image name")
	}
}

func TestCreateFromImageNameOverride(t *testing.T) {
	stubPrompts(t, "nginx", "")
	root := t.TempDir()

	src := &fakeSource{
		searchResults: []hub.SearchResult{{RepoName: "nginx", PullCount: 100}},
		imgCfg:        &hub.ImageConfig{EnvDefaults: map[string]string{}},
	}

	var out bytes.Buffer
	projectPath, err := createProject(&createOptions{
		name: "nginx", projectName: "my-web", projectsRoot: root, source: src, out: &out,
	})
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(projectPath) != "my-web" {
		t.Errorf("expected 'my-web', got %q", filepath.Base(projectPath))
	}
}

func TestCreateFromRecipe(t *testing.T) {
	stubPrompts(t, "", "")
	root := t.TempDir()

	recipes := map[string]*recipe.Recipe{
		"wordpress": {
			Meta: recipe.Meta{Name: "wordpress", Description: "WP + MySQL"},
			Prompts: []recipe.Prompt{
				{Name: "DB_PASSWORD", Required: true, Secret: true, Generate: "password"},
				{Name: "DB_NAME", Default: "wordpress"},
			},
			Services: map[string]recipe.Service{
				"wordpress": {Image: "wordpress", DefaultTag: "latest", Environment: map[string]string{"WORDPRESS_DB_PASSWORD": "${DB_PASSWORD}"}},
				"db":        {Image: "mysql", DefaultTag: "8.0", Environment: map[string]string{"MYSQL_ROOT_PASSWORD": "${DB_PASSWORD}"}},
			},
		},
	}

	var out bytes.Buffer
	projectPath, err := createProject(&createOptions{
		name:         "wordpress",
		projectsRoot: root,
		autoGenPass:  true,
		recipes:      recipes,
		out:          &out,
	})
	if err != nil {
		t.Fatalf("createProject (recipe) error: %v", err)
	}

	if !strings.Contains(out.String(), "Using recipe") {
		t.Error("should announce recipe use")
	}

	data, err := os.ReadFile(filepath.Join(projectPath, "docker-compose.yml"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if !strings.Contains(content, "wordpress:latest") || !strings.Contains(content, "mysql:8.0") {
		t.Error("recipe should scaffold both services")
	}
}

func TestCreateFromRecipeFallsBackToImage(t *testing.T) {
	stubPrompts(t, "nginx", "")
	root := t.TempDir()

	src := &fakeSource{
		searchResults: []hub.SearchResult{{RepoName: "nginx", PullCount: 100}},
		imgCfg:        &hub.ImageConfig{EnvDefaults: map[string]string{}},
	}

	// recipes map has wordpress, but we ask for nginx, so it falls back to the image flow
	recipes := map[string]*recipe.Recipe{
		"wordpress": {Meta: recipe.Meta{Name: "wordpress"}},
	}

	var out bytes.Buffer
	_, err := createProject(&createOptions{
		name: "nginx", projectsRoot: root, source: src, recipes: recipes, out: &out,
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !strings.Contains(out.String(), "Searching Docker Hub") {
		t.Error("should fall back to image search flow")
	}
}

func TestCreateAutoIncrementsName(t *testing.T) {
	stubPrompts(t, "nginx", "")
	root := t.TempDir()

	src := &fakeSource{
		searchResults: []hub.SearchResult{{RepoName: "nginx", PullCount: 100}},
		imgCfg:        &hub.ImageConfig{EnvDefaults: map[string]string{}},
	}

	var out bytes.Buffer
	opts := &createOptions{name: "nginx", projectsRoot: root, source: src, out: &out}

	p1, err := createProject(opts)
	if err != nil {
		t.Fatal(err)
	}
	p2, err := createProject(opts)
	if err != nil {
		t.Fatal(err)
	}

	if filepath.Base(p1) != "nginx" {
		t.Errorf("first project should be 'nginx', got %q", filepath.Base(p1))
	}
	if filepath.Base(p2) != "nginx-2" {
		t.Errorf("second project should be 'nginx-2', got %q", filepath.Base(p2))
	}
}

func TestUserRecipesDir(t *testing.T) {
	t.Run("with XDG", func(t *testing.T) {
		custom := filepath.Join(t.TempDir(), "data")
		t.Setenv("XDG_DATA_HOME", custom)
		dir := userRecipesDir()
		if dir != filepath.Join(custom, "updock", "recipes") {
			t.Errorf("got %q", dir)
		}
	})
	t.Run("without XDG", func(t *testing.T) {
		t.Setenv("XDG_DATA_HOME", "")
		dir := userRecipesDir()
		if dir != "" && !strings.Contains(dir, "updock") {
			t.Errorf("got %q", dir)
		}
	})
}
