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

- **Go 1.22+** (Go 1.22 or higher is required for standard library path wildcards).

## 🚀 Getting started

### 1. Clone and install dependencies

```bash
# Clone repository
git clone [https://github.com/bobellobo/todo-go-app.git](https://github.com/bobellobo/todo-go-app.git)
cd go-app

# Download dependencies
go mod download
```

### 2. Start the API server

Launch the REST server in a terminal window: 
```bash
go run ./cmd/server
```

Server now runs locally on `[buah](http://localhost:8080).` You can obsviously chose to run it anywhere. I hava personnaly chosen to deploy it on my VPS in order to add, edit and view tasks anywhere, anytime.

## 🖥️ CLI Usage (`taskctl`)

You can run `taskctl` directly via go run or compile the executable binary:

```bash
# Build binary
go build -o taskctl ./cmd/taskctl

# Optional: Install globally to system PATH
go install ./cmd/taskctl
```

### CLI Commands

| Action | Comand example |
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
Launch the full-screen terminal interface:

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
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Configure firewall", "group":"devops"}'
```
