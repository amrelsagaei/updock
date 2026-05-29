# Recipes

Some apps are a single container. updock infers everything it needs from the
image. Other apps need a friend: WordPress needs a database, Nextcloud needs
MariaDB. You cannot guess that from the WordPress image alone.

**Recipes** are small curated definitions for these multi-service apps. When you
run `updock wordpress`, updock recognizes the recipe, scaffolds both services
into one project, wires them together, and asks for shared values (like the
database password) once.

Recipes are data, not code, so they are easy to add and contribute.

## Using a recipe

Just run updock with the recipe name:

```bash
updock wordpress
```

If a recipe matches the name, updock uses the recipe flow. If not, it falls back
to the single-image flow and tells you which path it took.

## Built-in recipes

| Recipe | Services |
|---|---|
| `wordpress` | WordPress + MySQL |
| `ghost` | Ghost + MySQL |
| `nextcloud` | Nextcloud + MariaDB |
| `gitea` | Gitea + PostgreSQL |
| `n8n` | n8n + PostgreSQL |
| `plausible` | Plausible + PostgreSQL + ClickHouse |

These are embedded in the binary, so they work offline.

## Recipe format

A recipe is a YAML file with three blocks: `meta`, `prompts`, and `services`.
Shared values are defined once in `prompts` and referenced as `${VAR}` in any
service. Here is the shape (trimmed):

```yaml
meta:
  name: wordpress
  description: WordPress with MySQL database
  version: "1.0"
  tags: [cms, blog, php]

prompts:
  - name: DB_PASSWORD
    description: Database password
    required: true
    secret: true
    generate: password      # updock generates a strong value
  - name: WP_PORT
    description: WordPress HTTP port
    default: "8080"
    type: port

services:
  wordpress:
    image: wordpress
    default_tag: latest
    ports:
      - "${WP_PORT}:80"
    environment:
      WORDPRESS_DB_PASSWORD: "${DB_PASSWORD}"
    depends_on:
      - db
  db:
    image: mysql
    default_tag: "8.0"
    environment:
      MYSQL_ROOT_PASSWORD: "${DB_PASSWORD}"
```

Notice `DB_PASSWORD` appears in both services. It is asked once and substituted
everywhere. Prompts marked `secret: true` are masked in output and written only
to `.env`; with `generate: password` updock creates a strong random value for
you.

## Adding your own recipe

You can drop your own recipes in:

```
~/.local/share/updock/recipes/<name>.yaml
```

(or `$XDG_DATA_HOME/updock/recipes/` if that variable is set). A user recipe
with the same `meta.name` as a built-in one overrides it. updock validates every
recipe on load: required fields must be present, every service needs an image,
and every `${VAR}` reference must resolve to a defined prompt.

## Contributing a recipe

To add a recipe to the built-in catalog for everyone, contribute it upstream.
See [CONTRIBUTING.md](https://github.com/amrelsagaei/updock/blob/main/CONTRIBUTING.md) for the recipe contribution flow.
