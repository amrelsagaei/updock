# Projects and file layout

Every app you run is a **project** with its own folder. Nothing is scattered
around your system, and nothing is hidden: the files updock writes are standard
Docker Compose files you can read, edit, and even run without updock.

## Where projects live

By default, projects are created under `~/updock/`. Change the root with
`projects_root` in the [config](configuration.md).

```
~/updock/juice-shop/
├── docker-compose.yml   # generated for you, yours to edit
├── .env                 # values and secrets, chmod 0600, gitignored
├── .gitignore           # ignores .env and data/
├── updock.json          # metadata: image, tag, ports, created date, state
└── data/                # persistent volumes for this app
```

## The files

- **`docker-compose.yml`** is a normal Compose file. Ports are wired through
  `${HOST_PORT_n}` variables and values come from `.env` via `env_file`, so
  secrets never appear here. It works with plain `docker compose` too, so you
  are never locked in.
- **`.env`** holds the values and secrets, written with `0600` permissions
  (owner read/write only) and listed in `.gitignore` so it is never committed by
  accident.
- **`.gitignore`** ignores `.env` and `data/` by default.
- **`updock.json`** is how updock remembers the project: image, chosen tag,
  optional digest, mapped ports, and timestamps. This is what powers control by
  number.
- **`data/`** keeps databases and uploads so they survive `down` and `up`.

## Control by number

`updock ls` discovers projects by scanning `~/updock/*/updock.json` and assigns
each a stable number, ordered alphabetically by folder name. Every lifecycle
command resolves that number back to a project path. There is no separate
database: the filesystem is the source of truth, and the `updock.json` files are
the index.

State (running, stopped, and so on) is **not** read from `updock.json`. It is
queried live from Docker every time you run `updock ls` or `updock status`, so
it always matches reality even if a container crashed or was changed outside
updock.

## Editing a project

Open `docker-compose.yml` or `.env` and edit them by hand whenever you like.
updock respects your edits. To apply changes:

- `updock restart <n>` for changes that need a fresh start,
- `updock rebuild <n>` after changing the image or recreating containers,
- `updock config <n>` to change values through prompts instead of editing files.

## Global state and config paths

Non-project data lives under the standard XDG paths:

| What | Path |
|---|---|
| Global config | `~/.config/updock/config.toml` |
| User recipes | `~/.local/share/updock/recipes/` |
