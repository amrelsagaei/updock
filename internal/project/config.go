// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package project

// Config holds all the configuration collected for a project,
// bridging the configuration prompts and the scaffold module.
type Config struct {
	Image       string
	Tag         string
	Digest      string
	ProjectName string
	Ports       []PortMapping
	Env         []EnvVar
	Volumes     []VolumeMapping
	Labels      map[string]string
}

// PortMapping describes a host:container port binding.
type PortMapping struct {
	Host      int
	Container int
	Protocol  string
}

// EnvVar describes an environment variable with metadata.
type EnvVar struct {
	Key      string
	Value    string
	Secret   bool
	Required bool
}

// VolumeMapping describes a host:container volume binding.
type VolumeMapping struct {
	HostPath      string
	ContainerPath string
}
