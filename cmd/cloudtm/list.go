package cloudtm

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/raxkumar/cloudtm/helper"
	"github.com/spf13/cobra"
)

type VersionMetadata struct {
	Version   string
	Timestamp string
	Added     string
	Changed   string
	Destroyed string
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list available state snapshots and versions",
	Long: `Lists all available CloudTimeMachine snapshot versions with metadata.
Shows version number, timestamp, and resource change statistics for each snapshot.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Step 1: Check CloudTM is initialized
		cwd, _ := os.Getwd()
		cloudtmDir := filepath.Join(cwd, ".cloudtm")
		metaDir := filepath.Join(cloudtmDir, "meta")

		if _, err := os.Stat(cloudtmDir); os.IsNotExist(err) {
			fmt.Println("❌ CloudTimeMachine not initialized. Run: cloudtm init")
			os.Exit(1)
		}

		// Step 2: Get current version and status
		currentVersion, currentStatus, err := helper.GetCurrentVersion(cloudtmDir)
		if err != nil {
			fmt.Println("⚠️  Warning: Could not read current.json:", err)
			currentVersion = ""
			currentStatus = false
		}

		// Step 3: Read all metadata files
		files, err := os.ReadDir(metaDir)
		if err != nil {
			fmt.Println("❌ Error reading meta directory:", err)
			os.Exit(1)
		}

		// Filter and collect version metadata
		var versions []VersionMetadata
		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
				continue
			}

			metaPath := filepath.Join(metaDir, file.Name())
			data, err := os.ReadFile(metaPath)
			if err != nil {
				continue
			}

			var meta map[string]interface{}
			if err := json.Unmarshal(data, &meta); err != nil {
				continue
			}

			resources, _ := meta["resources"].(map[string]interface{})
			version := VersionMetadata{
				Version:   meta["version"].(string),
				Timestamp: meta["timestamp"].(string),
				Added:     resources["added"].(string),
				Changed:   resources["changed"].(string),
				Destroyed: resources["destroyed"].(string),
			}
			versions = append(versions, version)
		}

		// Step 4: Check if any versions exist
		if len(versions) == 0 {
			fmt.Println("ℹ️  No versions found. Run 'cloudtm apply' to create your first snapshot.")
			return
		}

		// Step 5: Sort versions (v1, v2, v3... in descending order for display)
		sort.Slice(versions, func(i, j int) bool {
			numI := extractVersionNumber(versions[i].Version)
			numJ := extractVersionNumber(versions[j].Version)
			return numI > numJ
		})

		// Step 6: Display header
		fmt.Println("\n📦 CloudTimeMachine Versions")
		fmt.Println("──────────────────────────────────────────────────────────────")

		// Show current version status
		if currentVersion != "" {
			statusText := "Inactive"
			if currentStatus {
				statusText = "Active"
			}
			fmt.Printf("Current: %s (%s)\n\n", currentVersion, statusText)
		} else {
			fmt.Println("Current: None\n")
		}

		// Step 7: Create table with tabwriter
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Version\tTimestamp\tAdded\tChanged\tDestroyed\tStatus")
		fmt.Fprintln(w, "────────\t─────────────────────\t─────\t───────\t─────────\t───────")

		// Step 8: Print each version
		for _, v := range versions {
			// Mark current version with asterisk
			versionDisplay := v.Version
			if v.Version == currentVersion {
				versionDisplay = v.Version + " *"
			}

			// Determine status
			status := "-"
			if v.Version == currentVersion && currentStatus {
				status = "Active"
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
				versionDisplay,
				v.Timestamp,
				v.Added,
				v.Changed,
				v.Destroyed,
				status)
		}

		w.Flush()

		// Step 9: Display footer
		fmt.Println("──────────────────────────────────────────────────────────────")
		fmt.Println("Use: cloudtm rollback --to <version>")
		fmt.Println()
	},
}

// extractVersionNumber extracts numeric part from version string (e.g., "v10" -> 10)
func extractVersionNumber(version string) int {
	numStr := strings.TrimPrefix(version, "v")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0
	}
	return num
}

func init() {
	rootCmd.AddCommand(listCmd)
}

