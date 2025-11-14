package cloudtm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/raxkumar/cloudtm/helper"
	"github.com/spf13/cobra"
)

var rollbackTo string

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "rollback infrastructure to a previous snapshot version",
	Long: `Rollback infrastructure to a previous snapshot version.
Usage:
    cloudtm rollback --to vN

Prerequisites:
1. All resources must be destroyed (terraform.tfstate resources should be empty)
2. No active rollback should be in progress (rollback.json should be empty)`,
	Run: func(cmd *cobra.Command, args []string) {
		// Step 1: Validate --to flag
		if rollbackTo == "" {
			fmt.Println("❌ Error: --to flag is required")
			fmt.Println("Usage: cloudtm rollback --to vN")
			os.Exit(1)
		}

		// Step 2: Get current working directory
		cwd, _ := os.Getwd()
		cloudtmDir := filepath.Join(cwd, ".cloudtm")

		// Step 3: Verify CloudTimeMachine is initialized
		if _, err := os.Stat(cloudtmDir); os.IsNotExist(err) {
			fmt.Println("❌ CloudTimeMachine not initialized. Run: cloudtm init")
			os.Exit(1)
		}

		// Step 4: Check if terraform.tfstate has empty resources
		fmt.Println("🔍 Checking terraform.tfstate...")
		isEmpty, err := helper.IsStateEmpty(cwd)
		if err != nil {
			fmt.Println("❌ Error reading terraform.tfstate:", err)
			os.Exit(1)
		}
		if !isEmpty {
			fmt.Println("❌ Error: Resources still exist in terraform.tfstate")
			fmt.Println("⚠️  You must destroy all resources before rollback")
			fmt.Println("💡 Run: terraform destroy")
			fmt.Println("💡 Or: cloudtm destroy (coming soon)")
			os.Exit(1)
		}
		fmt.Println("✅ Terraform state is empty")

		// Step 5: Check if rollback.json is empty
		fmt.Println("🔍 Checking rollback status...")
		isRollbackEmpty, err := helper.IsRollbackEmpty(cloudtmDir)
		if err != nil {
			fmt.Println("❌ Error reading rollback.json:", err)
			os.Exit(1)
		}
		if !isRollbackEmpty {
			existingVersion, _ := helper.GetRollbackVersion(cloudtmDir)
			fmt.Printf("❌ Error: Rollback to version '%s' is already applied\n", existingVersion)
			fmt.Println("⚠️  You must destroy the rollback first")
			fmt.Println("💡 Destroy resources in the rollback/ directory and reset rollback.json")
			os.Exit(1)
		}
		fmt.Println("✅ No active rollback in progress")

		// Step 6: Verify requested version exists
		versionPath := filepath.Join(cloudtmDir, "versions", rollbackTo)
		if _, err := os.Stat(versionPath); os.IsNotExist(err) {
			fmt.Printf("❌ Error: Version '%s' does not exist\n", rollbackTo)
			os.Exit(1)
		}
		fmt.Printf("✅ Found version '%s'\n", rollbackTo)

		// Step 7: Create rollback directory
		rollbackDir := filepath.Join(cloudtmDir, "rollback")
		if err := os.RemoveAll(rollbackDir); err != nil {
			fmt.Println("❌ Error cleaning rollback directory:", err)
			os.Exit(1)
		}
		if err := os.MkdirAll(rollbackDir, 0755); err != nil {
			fmt.Println("❌ Error creating rollback directory:", err)
			os.Exit(1)
		}
		fmt.Println("✅ Created rollback directory")

		// Step 8: Copy tf_configs from version to rollback directory
		tfConfigsSrc := filepath.Join(versionPath, "tf_configs")
		if err := helper.CopyDirectory(tfConfigsSrc, rollbackDir, []string{}, []string{}, []string{}); err != nil {
			fmt.Println("❌ Error copying configs to rollback directory:", err)
			os.Exit(1)
		}
		fmt.Printf("✅ Copied configs from '%s' to rollback directory\n", rollbackTo)

		// Step 9: Copy metadata file to rollback directory
		metaSrc := filepath.Join(cloudtmDir, "meta", rollbackTo+".json")
		metaDest := filepath.Join(rollbackDir, rollbackTo+".json")
		if err := helper.CopyFile(metaSrc, metaDest); err != nil {
			fmt.Println("⚠️  Warning: Could not copy metadata file:", err)
		} else {
			fmt.Printf("✅ Copied metadata '%s.json' to rollback directory\n", rollbackTo)
		}

		// Step 10: Run terraform init in rollback directory
		fmt.Println("\n🚀 Running 'terraform init' in rollback directory...")
		initCmd := exec.Command("terraform", "init")
		initCmd.Dir = rollbackDir
		initCmd.Stdout = os.Stdout
		initCmd.Stderr = os.Stderr
		initCmd.Stdin = os.Stdin

		if err := initCmd.Run(); err != nil {
			fmt.Println("\n❌ Terraform init failed in rollback directory:", err)
			fmt.Println("⚠️  Rollback directory preserved for investigation")
			os.Exit(1)
		}
		fmt.Println("✅ Terraform initialized successfully")

		// Step 11: Run terraform apply --auto-approve in rollback directory
		fmt.Println("\n🚀 Running 'terraform apply --auto-approve' in rollback directory...")
		applyCmd := exec.Command("terraform", "apply", "--auto-approve")
		applyCmd.Dir = rollbackDir
		applyCmd.Stdout = os.Stdout
		applyCmd.Stderr = os.Stderr
		applyCmd.Stdin = os.Stdin

		if err := applyCmd.Run(); err != nil {
			fmt.Println("\n❌ Terraform apply failed in rollback directory:", err)
			fmt.Println("⚠️  Rollback directory preserved for investigation")
			os.Exit(1)
		}

		// Step 12: Update rollback.json
		if err := helper.UpdateRollbackVersion(cloudtmDir, rollbackTo); err != nil {
			fmt.Println("⚠️  Warning: Failed to update rollback.json:", err)
		} else {
			fmt.Printf("\n✅ Updated rollback.json to version: %s\n", rollbackTo)
		}

		fmt.Println("\n🎉 Rollback completed successfully!")
		fmt.Printf("✅ Infrastructure rolled back to version: %s\n", rollbackTo)
		fmt.Println("📁 Rollback configs available in: .cloudtm/rollback/")
	},
}

func init() {
	rollbackCmd.Flags().StringVar(&rollbackTo, "to", "", "Version to rollback to (e.g., v1, v2)")
	rollbackCmd.MarkFlagRequired("to")
	rootCmd.AddCommand(rollbackCmd)
}
