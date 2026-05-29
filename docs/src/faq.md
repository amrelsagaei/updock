# FAQ

### Why not just write `docker compose` myself?

You can, and updock's output is plain Compose so you still can. updock removes
the tedious parts: finding the right image, picking a version, looking up which
ports and environment variables it needs, generating passwords, and wiring
multi-service apps together. One word in, a running app out.

### Why control by number instead of by name?

Because it is faster and you never have to remember or retype names. `updock ls`
numbers everything, and every command takes that number: `updock logs 2`,
`updock stop 2`. The numbers are stable and assigned alphabetically.

### Where are my files?

Projects live in `~/updock/<name>/` (configurable). Global config is at
`~/.config/updock/config.toml` and your own recipes go in
`~/.local/share/updock/recipes/`. See [Projects](projects.md).

### Is updock locked in? Can I stop using it?

No lock-in. The generated `docker-compose.yml` and `.env` are standard files
that work with plain `docker compose`. You can keep the folders and manage them
directly at any time.

### Does updock store or send my secrets anywhere?

No. Secrets are written only to the project's `.env` file (`0600`, gitignored),
masked in all output, and never logged. updock has no telemetry and never phones
home. See the [security model](security.md).

### Do I need root?

No. updock never needs root for itself; only Docker needs privileges. updock
never runs sudo on your behalf.

### How do I run a specific version of an image?

At the picker, choose the "choose version" option to see the tag list, sorted
newest-first with `latest` pinned to the top.

### Can I add an app that needs multiple containers?

Yes, through recipes. Several are built in (WordPress, Nextcloud, Gitea, and
more), and you can add your own. See [Recipes](recipes.md).

### Does it work offline?

Search and image inspection need network access to Docker Hub. Once a project is
scaffolded, lifecycle commands (`up`, `stop`, `logs`, and so on) work against
your local Docker. Built-in recipes are embedded in the binary, so they need no
network.

### Which platforms are supported?

Linux, macOS, and Windows, on amd64 and arm64. See
[Installation](installation.md).
