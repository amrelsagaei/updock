# Quick start

This walks you from an empty machine to a running app in a couple of minutes.
If you have not installed updock yet, see [Installation](installation.md).

## Prerequisites

updock drives Docker, so you need Docker and the Compose plugin installed and
the daemon running. Check everything at once:

```bash
updock doctor
```

Every failed check tells you exactly what is wrong and how to fix it. Fix
anything it flags before continuing.

## Run your first app

Pick an app by name. updock searches Docker Hub, lets you choose, asks the few
things that matter, writes the files, and starts it.

```bash
updock postgres
```

What happens:

1. updock searches Docker Hub and shows the ranked matches.
2. You pick one (option 1 is the best match at `latest`; option 2 lets you
   choose a specific version).
3. It reads the image and asks about the port, then offers to generate a strong
   password for `POSTGRES_PASSWORD`.
4. It writes a project folder under `~/updock/postgres/` and offers to start it.

Say yes, and Postgres is running in the background.

## Control it by number

Once a project exists you never type its name again. Everything is by number:

```bash
updock ls            # list projects, numbered, with live state
updock status 1      # details for project 1 (secrets masked)
updock logs 1        # show logs
updock logs 1 -f     # follow logs
updock stop 1        # stop it (keeps data)
updock up 1          # start it again
updock open 1        # open the mapped port in your browser
updock rm 1          # delete it (asks first)
```

## Try a multi-service app

Some apps need more than one container. updock ships **recipes** that scaffold
and wire them together for you. For example, WordPress with its database:

```bash
updock wordpress
```

updock recognizes the recipe, generates the database passwords, and brings up
both WordPress and MySQL as one project. See [Recipes](recipes.md) for the full
list and how to add your own.

## What got created

Have a look at what updock wrote:

```bash
ls ~/updock/postgres/
cat ~/updock/postgres/docker-compose.yml   # a normal Compose file, yours to edit
cat ~/updock/postgres/.env                 # values and secrets (chmod 0600)
```

These are standard files. You can run `docker compose` in that folder directly
if you ever want to. See [Projects and file layout](projects.md) for details.

## Next steps

- [Usage](usage.md) - how the whole flow works
- [Commands](commands.md) - every command and flag
- [Configuration](configuration.md) - change defaults like the projects folder
