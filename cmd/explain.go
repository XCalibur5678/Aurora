package cmd

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/abhigyan-chatterjee/aurora/internal/ui"
	"github.com/spf13/cobra"
)

type conceptExplanation struct {
	Title       string
	Summary     string
	Description string
	Aliases     []string
}

var explanations = []conceptExplanation{
	{
		Title:   "-Syu (Full System Upgrade)",
		Summary: "Refreshes database indexes and upgrades all installed packages.",
		Description: `Command: sudo pacman -Syu

What each flag means:
  -S : synchronize (operation mode)
  -y : refresh package databases from official mirrors
  -u : upgrade all installed packages to their newest versions

Why it matters:
  Arch Linux is a rolling-release distribution. Performing full system upgrades
  regularly ensures system dependencies remain synchronized and safe. Partial upgrades
  (installing new software without updating the rest of the system) are discouraged.`,
		Aliases: []string{"-Syu", "syu", "system-upgrade", "upgrade"},
	},
	{
		Title:   "-Rns (Safe Package Removal)",
		Summary: "Removes a package, its unneeded dependencies, and configuration files.",
		Description: `Command: sudo pacman -Rns <package>

What each flag means:
  -R : remove the specified package
  -n : nosave (remove global configuration files)
  -s : recursive (remove dependencies that are no longer required by any other installed package)

Why it matters:
  Plain 'pacman -R' leaves orphan dependencies behind, cluttering your disk space over time.
  'pacman -Rns' cleans up both the target application and its leftover dependencies safely.`,
		Aliases: []string{"-Rns", "rns", "remove-flags"},
	},
	{
		Title:   "AUR (Arch User Repository)",
		Summary: "Community-driven repository containing build scripts (PKGBUILDs).",
		Description: `The Arch User Repository (AUR) is a community-driven repository for Arch users.

Key points:
  1. AUR does NOT host prebuilt binary packages.
  2. It hosts build scripts called PKGBUILDs.
  3. Aurora clones the repository, reads the PKGBUILD, and runs 'makepkg' to build the binary on your system.

Safety advice:
  Always prefer official repository packages when available. Official packages are prebuilt,
  tested, and maintained directly by Arch developers.`,
		Aliases: []string{"aur", "arch-user-repository"},
	},
	{
		Title:   "Official Repositories",
		Summary: "Core, Extra, and Multilib repositories maintained by Arch Linux developers.",
		Description: `Official repositories (core, extra, multilib) are the primary software sources on Arch Linux.

Key advantages:
  - Packages are prebuilt binaries (no long build times).
  - Maintained and digitally signed by Arch Linux developers.
  - Automatically checked and updated as part of system upgrades.

Aurora's Philosophy:
  Official repositories are always preferred over the AUR.`,
		Aliases: []string{"official", "official-repos", "repos"},
	},
	{
		Title:   "PKGBUILD & makepkg",
		Summary: "Build recipe and build tool used to create Arch Linux packages.",
		Description: `PKGBUILD is a text script containing metadata (version, dependencies, source URLs) and build instructions.

How makepkg works:
  1. Downloads source files from upstream developers.
  2. Verifies integrity using checksums (SHA256/MD5).
  3. Compiles the software locally.
  4. Packages the compiled files into a '.pkg.tar.zst' archive.
  5. Installs the compiled package onto your system via pacman.`,
		Aliases: []string{"pkgbuild", "makepkg"},
	},
	{
		Title:   "Partial Upgrades",
		Summary: "Installing new software without upgrading the rest of the system.",
		Description: `A partial upgrade occurs when you update package databases ('pacman -Sy') and install a package without upgrading existing installed packages ('pacman -u').

Why it causes issues:
  Shared libraries (e.g. openssl, glibc) on Arch are continually updated. If a new application expects
  a newer shared library version than what is currently installed on your machine, programs can crash.

Best practice:
  Always perform a full system update ('aurora update' or 'sudo pacman -Syu') regularly.`,
		Aliases: []string{"partial-upgrades", "partial-upgrade"},
	},
	{
		Title:   "Foreign Packages & Orphans",
		Summary: "Packages installed outside official repos or dependencies no longer needed.",
		Description: `Foreign Packages:
  Packages present on your system that do not exist in official repositories.
  These are usually installed from the AUR or installed manually.
  View them using: 'pacman -Qm'

Orphans:
  Packages that were automatically installed as dependencies for another application,
  but are no longer required because that application was uninstalled.
  Remove them using: 'aurora remove <package>' or 'pacman -Rns'.`,
		Aliases: []string{"foreign", "foreign-packages", "orphans"},
	},
}

func explain(cmd *cobra.Command, args []string) {
	reader := bufio.NewReader(os.Stdin)

	// Since DisableFlagParsing is true (so topics like -Syu and -Rns can be passed without being parsed
	// as unknown flags), we manually inspect arguments for standard flags like --help, -h, --concise, -c.
	var cleanArgs []string
	for _, arg := range args {
		switch arg {
		case "--help", "-h":
			_ = cmd.Help()
			return
		case "--concise", "-c":
			conciseFlag = true
		case "--yes", "-y":
			yesFlag = true
		default:
			cleanArgs = append(cleanArgs, arg)
		}
	}

	if len(cleanArgs) > 0 {
		query := strings.ToLower(strings.TrimSpace(cleanArgs[0]))
		for _, exp := range explanations {
			for _, alias := range exp.Aliases {
				if strings.EqualFold(alias, query) {
					printExplanation(exp)
					return
				}
			}
		}
		fmt.Printf("No exact explanation topic found for \"%s\".\n\n", cleanArgs[0])
	}

	fmt.Println("=== Aurora Teaching Mode — Concepts & Walkthroughs ===")
	fmt.Println("Select a topic to learn more about Arch Linux package management:")
	fmt.Println()

	var items []ui.SelectionItem
	for i, exp := range explanations {
		item := ui.SelectionItem{
			Index:       i + 1,
			Name:        exp.Title,
			SourceLabel: exp.Summary,
		}
		items = append(items, item)
		fmt.Printf("  [%d] %-30s - %s\n", i+1, exp.Title, exp.Summary)
	}

	_, selected := ui.PromptSelection(reader, "\nEnter topic number or name (or press Enter to exit): ", items)
	if selected == nil {
		return
	}

	for _, exp := range explanations {
		if exp.Title == selected.Name {
			printExplanation(exp)
			return
		}
	}
}

func printExplanation(exp conceptExplanation) {
	fmt.Printf("\n=======================================================\n")
	fmt.Printf(" Topic: %s\n", exp.Title)
	fmt.Printf("=======================================================\n\n")
	if conciseFlag {
		fmt.Println(exp.Summary)
	} else {
		fmt.Println(exp.Description)
	}
	fmt.Println()
}

var explainCmd = &cobra.Command{
	Use:                "explain [topic]",
	Short:              "Explains Arch Linux package concepts, flags, and workflows in plain English",
	DisableFlagParsing: true,
	Long: `Teaching mode command that provides plain-English explanations of Arch Linux concepts,
pacman command flags (-Syu, -Rns), AUR build processes, and safety best practices.

Available topics:
  - -Syu / upgrade
  - -Rns / remove-flags
  - aur / arch-user-repository
  - official / official-repos
  - pkgbuild / makepkg
  - partial-upgrades
  - foreign / orphans`,
	Run: explain,
}

func init() {
	rootCmd.AddCommand(explainCmd)
}

func listAvailableExplainTopics() []string {
	var topics []string
	for _, exp := range explanations {
		topics = append(topics, exp.Title)
	}
	sort.Strings(topics)
	return topics
}
