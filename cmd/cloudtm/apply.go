package cloudtm

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var autoApprove bool

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "apply infrastructure changes (wrapper around Terraform apply)",
	Long: `Applies Terraform infrastructure changes and ensures CloudTimeMachine is initialized.
Behaviors:
- 'cloudtm apply' runs interactively like Terraform (auto-confirms with 'yes').
- 'cloudtm apply -auto-approve' passes flag directly to Terraform.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Step 1: Check if Terraform is installed
		_, err := exec.LookPath("terraform")
		if err != nil {
			fmt.Println("❌ Terraform not found in PATH.")
			fmt.Println("Please install Terraform: https://developer.hashicorp.com/terraform/downloads")
			os.Exit(1)
		}

		// Step 2: Ensure .cloudtm exists
		cwd, _ := os.Getwd()
		cloudtmDir := filepath.Join(cwd, ".cloudtm")

		if _, err := os.Stat(cloudtmDir); os.IsNotExist(err) {
			fmt.Println("╷")
			fmt.Println("│ Error: CloudTimeMachine not initialized in this directory.")
			fmt.Println("│")
			fmt.Println("│ To initialize CloudTimeMachine, run:")
			fmt.Println("│   cloudtm init")
			fmt.Println("╵")
			os.Exit(1)
		}

		// Step 3: Detect if user passed -auto-approve
		autoApprove := false
		for _, arg := range args {
			if strings.Contains(arg, "-auto-approve") {
				autoApprove = true
				break
			}
		}

		// Step 4: Prepare Terraform apply command
		tfArgs := append([]string{"apply"}, args...)
		tfCmd := exec.Command("terraform", tfArgs...)

		// Capture output (directly stream to console)
		tfCmd.Stdout = os.Stdout
		tfCmd.Stderr = os.Stderr

		// Step 5: Handle input stream (interactive vs auto)
		stdinPipe, err := tfCmd.StdinPipe()
		if err != nil {
			fmt.Println("Error creating stdin pipe:", err)
			os.Exit(1)
		}

		if autoApprove {
			fmt.Println("🚀 Running 'terraform apply -auto-approve'...")
		} else {
			fmt.Println("🚀 Running 'terraform apply' (interactive)...")
		}

		// Start process
		if err := tfCmd.Start(); err != nil {
			fmt.Println("Error starting terraform apply:", err)
			os.Exit(1)
		}

		// If not auto-approve, send "yes\n" when prompted
		if !autoApprove {
			// Wait a bit for the prompt to appear
			go func() {
				scanner := bufio.NewScanner(os.Stdin)
				for scanner.Scan() {
					input := scanner.Text()
					if strings.ToLower(strings.TrimSpace(input)) == "yes" {
						stdinPipe.Write([]byte("yes\n"))
					} else {
						stdinPipe.Write([]byte(input + "\n"))
					}
				}
			}()
		}

		// Wait for Terraform to finish
		if err := tfCmd.Wait(); err != nil {
			fmt.Println("\n❌ Terraform apply failed:", err)
			os.Exit(1)
		}

		fmt.Println("\n✅ Terraform apply completed successfully.")
		fmt.Println("CloudTimeMachine is ready to snapshot the updated state (feature coming soon).")
	},
}

func init() {
	// Define a flag for the command
    applyCmd.Flags().BoolVar(&autoApprove, "auto-approve", false, "Skip interactive approval")
	rootCmd.AddCommand(applyCmd)
}
