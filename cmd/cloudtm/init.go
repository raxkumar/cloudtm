package cloudtm

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize cloudtm in current Terraform project",
	Long: `Initializes CloudTimeMachine in the current Terraform project.
This command:
1. Checks for Terraform installation.
2. Creates the .cloudtm/ directory with versions/ and meta/ subfolders.
3. Runs 'terraform init' as a wrapper.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Step 1: Check if terraform is installed
		_, err := exec.LookPath("terraform")
		if err != nil {
			fmt.Println("❌ Terraform not found in PATH.")
			fmt.Println("Please install Terraform: https://developer.hashicorp.com/terraform/downloads")
			os.Exit(1)
		}

		// Step 2: Create .cloudtm/ folder structure
		cwd, _ := os.Getwd()
		cloudtmDir := filepath.Join(cwd, ".cloudtm")
		versionsDir := filepath.Join(cloudtmDir, "versions")
		metaDir := filepath.Join(cloudtmDir, "meta")
		currentFile := filepath.Join(cloudtmDir, "current.json")
		rollbackFile := filepath.Join(cloudtmDir, "rollback.json")

		if _, err := os.Stat(cloudtmDir); os.IsNotExist(err) {
			if err := os.MkdirAll(versionsDir, 0755); err != nil {
				fmt.Println("Error creating versions directory:", err)
				os.Exit(1)
			}
			if err := os.MkdirAll(metaDir, 0755); err != nil {
				fmt.Println("Error creating meta directory:", err)
				os.Exit(1)
			}
			fmt.Println("✅ Created .cloudtm/ directory with versions/ and meta/ folders.")
		} else {
			// Ensure subfolders exist even if .cloudtm already does
			os.MkdirAll(versionsDir, 0755)
			os.MkdirAll(metaDir, 0755)
			fmt.Println("ℹ️ .cloudtm/ directory already exists. Verified subfolders.")
		}

		// Create 'current.json' file to track the active snapshot version
		if _, err := os.Stat(currentFile); os.IsNotExist(err) {
			currentData := map[string]interface{}{
				"current": "",
				"status":  false,
			}
			currentJSON, _ := json.MarshalIndent(currentData, "", "  ")
			if err := os.WriteFile(currentFile, currentJSON, 0644); err != nil {
				fmt.Println("Error creating current.json file:", err)
				os.Exit(1)
			}
			fmt.Println("✅ Created 'current.json' file to track snapshot versions.")
		}

		// Create 'rollback.json' file to track rollback status
		if _, err := os.Stat(rollbackFile); os.IsNotExist(err) {
			rollbackData := map[string]string{
				"rollback": "",
			}
			rollbackJSON, _ := json.MarshalIndent(rollbackData, "", "  ")
			if err := os.WriteFile(rollbackFile, rollbackJSON, 0644); err != nil {
				fmt.Println("Error creating rollback.json file:", err)
				os.Exit(1)
			}
			fmt.Println("✅ Created 'rollback.json' file to track rollback status.")
		}

		// Step 3: Run terraform init
		fmt.Println("\n🚀 Running 'terraform init'...")

		tfCmd := exec.Command("terraform", append([]string{"init"}, args...)...)
		tfCmd.Stdout = os.Stdout
		tfCmd.Stderr = os.Stderr
		tfCmd.Stdin = os.Stdin

		err = tfCmd.Run()
		if err != nil {
			fmt.Println("\n❌ Terraform initialization failed:", err)
			os.Exit(1)
		}

		fmt.Println("\n✅ Terraform initialized successfully.")
		fmt.Println("CloudTimeMachine is now ready to manage state snapshots.")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
