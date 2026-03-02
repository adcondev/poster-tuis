# LEARNING.md — Internal Technical Summary (CV Context)

## Project Overview

**R2k Service Family Manager** is a self-contained Terminal User Interface (TUI) installer built in Go that manages the
full lifecycle of Windows POS services (Scale/Báscula and Ticket Printer) — from installation and binary deployment to
real-time status monitoring and uninstallation — all through an interactive, keyboard-driven console application.

## Tech Stack & Infrastructure

| Category              | Technologies                                                                                    |
|-----------------------|-------------------------------------------------------------------------------------------------|
| **Language**          | Go 1.25                                                                                         |
| **TUI Framework**     | Charm stack: Bubble Tea (Elm architecture), Bubbles (components), Lipgloss (styling)            |
| **Cryptography**      | `golang.org/x/crypto` — bcrypt password hashing with Base64 encoding for safe ldflags injection |
| **Build System**      | Taskfile v3 (YAML-based task runner) with modular task includes (`build`, `setup`, `ci`)        |
| **Binary Embedding**  | Go `//go:embed` directives for compile-time binary bundling (~15–20 MB self-extracting output)  |
| **OS Integration**    | Windows Service Control Manager (`sc.exe`), `taskkill.exe`, `%PROGRAMDATA%` log management      |
| **CI/CD**             | GitHub Actions (4 workflows: CI, CodeQL SAST, PR Automation, PR Status Dashboard)               |
| **Code Quality**      | golangci-lint v2 with 15+ linters (govet, errcheck, staticcheck, gosec, gocyclo, revive, etc.)  |
| **Security Scanning** | GitHub CodeQL with `security-extended` + `security-and-quality` query suites                    |
| **Version Control**   | Git with Conventional Commits, semantic PR title enforcement, auto-labeling                     |
| **Target Platform**   | Windows 10/11 (CGO_ENABLED=0 for fully static binaries, cross-compiled from any OS)             |
| **License**           | MIT License (© 2026 Adrián Constante)                                                           |

## Notable Libraries

| Library                               | Purpose                                                                                                                    |
|---------------------------------------|----------------------------------------------------------------------------------------------------------------------------|
| `charmbracelet/bubbletea`             | Elm-architecture framework powering the 6-screen state machine (Dashboard → Family → Logs → Processing → Result → Confirm) |
| `charmbracelet/bubbles`               | Pre-built TUI components: interactive lists, spinners, progress bars, contextual help overlays                             |
| `charmbracelet/lipgloss`              | Declarative terminal styling with Tokyo Night-inspired color palette and responsive layouts                                |
| `golang.org/x/crypto/bcrypt`          | Build-time password hashing utility; hashes are Base64-encoded to circumvent `$` escaping issues in Go ldflags             |
| `amannn/action-semantic-pull-request` | Enforces Conventional Commit PR titles with required scopes per component                                                  |
| `golangci/golangci-lint-action`       | Automated static analysis across 15+ linters including security (gosec), performance (prealloc), and complexity (gocyclo)  |
| `github/codeql-action`                | Weekly scheduled + per-PR SAST scanning with Go-specific security queries                                                  |

## CV-Ready Achievements

- **Architected** a modular Go application following clean architecture principles with strict package boundaries (
  `config` → `assets` → `service` → `ui`) that prevent circular dependencies and enforce separation of concerns across 4
  internal packages and 2 command entry points.

- **Engineered** a compile-time binary embedding pipeline using Go's `//go:embed` directives and Taskfile orchestration
  to bundle 4 platform-specific Windows service executables into a single ~15–20 MB self-extracting installer,
  eliminating external distribution dependencies.

- **Developed** a secure Windows Service Control Manager integration layer that programmatically manages the full
  service lifecycle (install, start, stop, restart, uninstall) via `sc.exe`, with comprehensive input validation against
  command injection, path traversal protection, and automatic failure recovery configuration.

- **Implemented** a build-time credential security pipeline using bcrypt hashing with Base64 encoding, injected via Go
  linker flags (`-ldflags -X`), ensuring plaintext passwords never appear in source code or compiled binaries while
  remaining compatible with the Go toolchain's string injection mechanism.

- **Designed** a 6-screen state machine TUI using the Elm architecture (Bubble Tea) with mutual exclusivity enforcement,
  background status polling every 5 seconds, animated progress indicators, and keyboard-driven navigation — delivering a
  zero-dependency, administrator-level installer experience.

- **Established** a comprehensive CI/CD pipeline with 4 GitHub Actions workflows spanning automated testing with race
  detection, 15+ linter static analysis, CodeQL SAST security scanning on a weekly schedule, PR size labeling, benchmark
  regression tracking with PR comment reporting, and Conventional Commit enforcement.

- **Optimized** the build process for cross-platform reproducibility by configuring `CGO_ENABLED=0` for fully static
  binaries, implementing a modular Taskfile structure (`build.yml`, `setup.yml`, `ci.yml`) with dependency-aware task
  chaining, and creating dummy asset generation for CI/linter compatibility.

## Skills Demonstrated

Go, Terminal User Interface (TUI) Development, Elm Architecture / Model-View-Update Pattern, Windows Service Management,
Build Automation (Taskfile), Binary Embedding & Self-Extracting Installers, Compile-Time Configuration Injection (
ldflags), Bcrypt Password Hashing, Input Validation & Command Injection Prevention, Path Traversal Protection, Static
Analysis (golangci-lint), SAST Security Scanning (CodeQL), GitHub Actions CI/CD, Conventional Commits, Concurrent
Programming, Cross-Compilation, Clean Architecture
