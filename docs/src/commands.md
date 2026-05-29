# Command reference

Every command that targets a project takes a **number** from `updock ls`. There
are no flags to memorize, just numbers.

| Command | What it does |
|---|---|
| [`updock <name>`](#updock-name) | Search, pick, configure, and run an app |
| [`updock search <name>`](#updock-search) | Show ranked Docker Hub results only |
| [`updock ls`](#updock-ls) | List all projects with live state |
| [`updock status <n>`](#updock-status) | Show details for one project |
| [`updock up <n>`](#updock-up) | Start a project |
| [`updock stop <n>`](#updock-stop) | Stop containers, keep data |
| [`updock restart <n>`](#updock-restart) | Restart a project |
| [`updock rebuild <n>`](#updock-rebuild) | Rebuild and recreate containers |
| [`updock down <n>`](#updock-down) | Remove containers, keep folder and data |
| [`updock logs <n>`](#updock-logs) | Show or follow logs |
| [`updock open <n>`](#updock-open) | Open the mapped port in a browser |
| [`updock config <n>`](#updock-config) | Re-run the configuration prompts |
| [`updock rm <n>`](#updock-rm) | Delete a project (asks first) |
| [`updock doctor`](#updock-doctor) | Run all environment checks |
| [`updock version`](#updock-version) | Print version and build info |

Run `updock <command> --help` for the same information at the terminal.

## updock &lt;name&gt;

Search, pick a version, configure, and run an app. This is the main entry point;
see [Usage](usage.md) for the full pipeline.

```bash
updock postgres
updock juice-shop --name ctf-lab    # override the project name
```

Flags:

- `--name <string>` - set the project name instead of deriving it from the image.

## updock search

Search Docker Hub and print the ranked matches without running anything. Useful
for browsing before you commit.

```bash
updock search nginx
```

## updock ls

List every project with its number, image, live state, and ports. The number is
what you pass to every other command.

```bash
updock ls
```

## updock status

Show full details for one project: image, live state, port mappings,
environment variables (secrets masked as `••••••`), and the project path.

```bash
updock status 1
```

## updock up

Start the project. Pulls images if needed and brings the stack up in the
background, then prints the URLs you can open.

```bash
updock up 1
```

## updock stop

Stop the containers but keep the project and its data. Start it again with
`updock up <n>`.

```bash
updock stop 1
```

## updock restart

Restart the containers. Use this to apply changes that only need a fresh start.

```bash
updock restart 1
```

## updock rebuild

Rebuild and recreate the containers. Use this after changing the image or
configuration.

```bash
updock rebuild 1
```

## updock down

Remove the containers but keep the project folder and the `data/` directory, so
`updock up <n>` brings everything back.

```bash
updock down 1
```

## updock logs

Show the logs for a project. Pass `-f` to follow them live.

```bash
updock logs 1
updock logs 1 -f
```

Flags:

- `-f`, `--follow` - stream new log output as it arrives.

## updock open

Open the project's first mapped port in your default browser. Set
`browser_command` in the [config](configuration.md) to use a specific browser.

```bash
updock open 1
```

## updock config

Re-run the configuration prompts for a project. Current values are pre-filled,
so press enter to keep them. updock regenerates the `.env` file and offers to
restart so the changes take effect.

```bash
updock config 1
```

## updock rm

Delete a project: remove its containers and volumes, then delete the project
folder. updock asks for confirmation first unless you pass `--yes`. This cannot
be undone.

```bash
updock rm 1
updock rm 1 --yes    # skip the confirmation
```

Flags:

- `-y`, `--yes` - skip the confirmation prompt.

## updock doctor

Run all environment checks: Docker installed, Compose available, daemon running,
socket accessible. For each failure it prints the problem and the fix. See
[Troubleshooting](troubleshooting.md).

```bash
updock doctor
```

## updock version

Print the version, commit, build date, Go version, and platform.

```bash
updock version
```
