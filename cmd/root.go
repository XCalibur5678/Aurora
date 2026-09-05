package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	yesFlag     bool
	conciseFlag bool
)

var rootCmd = &cobra.Command{
	Use:     "aurora",
	Version: "0.3.0",
	Short:   "Aurora — A powerful, beginner-first package assistant for Arch Linux",
	Long: `Aurora is a package manager and AUR helper for Arch Linux that combines
the power of Paru with the simplicity of apt. It prioritizes official repositories,
explains commands in plain English, and provides transparent package lifecycle management.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if _, err := exec.LookPath("pacman"); err != nil {
			fmt.Println("Aurora currently supports Arch-based systems only. This machine does not appear to be Arch-based.")
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&yesFlag, "yes", "y", false, "skip Aurora's confirmation prompts (pacman/sudo still prompt)")
	rootCmd.PersistentFlags().BoolVarP(&conciseFlag, "concise", "c", false, "enable concise output mode (suppress teaching explanations and notes)")
}

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
