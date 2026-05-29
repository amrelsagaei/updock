---
layout: home

hero:
  name: updock
  text: Run any Docker app from one word.
  tagline: No YAML. No flag-hunting. No copy-paste from Stack Overflow. Type a name, pick a version, answer a few prompts, and updock writes the Compose file and brings it up.
  image:
    src: /logo.svg
    alt: updock
  actions:
    - theme: brand
      text: Quick start
      link: /quickstart
    - theme: alt
      text: Installation
      link: /installation
    - theme: alt
      text: View on GitHub
      link: https://github.com/amrelsagaei/updock

features:
  - icon: '<i class="fa-solid fa-bolt"></i>'
    title: One word in, a running app out
    details: "updock postgres searches Docker Hub, ranks the matches, and walks you to a running container."
  - icon: '<i class="fa-solid fa-wand-magic-sparkles"></i>'
    title: No YAML, ever
    details: "It reads the image and asks only what matters: ports, passwords, env vars, volumes. Then it writes a clean docker-compose.yml and .env."
  - icon: '<i class="fa-solid fa-shield-halved"></i>'
    title: Strong secrets by default
    details: "Required passwords are generated with crypto/rand, stored in a 0600 .env, masked in every view, and never written into the Compose file."
  - icon: '<i class="fa-solid fa-list-ol"></i>'
    title: Control by number
    details: "updock ls numbers everything. updock up 2, updock logs 2, updock rm 2. No names to retype."
  - icon: '<i class="fa-solid fa-layer-group"></i>'
    title: Recipes for multi-service apps
    details: "WordPress with MySQL, Nextcloud with MariaDB, Gitea with Postgres, and more, scaffolded and wired together. Add your own as plain YAML."
  - icon: '<i class="fa-solid fa-box-open"></i>'
    title: Not a lock-in
    details: "Output is standard Compose. Keep the files and run docker compose yourself any time. No telemetry, ever."
---
