# Usage

updock has one core idea: you name an app, and updock does the rest. After that,
you manage everything by number. This page explains both halves.

## The pipeline

Running `updock <name>` walks through these stages:

1. **Preflight.** updock checks that Docker, the Compose plugin, and the daemon
   are present and reachable. If something is off it stops and tells you how to
   fix it. It never changes your system silently. Run the same checks any time
   with `updock doctor`.
2. **Search.** updock queries Docker Hub for your term and ranks the matches.
   Official images come first, then popular ones, then the closest name matches.
3. **Select.** You get a numbered, keyboard-driven list. Option 1 is the best
   match at `latest`. Option 2 is the same image with "choose version", which
   lists tags newest-first with proper semantic-version ordering. Pick by typing
   a number or using the arrow keys.
4. **Inspect.** updock reads the image's exposed ports, environment variables,
   and volumes straight from the registry, without pulling the whole image.
5. **Configure.** updock asks only what matters: which host port to use
   (auto-resolving conflicts), any required passwords (it offers to generate
   strong ones), common environment variables, and where to keep volumes.
6. **Scaffold.** updock creates a project folder with a `docker-compose.yml`, a
   locked-down `.env`, a `.gitignore`, and an `updock.json` metadata file.
7. **Up.** updock pulls the image and starts the stack in the background, then
   prints the URLs you can open.

If updock cannot read an image's metadata (a private or unusual image), it says
so and continues with a minimal configuration rather than failing.

## Choosing a version

At the select step, pick the "choose version" option to see the tag list. Tags
are sorted newest-first: `latest` is pinned to the top, then valid semantic
versions in descending order (so `1.10.0` sits above `1.9.0`), then any
remaining tags alphabetically.

## Naming projects

By default the project name comes from the image (for example `juice-shop`). If
that name is taken, updock appends a number (`juice-shop-2`). Override it at
create time:

```bash
updock juice-shop --name ctf-lab
```

## Control by number

Once a project exists, `updock ls` numbers everything and every action takes
that number:

```
  ┌───┬────────────┬──────────────────────────────┬─────────┬──────────────┐
  │ # │ PROJECT    │ IMAGE                        │ STATE   │ PORTS        │
  ├───┼────────────┼──────────────────────────────┼─────────┼──────────────┤
  │ 1 │ juice-shop │ bkimminich/juice-shop:latest │ running │ 3000->3000   │
  │ 2 │ postgres   │ postgres:16                  │ stopped │ 5432->5432   │
  └───┴────────────┴──────────────────────────────┴─────────┴──────────────┘
```

The numbers are stable and assigned alphabetically. State is queried live from
Docker each time you run `updock ls`, so it always reflects reality.

See [Commands](commands.md) for the full reference. The common lifecycle:

```bash
updock up 2        # start
updock stop 2      # stop, keep data
updock restart 2   # restart
updock rebuild 2   # recreate after an image or config change
updock down 2      # remove containers, keep folder and data
updock rm 2        # delete the project entirely (asks first)
```

## Background by default

`updock up` (and the initial run) start the stack detached and hand your shell
back immediately. To watch a project live, follow its logs:

```bash
updock logs 2 -f
```

## Editing by hand

The generated `docker-compose.yml` and `.env` are normal files. Open and edit
them whenever you want. To apply changes, `updock restart <n>` or
`updock rebuild <n>`. To change values through prompts instead, run
`updock config <n>`. See [Configuration](configuration.md) and
[Projects](projects.md).
