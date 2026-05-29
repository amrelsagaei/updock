// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package prompt

import (
	"testing"

	"github.com/amrelsagaei/updock/internal/hub"
	"github.com/amrelsagaei/updock/internal/project"
)

func TestClassifyEnvVarsPostgres(t *testing.T) {
	imgCfg := &hub.ImageConfig{
		EnvDefaults: map[string]string{
			"POSTGRES_PASSWORD": "",
			"POSTGRES_USER":     "postgres",
			"POSTGRES_DB":       "",
			"PATH":              "/usr/local/bin",
			"GOPATH":            "/go",
		},
	}

	vars := ClassifyEnvVars("postgres", imgCfg)

	var pwVar *project.EnvVar
	for i := range vars {
		if vars[i].Key == "POSTGRES_PASSWORD" {
			pwVar = &vars[i]
			break
		}
	}

	if pwVar == nil {
		t.Fatal("expected POSTGRES_PASSWORD in vars")
	}
	if !pwVar.Required {
		t.Error("POSTGRES_PASSWORD should be required")
	}
	if !pwVar.Secret {
		t.Error("POSTGRES_PASSWORD should be secret")
	}

	for _, v := range vars {
		if v.Key == "PATH" || v.Key == "GOPATH" {
			t.Errorf("boring var %q should be filtered out", v.Key)
		}
	}
}

func TestClassifyEnvVarsUnknownImage(t *testing.T) {
	imgCfg := &hub.ImageConfig{
		EnvDefaults: map[string]string{
			"PORT":       "3000",
			"API_SECRET": "",
			"NODE_ENV":   "production",
		},
	}

	vars := ClassifyEnvVars("some-unknown-app", imgCfg)

	var secretVar *project.EnvVar
	for i := range vars {
		if vars[i].Key == "API_SECRET" {
			secretVar = &vars[i]
			break
		}
	}

	if secretVar == nil {
		t.Fatal("expected API_SECRET")
	}
	if !secretVar.Secret {
		t.Error("API_SECRET should be detected as secret")
	}
	if !secretVar.Required {
		t.Error("empty secret should be required")
	}
}

func TestClassifyEnvVarsEmpty(t *testing.T) {
	imgCfg := &hub.ImageConfig{
		EnvDefaults: map[string]string{},
	}

	vars := ClassifyEnvVars("nginx", imgCfg)
	if len(vars) != 0 {
		t.Errorf("expected 0 vars for nginx with no env, got %d", len(vars))
	}
}

func TestCollectEnvVarsAutoGenerate(t *testing.T) {
	vars := []project.EnvVar{
		{Key: "POSTGRES_PASSWORD", Value: "", Secret: true, Required: true},
		{Key: "POSTGRES_USER", Value: "postgres"},
	}

	result, err := CollectEnvVars(vars, true)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if result[0].Value == "" {
		t.Error("password should have been auto-generated")
	}
	if len(result[0].Value) < 32 {
		t.Errorf("generated password should be 32+ chars, got %d", len(result[0].Value))
	}
	if result[1].Value != "postgres" {
		t.Errorf("non-empty value should be preserved, got %q", result[1].Value)
	}
}

func TestCollectEnvVarsRequiredNoAutoGenerate(t *testing.T) {
	orig := InputFunc
	defer func() { InputFunc = orig }()
	InputFunc = func(_, _, _ string) (string, error) {
		return "my-password", nil
	}

	vars := []project.EnvVar{
		{Key: "DB_PASSWORD", Value: "", Secret: true, Required: true},
	}

	result, err := CollectEnvVars(vars, false)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if result[0].Value != "my-password" {
		t.Errorf("expected 'my-password', got %q", result[0].Value)
	}
}

func TestCollectEnvVarsRequiredEmpty(t *testing.T) {
	orig := InputFunc
	defer func() { InputFunc = orig }()
	InputFunc = func(_, _, _ string) (string, error) {
		return "", nil
	}

	vars := []project.EnvVar{
		{Key: "DB_PASSWORD", Value: "", Secret: true, Required: true},
	}

	_, err := CollectEnvVars(vars, false)
	if err == nil {
		t.Error("expected error when required var is left empty")
	}
}

func TestCollectEnvVarsOptional(t *testing.T) {
	orig := InputFunc
	defer func() { InputFunc = orig }()
	InputFunc = func(_, _, _ string) (string, error) {
		return "custom-value", nil
	}

	vars := []project.EnvVar{
		{Key: "OPTIONAL_VAR", Value: ""},
	}

	result, err := CollectEnvVars(vars, true)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if result[0].Value != "custom-value" {
		t.Errorf("expected 'custom-value', got %q", result[0].Value)
	}
}

func TestIsBoringEnvVar(t *testing.T) {
	if !isBoringEnvVar("PATH") {
		t.Error("PATH should be boring")
	}
	if isBoringEnvVar("DATABASE_URL") {
		t.Error("DATABASE_URL should not be boring")
	}
}
