# Security model

updock runs other people's code on your machine and handles your secrets, so it
treats both carefully. This page explains the security model. To **report a
vulnerability**, see [SECURITY.md](https://github.com/amrelsagaei/updock/blob/main/SECURITY.md) instead.

## Running the tool

- **updock never needs root for itself.** Only Docker needs privileges. updock
  does not ask for sudo and does not run as root.
- **The `docker` group is root in disguise.** Anyone who can reach the Docker
  socket can get root on the host. When `updock doctor` suggests adding yourself
  to the `docker` group, it states this plainly so you make an informed choice.
  It never adds you silently.

## Secrets

- **Secrets live in `.env`, never in the Compose file.** Values you mark as
  secret (passwords, tokens, keys) are written to `.env` and referenced from
  `docker-compose.yml` as `${VAR}`. They are never inlined into the YAML,
  because anything in the YAML shows up in `docker inspect` and the process
  list.
- **`.env` is locked down.** It is created with `0600` permissions (owner read
  and write only) and added to the project's `.gitignore`.
- **Secrets are masked in output.** updock never prints a secret to the terminal
  or logs. In `updock status` and other views, secret values show as `••••••`.
- **Generated passwords are strong.** When updock generates a password it uses
  `crypto/rand` (never `math/rand`), producing a long, random value.

## Running images safely

- **You are executing third-party code.** updock shows whether an image is
  official or popular before you run it.
- **Pin by digest.** updock can record an image's digest so a rebuild uses
  exactly what you reviewed.
- **Optional scanning.** With `scan_before_run` enabled in the
  [config](configuration.md), updock can run a vulnerability scan (trivy or
  Docker Scout) before the first run.

## Input handling

- **No shell injection.** updock builds Docker commands with argument arrays,
  not by gluing your input into a shell string. An image name with a semicolon
  in it cannot run a command. This is enforced by an automated test that fails
  the build if any shell-string command construction is introduced.
- **Everything is validated.** Image names, tags, project names, ports, and
  environment variable names are checked against strict patterns before anything
  runs.

## Supply chain

- **Signed releases.** Every release is signed with cosign using keyless signing
  through GitHub Actions OIDC. No long-lived signing keys to leak. See
  [Verifying a download](installation.md#verifying-a-download).
- **Checksums and SBOM.** Each release ships a `checksums.txt` and a software
  bill of materials, so anyone can verify what they downloaded.
- **No telemetry.** updock does not phone home. There is no tracking.
