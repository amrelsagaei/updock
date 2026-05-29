# Security Policy

## Supported Versions

| Version | Supported |
|---------|-----------|
| latest  | Yes       |
| < latest | No       |

Only the most recent release receives security updates.

## Reporting a Vulnerability

**Do NOT open a public GitHub issue for security vulnerabilities.**

Instead, email **amr.elsagaei@gmail.com** with:

1. A description of the vulnerability.
2. Steps to reproduce.
3. The potential impact.
4. Any suggested fix (optional).

## Response Timeline

- **Acknowledgment**: Within 48 hours.
- **Assessment**: Within 72 hours we'll confirm whether it's a valid security issue.
- **Fix plan**: Within 7 days of confirmation.
- **Release**: As soon as the fix is ready and tested.

## Disclosure Policy

We follow **coordinated disclosure**:

- We work with you to understand and fix the issue before any public disclosure.
- You will be credited in the release notes (unless you prefer to remain anonymous).
- We ask that you do not disclose publicly until a fix is released.

## Scope

The following are considered security issues:

- Shell injection through user input
- Secret values leaking to stdout, logs, or generated files other than `.env`
- `.env` file created with permissions other than `0600`
- Secrets appearing in `docker-compose.yml` or `updock.json`
- Bypassing input validation
- Unsafe Docker operations (e.g., running containers with unintended privileges)
- Supply chain issues (compromised dependencies, unsigned releases)

The following are NOT security issues:

- The Docker group granting root-equivalent access (this is documented behavior)
- Vulnerabilities in Docker images that users choose to run (use `updock` with `scan_before_run` enabled)
