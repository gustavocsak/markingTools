# 🧰 Marking Tools (Go Edition)

A collection of command-line tools designed to assist with marking programming assignments.

I am migrating all my previous marking utilities to **Go** to:

- Learn the Go programming language  
- Improve performance, portability, and simplicity  
- Streamline my marking workflow across platforms  

---

## 📍 Overview

This repository contains small, focused utilities that help automate repetitive marking tasks for coding assignments.  
Each tool lives inside the `cmd/` directory as its own executable.

Examples of tasks these tools aim to solve:

- Flattening student submission folders  
- Cleaning and normalizing file structures  
- Automating common marking steps  
- Reducing time spent on repetitive review tasks  

More tools will be added as they are migrated or re-written in Go.

---

## Getting Started

### Prerequisites

- Go 1.22+ installed  
- `$GOPATH/bin` added to your system `PATH` (required for `go install`)  

Check your Go setup:
```bash
go version
go env GOPATH
```

### Install a Tool

- Install a tool globally (so it can be run from anywhere):
```bash
go install ./cmd/<tool_name>
```


## Roadmap

- Migrate more existing marking scripts to Go

- Add automated testing for tools

- Add documentation for each tool in cmd/<tool>/README.md

- Package common shared logic into internal/
