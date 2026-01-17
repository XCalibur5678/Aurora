---

# Aurora

**Aurora** is a simple, opinionated CLI tool for installing AUR packages on Arch Linux.

It exists for one reason: **to make AUR usage understandable, not magical.**

---

## Why Aurora exists

Arch Linux is not difficult because it is complex. It is difficult because its tools assume prior knowledge.

Most existing AUR helpers optimize for power users:

* **Short flags** instead of descriptive words.
* **Implicit behavior** that happens behind the scenes.
* **Noisy, cryptic error messages** that require a wiki search.
* **Minimal explanation** of what is actually happening to the system.

**Aurora** deliberately chooses a different path.

---

## Design Philosophy

Aurora follows a small set of strict rules to ensure transparency and ease of use.

### 1. No Surprises

Every action is explained before it happens.

* Dependencies are shown before installation.
* Package sources are clearly identified.
* Nothing installs silently.

> If aurora is about to change your system, you will know exactly **how** and **why**.

### 2. Human-Readable Commands

aurora uses words, not flags. You should not need to memorize abbreviations to manage your software.

```bash
aurora install discord
aurora remove discord
aurora search neovim
aurora update

```

### 3. Explain, Don’t Obscure

Aurora does not hide Arch internals; it translates them. If something fails, Aurora explains the likely cause in plain language and offers the full technical output only if requested.

**Errors should guide the user, not send them to a search engine.**

### 4. Safe Defaults

Aurora is opinionated by design to protect the system:

* Partial upgrades are discouraged.
* AUR packages are clearly labeled as **unofficial**.
* PKGBUILDs are summarized before execution.
* If an action is risky, Aurora will explicitly say so.

### 5. Pacman Does the Real Work

Aurora does not replace `pacman`; it orchestrates it. All system changes are ultimately handled by `pacman` and `makepkg`. Aurora exists to make those tools usable, not to reinvent them.

---

## Target Audience

## Who aurora is FOR 

* New Arch users 
* Curious Linux users 
* People who want to understand their system 

---

## Project Status

Aurora is a **learning-focused, community-driven project**.

* Expect bugs.
* Expect rough edges.
* **Expect clarity.**

---
