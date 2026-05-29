# Troubleshooting

Start with the built-in checker, which diagnoses the most common problems and
prints the exact fix for each:

```bash
updock doctor
```

## Preflight failures

`updock doctor` (and the start of every run) checks four things.

### Docker is not installed

```
✗ Docker installed - docker is not on your PATH
  Fix: Install Docker: https://docs.docker.com/get-docker/
```

Install Docker, then re-run `updock doctor`.

### Docker Compose is not available

```
✗ Docker Compose available - neither 'docker compose' plugin nor 'docker-compose' found
```

Install the Compose plugin: https://docs.docker.com/compose/install/

### The Docker daemon is not running

```
✗ Docker daemon running - cannot connect to the Docker daemon
  Fix: Start Docker Desktop or the Docker daemon
```

On Linux: `sudo systemctl start docker`. On macOS or Windows: start Docker
Desktop.

### You cannot reach the Docker socket

```
✗ Docker socket accessible - cannot reach /var/run/docker.sock
  Fix: Add yourself to the docker group: sudo usermod -aG docker <you> && newgrp docker
```

Adding yourself to the `docker` group grants root-equivalent access to your
machine. updock states this so you can make an informed choice; it never does it
for you. See the [security model](security.md).

## Common errors

### `command not found: updock`

The binary is not on your `PATH`. If you used `go install`, add your Go bin
directory:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

See [Installation](installation.md#go-install).

### `no project with number N`

The number does not match any project. Run `updock ls` to see the current
numbers. They are assigned alphabetically and can change as you add or remove
projects.

### `no images found for "..."`

The search returned nothing. Check the spelling or try a shorter, more general
term. You can browse results first with `updock search <name>`.

### A port is already in use

updock detects this during configuration and proposes the next free port
automatically, telling you what it picked. To change a port later, run
`updock config <n>` and set a different host port, then restart.

### Rate limited by Docker Hub

Anonymous Docker Hub access is rate limited. If you hit the limit, updock tells
you and you can wait for the window to reset, or authenticate with Docker to
raise the limit.

### "could not read image metadata"

updock could not inspect the image (it may be private or unusual). It continues
with a minimal configuration, so you can still scaffold and run the project;
you may just need to set ports and environment variables yourself.

## Still stuck?

Open a discussion or issue:
https://github.com/amrelsagaei/updock/issues
