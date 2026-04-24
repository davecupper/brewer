# brewer

A lightweight CLI for managing local service dependencies and startup order during development.

---

## Installation

```bash
go install github.com/yourusername/brewer@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/brewer.git && cd brewer && go build -o brewer .
```

---

## Usage

Define your services in a `brewer.yaml` file at the root of your project:

```yaml
services:
  postgres:
    cmd: "docker run -p 5432:5432 postgres"
  redis:
    cmd: "docker run -p 6379:6379 redis"
  api:
    cmd: "go run ./cmd/api"
    depends_on:
      - postgres
      - redis
```

Then start everything in the correct order:

```bash
brewer up
```

Stop all running services:

```bash
brewer down
```

Check service status:

```bash
brewer status
```

Brewer resolves dependency order automatically and streams logs from each service to your terminal.

---

## Commands

| Command         | Description                          |
|-----------------|--------------------------------------|
| `brewer up`     | Start all services in dependency order |
| `brewer down`   | Stop all running services            |
| `brewer status` | Show current status of each service  |
| `brewer logs`   | Tail logs from a specific service    |

---

## License

MIT © 2024 yourusername