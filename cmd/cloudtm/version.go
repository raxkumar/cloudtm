package cloudtm

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "v0.1.0" // change this whenever you release new version

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print cloudtm CLI version",
	Long:  `Prints the current version of the CloudTimeMachine CLI tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("cloudtm CLI version %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
