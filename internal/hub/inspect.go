// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package hub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

// ImageConfig holds the metadata extracted from a Docker image without pulling it.
type ImageConfig struct {
	ExposedPorts []int
	EnvDefaults  map[string]string
	Labels       map[string]string
	Volumes      []string
	Cmd          []string
	Entrypoint   []string
}

type tokenResponse struct {
	Token string `json:"token"`
}

type manifest struct {
	SchemaVersion int    `json:"schemaVersion"`
	MediaType     string `json:"mediaType"`
	Config        struct {
		Digest string `json:"digest"`
	} `json:"config"`
}

type manifestList struct {
	SchemaVersion int    `json:"schemaVersion"`
	MediaType     string `json:"mediaType"`
	Manifests     []struct {
		Digest    string `json:"digest"`
		MediaType string `json:"mediaType"`
		Platform  struct {
			Architecture string `json:"architecture"`
			OS           string `json:"os"`
		} `json:"platform"`
	} `json:"manifests"`
}

type imageConfigBlob struct {
	Config struct {
		ExposedPorts map[string]struct{} `json:"ExposedPorts"`
		Env          []string            `json:"Env"`
		Labels       map[string]string   `json:"Labels"`
		Volumes      map[string]struct{} `json:"Volumes"`
		Cmd          []string            `json:"Cmd"`
		Entrypoint   []string            `json:"Entrypoint"`
	} `json:"config"`
}

// Inspect fetches image metadata from the registry without pulling the full image.
func (c *Client) Inspect(image, tag string) (*ImageConfig, error) {
	return c.inspectWith(image, tag, authURL, registryURL)
}

func (c *Client) inspectWith(image, tag, authBase, regBase string) (*ImageConfig, error) {
	repo := normalizeRepo(image)

	token, err := c.getRegistryTokenFrom(repo, authBase)
	if err != nil {
		return nil, fmt.Errorf("getting auth token: %w", err)
	}

	configDigest, err := c.resolveConfigDigestFrom(repo, tag, token, regBase)
	if err != nil {
		return nil, fmt.Errorf("resolving manifest: %w", err)
	}

	cfg, err := c.fetchConfigBlobFrom(repo, configDigest, token, regBase)
	if err != nil {
		return nil, fmt.Errorf("fetching config: %w", err)
	}

	return cfg, nil
}

func normalizeRepo(image string) string {
	if !strings.Contains(image, "/") {
		return "library/" + image
	}
	return image
}

func (c *Client) getRegistryTokenFrom(repo, authBase string) (string, error) {
	u := fmt.Sprintf("%s?service=registry.docker.io&scope=repository:%s:pull", authBase, repo)

	var resp tokenResponse
	if err := c.getJSON(u, &resp); err != nil {
		return "", err
	}
	if resp.Token == "" {
		return "", fmt.Errorf("empty token received")
	}
	return resp.Token, nil
}

func (c *Client) registryGet(url, token, accept string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Authorization", "Bearer "+token)
	if accept != "" {
		req.Header.Set("Accept", accept)
	}
	return c.http.Do(req)
}

func (c *Client) resolveConfigDigestFrom(repo, tag, token, regBase string) (string, error) {
	u := fmt.Sprintf("%s/v2/%s/manifests/%s", regBase, repo, tag)

	accept := strings.Join([]string{
		"application/vnd.docker.distribution.manifest.v2+json",
		"application/vnd.oci.image.manifest.v1+json",
		"application/vnd.docker.distribution.manifest.list.v2+json",
		"application/vnd.oci.image.index.v1+json",
	}, ", ")

	resp, err := c.registryGet(u, token, accept)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("manifest request returned HTTP %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")

	if strings.Contains(ct, "manifest.list") || strings.Contains(ct, "image.index") {
		return c.resolveFromManifestList(resp, repo, token, regBase)
	}

	var m manifest
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return "", fmt.Errorf("decoding manifest: %w", err)
	}
	if m.Config.Digest == "" {
		return "", fmt.Errorf("manifest has no config digest")
	}
	return m.Config.Digest, nil
}

func (c *Client) resolveFromManifestList(resp *http.Response, repo, token, regBase string) (string, error) {
	var ml manifestList
	if err := json.NewDecoder(resp.Body).Decode(&ml); err != nil {
		return "", fmt.Errorf("decoding manifest list: %w", err)
	}

	targetOS := runtime.GOOS
	targetArch := runtime.GOARCH

	for _, m := range ml.Manifests {
		if m.Platform.OS == targetOS && m.Platform.Architecture == targetArch {
			return c.fetchManifestDigest(repo, m.Digest, token, regBase)
		}
	}

	if len(ml.Manifests) > 0 {
		return c.fetchManifestDigest(repo, ml.Manifests[0].Digest, token, regBase)
	}

	return "", fmt.Errorf("no suitable platform found in manifest list")
}

func (c *Client) fetchManifestDigest(repo, digest, token, regBase string) (string, error) {
	u := fmt.Sprintf("%s/v2/%s/manifests/%s", regBase, repo, digest)

	resp, err := c.registryGet(u, token, "application/vnd.docker.distribution.manifest.v2+json")
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("platform manifest returned HTTP %d", resp.StatusCode)
	}

	var m manifest
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return "", err
	}
	if m.Config.Digest == "" {
		return "", fmt.Errorf("platform manifest has no config digest")
	}
	return m.Config.Digest, nil
}

func (c *Client) fetchConfigBlobFrom(repo, digest, token, regBase string) (*ImageConfig, error) {
	u := fmt.Sprintf("%s/v2/%s/blobs/%s", regBase, repo, digest)

	resp, err := c.registryGet(u, token, "")
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("config blob returned HTTP %d", resp.StatusCode)
	}

	var blob imageConfigBlob
	if err := json.NewDecoder(resp.Body).Decode(&blob); err != nil {
		return nil, fmt.Errorf("decoding config blob: %w", err)
	}

	return parseConfigBlob(&blob), nil
}

func parseConfigBlob(blob *imageConfigBlob) *ImageConfig {
	cfg := &ImageConfig{
		EnvDefaults: make(map[string]string),
		Labels:      blob.Config.Labels,
	}

	for portSpec := range blob.Config.ExposedPorts {
		port := strings.TrimSuffix(portSpec, "/tcp")
		port = strings.TrimSuffix(port, "/udp")
		if p, err := strconv.Atoi(port); err == nil {
			cfg.ExposedPorts = append(cfg.ExposedPorts, p)
		}
	}
	sort.Ints(cfg.ExposedPorts)

	for _, env := range blob.Config.Env {
		if k, v, ok := strings.Cut(env, "="); ok {
			cfg.EnvDefaults[k] = v
		}
	}

	for vol := range blob.Config.Volumes {
		cfg.Volumes = append(cfg.Volumes, vol)
	}
	sort.Strings(cfg.Volumes)

	cfg.Cmd = blob.Config.Cmd
	cfg.Entrypoint = blob.Config.Entrypoint

	if cfg.Labels == nil {
		cfg.Labels = make(map[string]string)
	}

	return cfg
}
