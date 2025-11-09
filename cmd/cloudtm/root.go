package cloudtm

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cloudtm",
	Short: "CloudTimeMachine is a tool for managing Terraform state versions and safe infrastructure rollbacks.",
	Long: `
CloudTimeMachine is a tool for managing Terraform state versions and safe infrastructure rollbacks.

Usage:
    cloudtm <command> [arguments]

The commands are:
    init         initialize cloudtm in current Terraform project
    apply        apply infrastructure changes (wrapper around Terraform apply)
    snapshot     manually create a versioned snapshot of the current Terraform state
    list         list available state snapshots and versions
    rollback     restore infrastructure to a previous snapshot
    version      print cloudtm CLI version
    help         show help for a command

Use "cloudtm help <command>" for more information about a command.
`,
	// No Run function → ensures that just typing `cloudtm` shows this help text
}

// Execute is called by main.go
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
