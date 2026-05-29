// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package docker

import "strings"

// ImageHint provides known env var metadata for popular Docker images.
type ImageHint struct {
	Required    []string
	Common      []string
	Description map[string]string
}

var imageHints = map[string]ImageHint{
	"postgres": {
		Required: []string{"POSTGRES_PASSWORD"},
		Common:   []string{"POSTGRES_USER", "POSTGRES_DB"},
		Description: map[string]string{
			"POSTGRES_PASSWORD": "Superuser password (required)",
			"POSTGRES_USER":     "Superuser name (default: postgres)",
			"POSTGRES_DB":       "Default database name (default: same as user)",
		},
	},
	"mysql": {
		Required: []string{"MYSQL_ROOT_PASSWORD"},
		Common:   []string{"MYSQL_DATABASE", "MYSQL_USER", "MYSQL_PASSWORD"},
		Description: map[string]string{
			"MYSQL_ROOT_PASSWORD": "Root password (required)",
			"MYSQL_DATABASE":      "Database to create on startup",
			"MYSQL_USER":          "Non-root user to create",
			"MYSQL_PASSWORD":      "Password for the non-root user",
		},
	},
	"mariadb": {
		Required: []string{"MARIADB_ROOT_PASSWORD"},
		Common:   []string{"MARIADB_DATABASE", "MARIADB_USER", "MARIADB_PASSWORD"},
		Description: map[string]string{
			"MARIADB_ROOT_PASSWORD": "Root password (required)",
			"MARIADB_DATABASE":      "Database to create on startup",
			"MARIADB_USER":          "Non-root user to create",
			"MARIADB_PASSWORD":      "Password for the non-root user",
		},
	},
	"redis": {
		Common: []string{"REDIS_PASSWORD"},
		Description: map[string]string{
			"REDIS_PASSWORD": "Optional authentication password",
		},
	},
	"mongo": {
		Common: []string{"MONGO_INITDB_ROOT_USERNAME", "MONGO_INITDB_ROOT_PASSWORD", "MONGO_INITDB_DATABASE"},
		Description: map[string]string{
			"MONGO_INITDB_ROOT_USERNAME": "Admin username",
			"MONGO_INITDB_ROOT_PASSWORD": "Admin password",
			"MONGO_INITDB_DATABASE":      "Initial database name",
		},
	},
	"wordpress": {
		Required: []string{"WORDPRESS_DB_PASSWORD"},
		Common:   []string{"WORDPRESS_DB_HOST", "WORDPRESS_DB_USER", "WORDPRESS_DB_NAME"},
		Description: map[string]string{
			"WORDPRESS_DB_HOST":     "Database host (default: mysql)",
			"WORDPRESS_DB_USER":     "Database user (default: root)",
			"WORDPRESS_DB_PASSWORD": "Database password (required)",
			"WORDPRESS_DB_NAME":     "Database name (default: wordpress)",
		},
	},
	"nextcloud": {
		Common: []string{"NEXTCLOUD_ADMIN_USER", "NEXTCLOUD_ADMIN_PASSWORD", "MYSQL_HOST", "MYSQL_DATABASE", "MYSQL_USER", "MYSQL_PASSWORD"},
		Description: map[string]string{
			"NEXTCLOUD_ADMIN_USER":     "Admin username",
			"NEXTCLOUD_ADMIN_PASSWORD": "Admin password",
		},
	},
	"gitea": {
		Common: []string{"GITEA__database__DB_TYPE", "GITEA__database__HOST", "GITEA__database__NAME", "GITEA__database__USER", "GITEA__database__PASSWD"},
		Description: map[string]string{
			"GITEA__database__DB_TYPE": "Database type (postgres, mysql, sqlite3)",
		},
	},
}

// GetHint returns the hint for a given image name, matching against the
// base name (without namespace or tag).
func GetHint(image string) (ImageHint, bool) {
	base := image
	if idx := strings.LastIndex(image, "/"); idx >= 0 {
		base = image[idx+1:]
	}
	base = strings.ToLower(base)

	hint, ok := imageHints[base]
	return hint, ok
}

// IsSecretKey returns true if the env var name looks like a secret.
func IsSecretKey(key string) bool {
	upper := strings.ToUpper(key)
	for _, pattern := range []string{"PASSWORD", "SECRET", "TOKEN", "KEY", "CREDENTIALS"} {
		if strings.Contains(upper, pattern) {
			return true
		}
	}
	return false
}
