# 🚀 taskctl and Todo REST API

A modular, high-performance Task Management system written in Go. Includes a native **REST API server**, a powerful **Cobra CLI**, and a terminal-native interactive **TUI (Terminal User Interface)** built with Charm's [Bubble Tea](https://github.com/charmbracelet/bubbletea).

---

## ✨ Features

- **RESTful API backend**: Fast standard-library (`net/http`) server support matching endpoints, routing, and wildcards.
- **Dynamic task ID remapping**: Automatically re-indexes task IDs (`1..N`) upon deletion.
- **Group support**: Organize tasks into categories (e.g. `devops`, `personal`) and filter views.
- **Command Line Interface (CLI)**: Full subcommand suite powered by Cobra (`list`, `add`, `done`, `delete`, `group`).
- **Terminal User Interface (TUI)**: Rich, interactive terminal view with status toggles, group submenus, and keyboard navigation.

---

## 📁 Project structure

```text
.
├── cmd/
│   ├── server/             # REST API server entrypoint
│   │   └── main.go
│   └── taskctl/            # CLI / TUI application entrypoint
│       └── main.go
├── internal/
│   ├── client/             # HTTP API client used by CLI and TUI
│   │   └── client.go
│   ├── command/            # Cobra command definitions
│   │   ├── root.go
│   │   ├── list.go
│   │   ├── add.go
│   │   ├── done.go
│   │   ├── delete.go
│   │   ├── group.go
│   │   └── tui.go
│   ├── task/               # Shared task models
│   │   └── task.go
│   └── ui/                 # Bubble Tea TUI components and styling
│       ├── model.go
│       ├── update.go
│       └── view.go
├── go.mod
└── README.md

```

## 🛠️ Prerequisites

- **Go 1.27+** (the version declared by `go.mod`).

## 🚀 Getting started

### 1. Clone and install dependencies

```bash
# Clone repository
git clone https://github.com/bobellobo/todo-go-app.git
cd go-app

# Download dependencies
go mod download
```

## 🌐 Deploying the API server

The server is a standalone Go process backed by SQLite. It listens on port `8081`
and creates `app.db` in its current working directory. Keep that file on persistent
storage and back it up; it contains all tasks.

### Run locally

Start the server from the repository root:

```bash
go run ./cmd/server
```

The API is then available at `http://localhost:8081`.

### Build and run on a Linux server

Install Go on the server, clone the repository, and build the server binary:

```bash
git clone https://github.com/bobellobo/todo-go-app.git
cd todo-go-app
go mod download
go build -o todo-server ./cmd/server
./todo-server
```

For a persistent deployment, create a dedicated service account and run the binary
with a service manager such as `systemd`. Example service file:

```ini
# /etc/systemd/system/todo-server.service
[Unit]
Description=Todo REST API
After=network.target

[Service]
User=todo
Group=todo
WorkingDirectory=/opt/todo-go-app
ExecStart=/opt/todo-go-app/todo-server
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Install and start the service:

```bash
sudo useradd --system --home /opt/todo-go-app --shell /usr/sbin/nologin todo
sudo chown -R todo:todo /opt/todo-go-app
sudo systemctl daemon-reload
sudo systemctl enable --now todo-server
sudo systemctl status todo-server
```

The server currently binds to `:8081`. For public access, put it behind a
reverse proxy with HTTPS and expose only the proxy port (normally `443`).
Do not expose the SQLite database file or run the process as `root`. Configure
firewall rules so port `8081` is reachable only from the reverse proxy or trusted
clients. The server requires a bearer token in the `Authorization` header.

### Configure API authentication

Generate a long random token on the server and keep it out of the repository:

```bash
openssl rand -hex 32
```

Set the resulting value as `TODO_API_TOKEN` in the server's environment. The
server refuses to start if this variable is missing:

```bash
export TODO_API_TOKEN="replace-with-a-long-random-token"
./todo-server
```

For the `systemd` deployment above, add an environment file readable only by
the service account:

```bash
sudo install -o todo -g todo -m 600 /dev/null /opt/todo-go-app/server.env
sudo sh -c 'printf "TODO_API_TOKEN=%s\n" "replace-with-a-long-random-token" > /opt/todo-go-app/server.env'
```

Add this line to the `[Service]` section:

```ini
EnvironmentFile=/opt/todo-go-app/server.env
```

Then restart the service:

```bash
sudo systemctl daemon-reload
sudo systemctl restart todo-server
```

The token is validated by the Go server using a constant-time comparison.
Requests without the exact header receive `401 Unauthorized`:

```bash
curl -H "Authorization: Bearer replace-with-a-long-random-token" \
  http://127.0.0.1:8081/tasks
```

### Caddy reverse proxy

Caddy terminates HTTPS and forwards the `Authorization` header to the Go
server. A minimal Caddyfile is:

```caddyfile
api.example.com {
    reverse_proxy 127.0.0.1:8081
}
```

Keep port `8081` closed to the public internet and allow it only from the
local machine or the private network. The Go server remains responsible for
checking the token, so authentication cannot be bypassed by connecting
directly to the backend. Do not put the token in the Caddyfile or in a
publicly tracked configuration file.

Verify the deployment from the server or a trusted client:

```bash
curl -H "Authorization: Bearer replace-with-a-long-random-token" \
  http://127.0.0.1:8081/tasks
curl -H "Authorization: Bearer replace-with-a-long-random-token" \
  http://127.0.0.1:8081/groups
```

## 🖥️ Build and use the CLI (`taskctl`)

Build the CLI from the repository root:

```bash
go build -o taskctl ./cmd/taskctl
```

On Windows, build `taskctl.exe` instead:

```powershell
go build -o taskctl.exe ./cmd/taskctl
```

The CLI uses `http://localhost:8081` by default. For a remote server, configure
the base URL with the `TASKCTL_API_URL` environment variable and the bearer
token with `TASKCTL_API_TOKEN`. Set them once per
terminal session:

**PowerShell:**

```powershell
$env:TASKCTL_API_URL = "https://api.example.com"
$env:TASKCTL_API_TOKEN = "replace-with-a-long-random-token"
.\taskctl.exe list
```

**Command Prompt:**

```cmd
set TASKCTL_API_URL=https://api.example.com
set TASKCTL_API_TOKEN=replace-with-a-long-random-token
taskctl.exe list
```

**Linux/macOS:**

```bash
export TASKCTL_API_URL="https://api.example.com"
export TASKCTL_API_TOKEN="replace-with-a-long-random-token"
./taskctl list
```

To persist the value, add the corresponding `export` command to your shell
profile, or set a Windows user environment variable and open a new terminal:

```powershell
[Environment]::SetEnvironmentVariable(
  "TASKCTL_API_URL",
  "https://api.example.com",
  "User"
)
[Environment]::SetEnvironmentVariable(
  "TASKCTL_API_TOKEN",
  "replace-with-a-long-random-token",
  "User"
)
```

For a one-time override, use the `--url` flag. The flag takes precedence over
the environment variable:

```bash
./taskctl --url https://api.example.com list
```

The URL should contain only the scheme and host (and, if needed, a path prefix);
do not add `/tasks` or `/groups`, because the CLI appends those paths.
The CLI and TUI attach `Authorization: Bearer <TASKCTL_API_TOKEN>` to every
API request. Never commit the token, print it in logs, or paste it into a
public issue or shell history.

### CLI Commands

| Action | Command example |
| ----------- | ----------- |
| List all tasks | `./taskctl list` |
| List tasks in group | `./taskctl list group-name` |
| Add a new task | `./taskctl add "Buy groceries" -g personal` |
| Toggle task status | `./taskctl done 1`  ***(1 being the task ID)*** |
| Delete a task | `./taskctl delete 1` ***(Triggers dynamic remapping)*** |
| List active groups | `./taskctl group list` |
| List group tasks | `./taskctl group tasks devops` |
| Launch TUI | `./taskctl tui` |

## 🎨 Interactive TUI
The TUI uses the same API URL configuration as the CLI. After setting
`TASKCTL_API_URL`, launch the full-screen terminal interface:

```bash
./taskctl tui
```

### Keyboards shortcuts in TUI 

- `j` / `k` or `↑` / `↓`: Navigate tasks up/down
- `Space` or `Enter`: Toggle task completion status (`[ ]` ↔ `[✓]`)
- `a` or `n`: Add a new task (supports typing `-g groupname`) 
- `g`: Open interactive Group Filter Submenu 
- `d` or Backspace: Delete task 
- `r`: Refresh tasks from API 
- `q`: Exit TUI

### TUI preview

The interface is a full-screen Bubble Tea application with colored task states,
group badges, an add-task input, and an interactive group filter:

<p>
  <img src="docs/tui-normal.svg" alt="taskctl TUI task list" width="800">
</p>
<p>
  <img src="docs/tui-add-task.svg" alt="taskctl TUI add-task mode" width="800">
  <img src="docs/tui-groups.svg" alt="taskctl TUI group filter menu" width="800">
</p>

## REST API Reference

| Method | Route | Description |
| ----------- | ----------- | ----------- |
| `GET` | `/tasks` | List all tasks (Supports query param `?group=name`) |
| `POST` | `/tasks` | Create a new task |
| `PATCH` | `/tasks/{id}/toggle` | Toggle task completion status |
| `DELETE` | `/tasks/{id}` | Delete task by ID and re-index remaining IDs |
| `GET` | `/groups` | Retrieve active groups and task counts |


```bash
# Example curl request
curl -X POST http://localhost:8081/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Configure firewall", "group":"devops"}'
```
