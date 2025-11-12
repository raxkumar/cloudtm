package cloudtm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"

	"github.com/raxkumar/cloudtm/helper"
	"github.com/spf13/cobra"
)

var autoApprove bool

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "apply infrastructure changes (wrapper around Terraform apply)",
	Long: `Applies Terraform infrastructure changes and snapshots state if any change occurs.
Behaviors:
- 'cloudtm apply' runs interactively like Terraform.
- 'cloudtm apply --auto-approve' skips manual approval automatically.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Step 1: Ensure Terraform exists
		if _, err := exec.LookPath("terraform"); err != nil {
			fmt.Println("❌ Terraform not found in PATH.")
			fmt.Println("Please install Terraform: https://developer.hashicorp.com/terraform/downloads")
			os.Exit(1)
		}

		// Step 2: Verify CloudTimeMachine directories
		cwd, _ := os.Getwd()
		cloudtmDir := filepath.Join(cwd, ".cloudtm")
		versionDir := filepath.Join(cloudtmDir, "versions")
		metaDir := filepath.Join(cloudtmDir, "meta")

		if _, err := os.Stat(cloudtmDir); os.IsNotExist(err) {
			fmt.Println("❌ CloudTimeMachine not initialized. Run: cloudtm init")
			os.Exit(1)
		}
		os.MkdirAll(versionDir, 0755)
		os.MkdirAll(metaDir, 0755)

		// Step 3: Build Terraform command
		tfArgs := []string{"apply"}
		if autoApprove {
			tfArgs = append(tfArgs, "--auto-approve")
			fmt.Println("🚀 Running 'terraform apply --auto-approve'...")
		} else {
			fmt.Println("🚀 Running 'terraform apply' (interactive)...")
		}

		tfCmd := exec.Command("terraform", tfArgs...)

		// Step 4: Stream + capture output simultaneously
		var outputBuf bytes.Buffer
		mw := io.MultiWriter(os.Stdout, &outputBuf)
		tfCmd.Stdout = mw
		tfCmd.Stderr = mw
		tfCmd.Stdin = os.Stdin

		// Step 5: Run Terraform
		if err := tfCmd.Run(); err != nil {
			fmt.Println("❌ Terraform apply failed:", err)
			os.Exit(1)
		}

		// Step 6: Analyze output for changes
		output := outputBuf.String()
		re := regexp.MustCompile(`Resources: (\d+) added, (\d+) changed, (\d+) destroyed`)
		matches := re.FindStringSubmatch(output)

		if len(matches) == 4 {
			added := matches[1]
			changed := matches[2]
			destroyed := matches[3]

			if added != "0" || changed != "0" || destroyed != "0" {
				// Step 7: Snapshot logic
				files, _ := os.ReadDir(versionDir)
				nextVersion := fmt.Sprintf("v%d", len(files)+1)

				// Create version directory structure
				versionPath := filepath.Join(versionDir, nextVersion)
				tfConfigsPath := filepath.Join(versionPath, "tf_configs")
				if err := os.MkdirAll(tfConfigsPath, 0755); err != nil {
					fmt.Println("⚠️ Failed to create version directory:", err)
					return
				}

				// Copy entire project directory excluding .terraform, .cloudtm, and unnecessary files
				excludeDirs := []string{".terraform", ".cloudtm"}
				excludeFiles := []string{"terraform.tfstate.backup"}
				excludePatterns := []string{"*.log", "*.tmp"}
				if err := helper.CopyDirectory(cwd, tfConfigsPath, excludeDirs, excludeFiles, excludePatterns); err != nil {
					fmt.Println("⚠️ Failed to copy project files:", err)
					return
				}

				// Create metadata JSON
				metaDest := filepath.Join(metaDir, nextVersion+".json")
				meta := map[string]interface{}{
					"version":   nextVersion,
					"timestamp": time.Now().UTC().Format(time.RFC3339),
					"resources": map[string]string{
						"added":     added,
						"changed":   changed,
						"destroyed": destroyed,
					},
				}

				metaJSON, _ := json.MarshalIndent(meta, "", "  ")
				if err := os.WriteFile(metaDest, metaJSON, 0644); err != nil {
					fmt.Println("⚠️ Failed to write metadata file:", err)
					return
				}

				// Update current.json
				if err := helper.UpdateCurrentVersion(cloudtmDir, nextVersion); err != nil {
					fmt.Println("⚠️ Failed to update current.json:", err)
					return
				}

				fmt.Printf("\n📦 Snapshot created: %s\n", nextVersion)
				fmt.Printf("🗂  Saved configs: %s\n", tfConfigsPath)
				fmt.Printf("🧾 Metadata: %s\n", metaDest)
				fmt.Printf("✅ Updated current version to: %s\n", nextVersion)
			} else {
				fmt.Println("✅ No resource changes detected — skipping snapshot.")
			}
		} else {
			fmt.Println("⚠️ Could not parse Terraform output for resource changes.")
		}

		fmt.Println("\n✅ Terraform apply completed successfully.")
	},
}

func init() {
	applyCmd.Flags().BoolVar(&autoApprove, "auto-approve", false, "Skip interactive approval")
	rootCmd.AddCommand(applyCmd)
}
