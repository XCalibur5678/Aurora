package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aurora",
	Short: "Aurora is a CLI based AUR helper , which enables users to easily search, install and manage packages from the Arch User Repository (AUR).",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if _, err := exec.LookPath("pacman"); err != nil {
			fmt.Println("Aurora currently supports Arch-based systems only. This machine does not appear to be Arch-based.")
			os.Exit(1)
		}
		return nil
	},
}

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
