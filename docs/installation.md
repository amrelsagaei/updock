# Installation

updock ships as a single static binary with no runtime dependencies (besides
Docker itself, which you already have). Pick whichever method fits your setup.

> Recommended: **Homebrew** on macOS/Linux, **`go install`** if you have a Go
> toolchain, or a **prebuilt binary** from the releases page for everyone else.

## Quick install

| Platform | Command |
|---|---|
| macOS / Linux (Homebrew) | `brew install amrelsagaei/tap/updock` |
| Go toolchain | `go install github.com/amrelsagaei/updock/cmd/updock@latest` |
| npm | `npm install -g updock` (or `npx updock <name>`) |
| Debian / Ubuntu | download the `.deb` from releases, then `sudo dpkg -i updock_*.deb` |

After installing, confirm it works and check your environment:

```bash
updock version
updock doctor
```

## macOS

### Homebrew

```bash
brew tap amrelsagaei/tap
brew install updock
```

Upgrade and uninstall:

```bash
brew upgrade updock
brew uninstall updock
```

If macOS Gatekeeper flags the binary on first run, the Homebrew install removes
the quarantine attribute automatically. For a manually downloaded binary, see
[Manual download](#manual-download) below.

## Linux

| Method | Install | Notes |
|---|---|---|
| Debian / Ubuntu (`.deb`) | `sudo dpkg -i updock_*.deb` or `sudo apt install ./updock_*.deb` | from the releases page |
| RPM (Fedora / RHEL) | `sudo rpm -i updock_*.rpm` | from the releases page |
| Alpine (`.apk`) | `sudo apk add --allow-untrusted updock_*.apk` | from the releases page |

## Windows

Download the `updock_windows_<arch>.zip` from the
[releases page](https://github.com/amrelsagaei/updock/releases), extract
`updock.exe`, and put it on your `PATH`. updock drives the `docker` CLI, so make
sure Docker Desktop is installed and running.

## go install

```bash
go install github.com/amrelsagaei/updock/cmd/updock@latest
```

This compiles from source straight into your Go bin directory. No release
tooling required, so it always tracks the latest tagged version.

> **PATH note.** `go install` puts binaries in `$GOBIN`, or `$(go env GOPATH)/bin`
> if `GOBIN` is unset (default `$HOME/go/bin`). If `updock: command not found`,
> add that directory to your `PATH`:
>
> ```bash
> export PATH="$PATH:$(go env GOPATH)/bin"
> ```

## npm

```bash
npm install -g updock     # global install
npx updock postgres       # or run once without installing
```

The npm package is a thin launcher that downloads the matching prebuilt binary
for your platform.

## Manual download

1. Download the archive for your OS and architecture from the
   [releases page](https://github.com/amrelsagaei/updock/releases)
   (`updock_<os>_<arch>.tar.gz`, or `.zip` on Windows).
2. **Verify it** (see [Verifying a download](#verifying-a-download)).
3. Extract and move the binary onto your `PATH`:

   ```bash
   tar -xzf updock_*.tar.gz
   sudo mv updock /usr/local/bin/
   chmod +x /usr/local/bin/updock
   ```

## Verifying a download

Every release ships a `checksums.txt`, a cosign signature, and a certificate.
Signing is keyless through GitHub Actions OIDC, so there are no long-lived keys.

```bash
cosign verify-blob \
  --certificate checksums.txt.pem \
  --signature   checksums.txt.sig \
  --certificate-identity-regexp 'https://github.com/amrelsagaei/updock' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com \
  checksums.txt
```

Then check your archive against the verified checksums:

```bash
sha256sum -c checksums.txt --ignore-missing
```

A software bill of materials (SBOM) is attached to every release as well.

## Shell completion

updock generates completion scripts for bash, zsh, fish, and PowerShell:

```bash
# bash
updock completion bash | sudo tee /etc/bash_completion.d/updock > /dev/null

# zsh (ensure compinit is loaded)
updock completion zsh > "${fpath[1]}/_updock"

# fish
updock completion fish > ~/.config/fish/completions/updock.fish
```

Run `updock completion --help` for per-shell details.

## Upgrading

Use the same channel you installed from: `brew upgrade updock`,
`go install ...@latest`, your package manager's update command, or download a
newer binary.

## Uninstalling

```bash
brew uninstall updock          # Homebrew
sudo apt remove updock         # apt
rm "$(command -v updock)"      # manual / go install
```

updock keeps your projects in `~/updock/` and config in `~/.config/updock/`.
Remove those directories too if you want a clean slate (this deletes your
generated Compose files and `.env` files).
