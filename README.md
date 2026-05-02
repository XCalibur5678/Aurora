# Aurora v0.2 🎉

**Aurora** is a beginner-first package workflow assistant for Arch-based Linux. It prefers official repositories over the AUR, explains every operation in plain English, and helps new users gradually understand Arch instead of punishing them for not already knowing it.

---

## What's New in v0.2

v0.2 is the foundation release — the first version that proves Aurora's philosophy is the right one.

### Official repositories first, always

Aurora now queries `pacman` before the AUR for every operation. When a package exists in both, Aurora recommends the official version and explains why. No other AUR helper does this as a design principle.

### New commands

- **`aurora info <package>`** — See package details, source (official or AUR), version, and description at a glance.
- **Automatic Arch detection** — Aurora refuses to run on non-Arch systems with a clear message. No cryptic errors.

### Rewritten internals

- **`search`** — Searches official repos first, offers AUR as an explicit opt-in. Results show source labels and sort by relevance.
- **`install`** — Exact-match resolution flow: tries pacman, falls back to AUR, offers fuzzy search when the package is unknown.
- **`remove`** — Resolves installed packages safely, shows the `-Rns` command with a plain-English explanation, confirms before acting.
- **`update`** — Checks AUR packages for newer versions first, then explains `pacman -Syu` flag-by-flag before running the full system upgrade.
- **Modular architecture** — All heavy lifting moved to `internal/pacman/` and `internal/aur/`. `cmd/` is a thin orchestration layer.

---

## Commands

| Command | What it does |
|---|---|
| `aurora search <query>` | Search official repos, optionally search AUR |
| `aurora info <package>` | Show detailed package information |
| `aurora install <package>` | Install from official repos (preferred) or AUR |
| `aurora remove <package>` | Remove a package safely with confirmation |
| `aurora update` | Check AUR for updates, then full system upgrade |

---

## Philosophy

- **Pacman is King** — Aurora delegates to `pacman` and `makepkg`. It orchestrates, never replaces.
- **Transparency First** — Every destructive command is previewed with a plain-English explanation before it runs.
- **No Jargon Without Explanation** — When Aurora says "`-Syu`", it explains what those mean right there.
- **Official Wins** — Official repositories are always preferred. AUR is an explicit choice, not a hidden default.

---

## Who This Is For

- New Arch users coming from Ubuntu, Debian, Mint, or Pop!_OS
- Users who want to learn Arch, not memorize flags
- Anyone who prefers guided, explained workflows over terse output

Aurora is **not** a `yay` or `paru` replacement for power users. It's a teaching tool that happens to be a capable package assistant.

---

## Installation

### 1. Download Pre-built Binary (Recommended)

Download the latest binary from the [releases page](https://github.com/abhigyan-chatterjee/aurora/releases), then:

```bash
sudo cp aurora /usr/local/bin/
```

That's it. `/usr/local/bin/` is already on your `PATH` across bash, zsh, and fish — no shell config changes needed.

### 2. Install with Go

Requires Go 1.25+:

```bash
go install github.com/abhigyan-chatterjee/aurora@latest
```

This places the binary in `$GOPATH/bin` (defaults to `~/go/bin`). Make sure that directory is on your `PATH`.

### 3. Build from Source

```bash
git clone https://github.com/abhigyan-chatterjee/aurora.git
cd aurora
go build -o aurora main.go
sudo cp aurora /usr/local/bin/
```

---

## Quick Start

```bash
aurora search neovim      # Search official repos + optional AUR
aurora info neovim        # See where it lives and what version
aurora install neovim     # Official repos if available, AUR otherwise
aurora update             # Update AUR packages + full system upgrade
aurora remove neovim      # Safe removal with confirmation
```

---

## Roadmap

- **v0.3 — Teaching mode**: `explain` command, beginner vs concise modes, Arch concept walkthroughs
- **v0.4 — Source intelligence**: `--from official|aur`, PKGBUILD inspection, source recommendations
- **v0.5 — Lifecycle**: `list` commands, orphan handling, cache cleanup
- **v0.6+ — Polish**: Tests, better errors, docs, demo flows

---

## Contributing

We value clarity. If you have an idea, open an issue first — we prioritize changes that keep Aurora transparent and predictable.
