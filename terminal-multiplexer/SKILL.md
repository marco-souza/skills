---
name: tmux
description: Terminal Multiplexer allowing to run multiple terminal sessions in a single window - can be used to run multiple processes in the background
---

# Tmux Instructions

Tmux is a terminal multiplexer that allows you to run multiple terminal sessions
in a single window. It can be used to run multiple processes in the background.

It's widely available on most Linux distributions and can be installed on macOS
using Homebrew. The host machine already have it installed

## Managing sessions

### Creating a new session

```bash
tmux new-session -s <session_name>
```

### Creating a detached session

```bash
tmux new-session -d -s <session_name>
```

### Creating a session with a specific window

```bash
tmux new-session -s <session_name> -n <window_name>
```

### Creating a session with a specific window and a specific layout

```bash
tmux new-session -s <session_name> -n <window_name> -d
```

## Controlling the session

### Sending commands to a session

```bash
tmux send-keys "<command>" C-m
```

### Capturing output

```bash
tmux capture-pane -p
```
