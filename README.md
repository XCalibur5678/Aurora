# Aurora

**Aurora** is a simple, opinionated CLI tool for installing AUR packages on Arch Linux or Arch Linux based distros like Manjaro , Garuda Linux etc.

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

> If Aurora is about to change your system, you will know exactly **how** and **why**.

### 2. Human-Readable Commands

Aurora uses words, not flags. You should not need to memorize abbreviations to manage your software.

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

Aurora does not replace `pacman`; it orchestrates it. All system changes are ultimately handled by `pacman` and `makepkg`. Aurora exists to make those tools user friendly, not to reinvent them.

---

## Target Audience

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
Good instinct. A **CONTRIBUTING section is a quiet flex** — most student projects skip it, and interviewers notice.

I’ll give you a **contributing section that matches aurora’s philosophy**: strict, welcoming, and not naïve.

You want three things simultaneously:

1. Encourage contributions
2. Filter low-effort noise
3. Set technical standards without sounding arrogant

Here’s a **ready-to-paste** section. Don’t edit it unless you’re very sure.

---

## Contributing

Aurora is an opinionated project.  
Contributions are welcome, but not everything will be accepted.

Before contributing, please read this section carefully.

---

### What kind of contributions are welcome

- Bug fixes
- Improvements to error explanations
- Better human-readable output
- Documentation improvements
- Small, focused features that align with the project philosophy

If you are unsure whether something fits, open an issue first.

---

### What will likely not be accepted

- Large feature additions without prior discussion
- Flag-heavy interfaces or short-option aliases
- Changes that hide or automate potentially risky behavior
- Anything that reduces clarity in favor of convenience

Aurora intentionally favors understanding over speed and automation.

---

### Code style and structure

- Keep changes small and focused
- Follow existing project structure
- Avoid clever abstractions
- Prefer explicit, readable code over brevity

If your change is hard to explain, it probably does not belong.

---

### Commit messages

Write clear, descriptive commit messages.

Good:
```

explain missing dependency errors more clearly

```

Bad:
```

fix stuff

```

---

### Testing and behavior changes

If your change alters user-visible behavior:
- Explain the reasoning in the pull request
- Include example terminal output when possible

Aurora values predictable behavior.

---

### Reporting bugs

When reporting bugs, include:
- Command used
- Expected behavior
- Actual behavior
- Relevant error output (redacted if necessary)

Vague bug reports are hard to act on.

---

### Code of conduct

Be respectful and constructive.

This project is maintained in spare time.
Disagreements are fine.
Disrespect is not.

---

