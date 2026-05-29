# Configuration

updock works with zero configuration. Every setting has a sensible default, so
the config file is entirely optional. Create it only when you want to change a
default.

## File location

| OS | Path |
|---|---|
| Linux / macOS | `~/.config/updock/config.toml` |
| With `XDG_CONFIG_HOME` set | `$XDG_CONFIG_HOME/updock/config.toml` |

The file is TOML. If it does not exist, updock uses the defaults below.

## Reference

| Key | Type | Default | Description |
|---|---|---|---|
| `projects_root` | path | `~/updock` | Where project folders, Compose files, and `.env` files are written. |
| `browser_command` | string | `open` (macOS), `xdg-open` (Linux) | Command used by `updock open` to launch URLs. |
| `auto_generate_passwords` | bool | `true` | Generate strong passwords for required secret fields. Set to `false` to be prompted for every secret. |

Values come from the config file when present, otherwise from these built-in
defaults. Invalid values (for example an empty `projects_root`) are rejected on
load with a clear error.

## Annotated example

```toml
# ~/.config/updock/config.toml
# All keys are optional. The values shown here are the defaults.

# Where updock writes each project (compose file, .env, data/).
projects_root = "~/updock"

# Command used to open service URLs. Linux default is "xdg-open".
browser_command = "open"

# Generate strong random passwords for required secrets instead of prompting.
auto_generate_passwords = true
```

## Per-project values

Settings above are global. The values for a specific project (its ports,
environment variables, and secrets) live inside that project's folder, in
`docker-compose.yml` and `.env`. To change them, edit those files directly or
run `updock config <n>` to re-run the prompts. See [Projects](projects.md).
