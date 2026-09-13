# ngumpul-host

A quiet, self-hostable community portal for side projects and experiments.

Stack: **Go** · **SvelteKit** · **PostgreSQL** · **Nginx** · **Docker**

---

## Quick Install

```bash
git clone https://github.com/Ijon6k/ngumpul-host.git
cd ngumpul-host
cp .env.example .env
```

Edit `.env` — at minimum change these two:

```env
POSTGRES_PASSWORD=your_secure_password
SESSION_SECRET=a_random_32_character_secret_here
```

Then pull images and run:

```bash
docker compose up -d
```

Open `http://localhost:1111` — on first boot you'll be guided through the
`/setup` page to configure the canonical node domain and create the primary
administrator account.

---

## Compose Files

| File | Purpose |
|---|---|
| `docker-compose.yml` | **Production** — pulls pre-built images from GHCR. Use this for self-hosting. |
| `docker-compose.dev.yml` | **Development** — builds from source. Use this when contributing. |

### Production (pre-built images)

```bash
docker compose up -d
```

### Development (build from source)

```bash
docker compose -f docker-compose.dev.yml up --build
```

---

## Environment Variables

Copy `.env.example` to `.env` and set:

| Variable | Required | Description |
|---|---|---|
| `NGINX_PORT` | — | Port to expose. Default `1111` |
| `POSTGRES_PASSWORD` | ✓ | PostgreSQL password |
| `SESSION_SECRET` | ✓ | Random string ≥ 32 chars |
| `SESSION_SECURE` | — | Set `true` when behind HTTPS |
| `SMTP_HOST` | — | Outbound email host (optional) |

> The canonical node domain and the first administrator account are configured
> through the `/setup` page (stored in the `instance_settings` table), not via
> environment variables.

---

## Architecture

```
Internet
    │
    ▼
Nginx :1111
    ├─ /api/v1/*  ──▶  Go backend  :8080
    └─ /*          ──▶  SvelteKit   :3000
                              │
                              ▼
                        PostgreSQL  :5432
```

Projects on this server are provisioned manually by an administrator — there is no automated container spawning.

---

## License

[MIT](LICENSE)
