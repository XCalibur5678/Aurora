# Aurora

**Aurora** is a transparent, opinionated CLI tool for managing AUR packages on Arch Linux and its derivatives (Manjaro, Garuda, etc.).

It exists for one reason: **to make AUR usage understandable, not magical.**

---

## The Philosophy: "No Surprises"

Arch Linux is powerful because it is explicit. Most AUR helpers hide this power behind cryptic flags and implicit, automated behavior. Aurora takes a different path.

*   **Transparency First:** Aurora is an orchestrator, not a black box. It wraps standard tools like `git`, `makepkg`, and `pacman` to perform tasks, meaning you can always see exactly what is happening to your system.
*   **Human-Readable:** No need to memorize obscure flag combinations. Aurora uses clear, descriptive commands.
*   **Pacman is King:** Aurora does not reinvent the wheel. It delegates the heavy lifting to `pacman` and `makepkg`, acting as a clean, intuitive layer on top.

---

## Features

Aurora currently provides the essential lifecycle management for AUR packages:

- `aurora search <pkg>`: Interactive search via the AUR RPC, allowing you to choose and inspect packages before installation.
- `aurora install <pkg>`: Clones the AUR repository, lets you inspect it, and builds/installs using `makepkg -si`.
- `aurora remove <pkg>`: Safely removes packages while helping you resolve naming discrepancies.
- `aurora update`: Checks your system for foreign (AUR) packages, compares versions against the AUR, and optionally triggers a rebuild for updates, followed by a system-wide `pacman -Syu`.

---

## Installation

### 1. The Simple Way (Recommended)
If you have [Go](https://go.dev/) installed, you can install the latest version directly:

```bash
go install github.com/abhigyan-chatterjee/aurora@latest
```

### 2. Build from Source
For a local build from source:

```bash
git clone https://github.com/abhigyan-chatterjee/aurora.git
cd aurora
go build -o aurora main.go
sudo cp aurora /usr/local/bin/
```

### 3. Using make
Using `make` to build and install:

```bash
git clone https://github.com/abhigyan-chatterjee/Aurora.git
cd aurora
make
sudo make install
```

---

## Getting Started

Getting started is as simple as:

```bash
aurora search <package-name>
# Select your package, then follow the prompts.
```

---

## Why use this?

Aurora is built for users who are transitioning to arch from Debian based Distros. The goal is to provide a familiar and simple flow for them to be comfortable with their system.  

It is a learning-focused, community-driven project. Expect clarity, simplicity, and a direct line to your system's operations.

---

## Contributing

We value clarity above all else. If you have an idea for a feature or a bug fix, please open an issue to discuss it first. We prioritize changes that keep the tool transparent, readable, and predictable.
