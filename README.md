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

Edit `.env` — at minimum change these three:

```env
POSTGRES_PASSWORD=your_secure_password
SESSION_SECRET=a_random_32_character_secret_here
ADMIN_PASSWORD=your_admin_password
```

Then pull images and run:

```bash
docker compose up -d
```

Open `http://localhost:1111` — done.

> First boot auto-creates the admin account from `ADMIN_EMAIL` / `ADMIN_PASSWORD` in `.env`.

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
| `APP_URL` | ✓ | Public URL of your instance e.g. `https://yourdomain.com` |
| `NGINX_PORT` | — | Port to expose. Default `1111` |
| `POSTGRES_PASSWORD` | ✓ | PostgreSQL password |
| `SESSION_SECRET` | ✓ | Random string ≥ 32 chars |
| `ADMIN_EMAIL` | ✓ | Initial admin email |
| `ADMIN_USERNAME` | ✓ | Initial admin username |
| `ADMIN_PASSWORD` | ✓ | Initial admin password — change after first login |
| `SESSION_SECURE` | — | Set `true` when behind HTTPS |
| `SMTP_HOST` | — | Outbound email host (optional) |

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
