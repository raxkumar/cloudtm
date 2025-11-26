package cloudtm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/raxkumar/cloudtm/helper"
	"github.com/spf13/cobra"
)

var autoApproveDestroy bool

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "destroy infrastructure (wrapper around Terraform destroy)",
	Long: `Destroys Terraform infrastructure resources.
Behaviors:
- 'cloudtm destroy' runs interactively like Terraform (requires user confirmation).
- 'cloudtm destroy --auto-approve' skips manual approval automatically.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Step 1: Ensure Terraform exists
		if _, err := exec.LookPath("terraform"); err != nil {
			fmt.Println("❌ Terraform not found in PATH.")
			fmt.Println("Please install Terraform: https://developer.hashicorp.com/terraform/downloads")
			os.Exit(1)
		}

		// Step 2: Verify CloudTimeMachine is initialized
		cwd, _ := os.Getwd()
		cloudtmDir := filepath.Join(cwd, ".cloudtm")

		if _, err := os.Stat(cloudtmDir); os.IsNotExist(err) {
			fmt.Println("❌ CloudTimeMachine not initialized. Run: cloudtm init")
			os.Exit(1)
		}

		// Step 3: Build Terraform command
		tfArgs := []string{"destroy"}
		if autoApproveDestroy {
			tfArgs = append(tfArgs, "--auto-approve")
			fmt.Println("🚀 Running 'terraform destroy --auto-approve'...")
		} else {
			fmt.Println("🚀 Running 'terraform destroy' (interactive)...")
		}

		tfCmd := exec.Command("terraform", tfArgs...)

		// Step 4: Stream output to user
		tfCmd.Stdout = os.Stdout
		tfCmd.Stderr = os.Stderr
		tfCmd.Stdin = os.Stdin

		// Step 5: Run Terraform
		if err := tfCmd.Run(); err != nil {
			fmt.Println("\n❌ Terraform destroy failed:", err)
			os.Exit(1)
		}

		fmt.Println("\n✅ Terraform destroy completed successfully.")

		// Update status to false after successful destroy
		if err := helper.SetCurrentStatus(cloudtmDir, false); err != nil {
			fmt.Println("⚠️  Warning: Failed to update current status:", err)
		}
	},
}

func init() {
	destroyCmd.Flags().BoolVar(&autoApproveDestroy, "auto-approve", false, "Skip interactive approval")
	rootCmd.AddCommand(destroyCmd)
}

