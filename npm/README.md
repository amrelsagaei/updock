# updock

Run any Docker app from one word. No YAML, no flag-hunting.

You type a name. updock finds the image, lets you pick a version, asks only the
few things that actually matter (ports, passwords, env vars), writes the Compose
file and `.env`, and brings it up. After that you control everything by number.

## Install

```bash
npm install -g updock
# or run it once, without installing:
npx updock postgres
```

Other install methods (Homebrew, go install, prebuilt binaries) are in the docs.

## Usage

```bash
updock postgres      # search, pick a version, configure, run
updock ls            # list your projects, numbered
updock logs 2        # everything after is controlled by number
updock stop 2
```

## Links

- Documentation: https://amrelsagaei.github.io/updock/
- Source and issues: https://github.com/amrelsagaei/updock
- License: MIT
