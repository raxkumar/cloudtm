# CloudTimeMachine (cloudtm) - Complete Overview

## Table of Contents

- [Introduction](#introduction)
- [Motivation](#motivation)
- [Architecture](#architecture)
- [Directory Structure](#directory-structure)
- [Installation](#installation)
- [Command Reference](#command-reference)
- [Workflow Examples](#workflow-examples)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)
- [Comparison](#comparison)

---

## Introduction

CloudTimeMachine (cloudtm) is a lightweight Terraform wrapper CLI designed to bring time-travel capabilities to your infrastructure management. It automatically creates versioned snapshots of your Terraform configurations and state files after every successful `apply`, enabling you to safely rollback to any previous state of your infrastructure.

### Why CloudTimeMachine?

Managing infrastructure with Terraform is powerful, but mistakes happen. Whether it's an accidental misconfiguration, an unexpected behavior from a provider, or simply the need to revert to a known-good state, CloudTimeMachine gives you the safety net you need.

**Key Problems CloudTimeMachine Solves:**

1. **Lost State Files** - Terraform state files can be accidentally deleted or corrupted
2. **No Native Versioning** - Terraform doesn't version your configurations automatically
3. **Risky Rollbacks** - Manually reverting infrastructure is error-prone
4. **Poor Audit Trail** - Difficult to track what changed and when
5. **Team Coordination** - Multiple team members need consistent state history

---

## Motivation

### The Problem with Vanilla Terraform

Terraform is excellent at managing infrastructure as code, but it has limitations:

1. **State Management Complexity**
   - State files must be manually backed up
   - No automatic versioning
   - Remote state requires additional setup

2. **No Built-in Rollback**
   - Reverting changes requires manual work
   - Must reconstruct old configurations
   - Risk of inconsistent state

3. **Limited History**
   - `terraform.tfstate.backup` only keeps one previous version
   - No metadata about what changed
   - No timestamps or audit trail

### The CloudTimeMachine Solution

CloudTimeMachine wraps Terraform to provide:

- **Automatic Snapshots** - Every successful apply creates a versioned backup
- **Complete Configuration Copies** - Captures all `.tf` files, state, and lock files
- **Metadata Tracking** - Records timestamps and resource change counts
- **Safe Rollback Mechanism** - Enforces safety checks before restoring
- **Simple CLI** - Familiar Terraform-like commands

---

## Architecture

### How It Works

```
┌─────────────────────────────────────────────────────────────┐
│              User runs: cloudtm apply                       │
└──────────────────┬──────────────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────────┐
│  1. Run terraform apply (interactive or --auto-approve)     │
└──────────────────┬──────────────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────────┐
│  2. Parse output for resource changes                       │
│     (X added, Y changed, Z destroyed)                       │
└──────────────────┬──────────────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────────┐
│  3. If changes detected:                                    │
│     - Create versions/vN/ directory                         │
│     - Copy all .tf files, state, lock file                  │
│     - Exclude .terraform/ and .cloudtm/                     │
│     - Generate metadata JSON                                │
└──────────────────┬──────────────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────────┐
│  4. Update current.json with:                               │
│     - New version number                                    │
│     - Active status = true                                  │
└─────────────────────────────────────────────────────────────┘
```

### Snapshot Mechanism

When you run `cloudtm apply`, the following happens:

1. **Pre-Flight Checks**
   - Verify Terraform is installed
   - Check CloudTM is initialized (`.cloudtm/` exists)

2. **Execute Terraform**
   - Run `terraform apply` with or without `--auto-approve`
   - Stream output to user in real-time
   - Capture output for parsing

3. **Analyze Changes**
   - Parse Terraform output using regex
   - Extract resource counts: added, changed, destroyed
   - Only snapshot if actual changes occurred

4. **Create Snapshot**
   - Determine next version number (v1, v2, v3...)
   - Create directory: `.cloudtm/versions/vN/tf_configs/`
   - Copy project files recursively
   - Exclude: `.terraform/`, `.cloudtm/`, `*.log`, `*.tmp`, `terraform.tfstate.backup`
   - Include: all `.tf` files, `terraform.tfstate`, `.terraform.lock.hcl`

5. **Generate Metadata**
   - Create `.cloudtm/meta/vN.json` with:
     - Version number
     - UTC timestamp
     - Resource change counts
   - Update `.cloudtm/current.json` with version and active status

### File Exclusions

CloudTimeMachine intelligently excludes files that don't need versioning:

**Excluded Directories:**
- `.terraform/` - Provider binaries (large, can be re-downloaded)
- `.cloudtm/` - CloudTM's own data (prevents recursion)

**Excluded Files:**
- `terraform.tfstate.backup` - Redundant (we have the full state)
- `*.log` - Temporary log files
- `*.tmp` - Temporary files

**Included Files:**
- `*.tf` - All Terraform configuration files
- `terraform.tfstate` - Current state file
- `.terraform.lock.hcl` - Provider version lock file
- Any custom files in your project directory

---

## Directory Structure

### Complete `.cloudtm/` Layout

```
project-root/
├── main.tf
├── variables.tf
├── outputs.tf
├── .terraform/                    # Not versioned
│   └── providers/
├── .terraform.lock.hcl           # Versioned
├── terraform.tfstate             # Versioned
└── .cloudtm/                     # CloudTM directory
    ├── versions/                 # Snapshot storage
    │   ├── v1/
    │   │   └── tf_configs/
    │   │       ├── main.tf
    │   │       ├── variables.tf
    │   │       ├── terraform.tfstate
    │   │       └── .terraform.lock.hcl
    │   ├── v2/
    │   │   └── tf_configs/
    │   │       └── ...
    │   └── v3/
    │       └── tf_configs/
    │           └── ...
    ├── meta/                     # Metadata storage
    │   ├── v1.json
    │   ├── v2.json
    │   └── v3.json
    ├── rollback/                 # Active rollback (temporary)
    │   ├── main.tf
    │   ├── terraform.tfstate
    │   └── .terraform/
    ├── current.json              # Current version tracker
    └── rollback.json             # Rollback status tracker
```

### File Descriptions

#### `current.json`
Tracks the currently active version and its deployment status.

```json
{
  "current": "v3",
  "status": true
}
```

- `current`: The latest applied version
- `status`: `true` if infrastructure is deployed, `false` if destroyed

#### `rollback.json`
Tracks active rollback operations.

```json
{
  "rollback": "v2"
}
```

- `rollback`: The version currently in rollback state
- Empty string `""` means no active rollback

#### `meta/vN.json`
Metadata for each version snapshot.

```json
{
  "version": "v3",
  "timestamp": "2025-11-26T17:36:35Z",
  "resources": {
    "added": "2",
    "changed": "1",
    "destroyed": "0"
  }
}
```

---

## Installation

### Method 1: Homebrew (Recommended for macOS/Linux)

```bash
# Add the tap
brew tap raxkumar/cloudtm

# Install cloudtm
brew install cloudtm

# Verify installation
cloudtm version
```

### Method 2: Download Binary

#### macOS (Apple Silicon)
```bash
curl -LO https://github.com/raxkumar/cloudtm/releases/latest/download/cloudtm_Darwin_arm64.tar.gz
tar -xzf cloudtm_Darwin_arm64.tar.gz
sudo mv cloudtm /usr/local/bin/
```

#### macOS (Intel)
```bash
curl -LO https://github.com/raxkumar/cloudtm/releases/latest/download/cloudtm_Darwin_x86_64.tar.gz
tar -xzf cloudtm_Darwin_x86_64.tar.gz
sudo mv cloudtm /usr/local/bin/
```

#### Linux (64-bit)
```bash
curl -LO https://github.com/raxkumar/cloudtm/releases/latest/download/cloudtm_Linux_x86_64.tar.gz
tar -xzf cloudtm_Linux_x86_64.tar.gz
sudo mv cloudtm /usr/local/bin/
```

#### Windows (64-bit)
```powershell
# Download from GitHub Releases
# Extract cloudtm.exe
# Add to PATH
```

### Method 3: Build from Source

```bash
git clone https://github.com/raxkumar/cloudtm.git
cd cloudtm
go build -o cloudtm
sudo mv cloudtm /usr/local/bin/
```

---

## Command Reference

### `cloudtm init`

Initialize CloudTimeMachine in your Terraform project.

**Usage:**
```bash
cloudtm init
```

**What it does:**
1. Checks if Terraform is installed
2. Creates `.cloudtm/` directory structure:
   - `versions/` - For snapshots
   - `meta/` - For metadata
3. Creates `current.json` with empty state
4. Creates `rollback.json` with empty state
5. Runs `terraform init` as a wrapper

**Example:**
```bash
$ cd my-terraform-project
$ cloudtm init

✅ Created .cloudtm/ directory with versions/ and meta/ folders.
✅ Created 'current.json' file to track snapshot versions.
✅ Created 'rollback.json' file to track rollback status.

🚀 Running 'terraform init'...

Initializing the backend...
Initializing provider plugins...
Terraform has been successfully initialized!

✅ Terraform initialized successfully.
CloudTimeMachine is now ready to manage state snapshots.
```

---

### `cloudtm apply`

Apply infrastructure changes and create automatic snapshots.

**Usage:**
```bash
cloudtm apply                # Interactive mode
cloudtm apply --auto-approve # Skip confirmation
```

**Flags:**
- `--auto-approve` - Skip interactive approval (non-interactive mode)

**What it does:**
1. Verifies CloudTM is initialized
2. Runs `terraform apply` (with or without auto-approve)
3. Parses output for resource changes
4. If changes detected:
   - Creates new version snapshot
   - Copies all project files
   - Generates metadata
   - Updates `current.json`

**Example (Interactive):**
```bash
$ cloudtm apply

🚀 Running 'terraform apply' (interactive)...

Terraform will perform the following actions:
  # aws_instance.web will be created
  + resource "aws_instance" "web" {
      + ami           = "ami-12345678"
      + instance_type = "t2.micro"
    }

Do you want to perform these actions? yes

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

📦 Snapshot created: v1
🗂  Saved configs: /project/.cloudtm/versions/v1/tf_configs
🧾 Metadata: /project/.cloudtm/meta/v1.json
✅ Updated current version to: v1

✅ Terraform apply completed successfully.
```

**Example (Auto-approve):**
```bash
$ cloudtm apply --auto-approve

🚀 Running 'terraform apply --auto-approve'...

Apply complete! Resources: 2 added, 1 changed, 0 destroyed.

📦 Snapshot created: v2
✅ Updated current version to: v2
```

---

### `cloudtm destroy`

Destroy infrastructure resources and update status.

**Usage:**
```bash
cloudtm destroy                # Interactive mode
cloudtm destroy --auto-approve # Skip confirmation
```

**Flags:**
- `--auto-approve` - Skip interactive approval

**What it does:**
1. Runs `terraform destroy` (with or without auto-approve)
2. Updates `current.json` status to `false`
3. Keeps version number intact (for reference)

**Example:**
```bash
$ cloudtm destroy

🚀 Running 'terraform destroy' (interactive)...

Terraform will destroy all resources.
Do you really want to destroy all resources? yes

Destroy complete! Resources: 3 destroyed.

✅ Terraform destroy completed successfully.
```

---

### `cloudtm list`

Display all snapshot versions with metadata.

**Usage:**
```bash
cloudtm list
```

**What it does:**
1. Reads all metadata files from `.cloudtm/meta/`
2. Displays in formatted table
3. Shows current version status
4. Marks active version with asterisk (*)

**Example:**
```bash
$ cloudtm list

📦 CloudTimeMachine Versions
──────────────────────────────────────────────────────────────
Current: v3 (Active)

Version   Timestamp              Added  Changed  Destroyed  Status
────────  ─────────────────────  ─────  ───────  ─────────  ───────
v3 *      2025-11-26 17:36:35Z   2      1        0          Active
v2        2025-11-26 15:02:22Z   0      1        0          -
v1        2025-11-26 13:55:04Z   5      0        0          -
──────────────────────────────────────────────────────────────
Use: cloudtm rollback --to <version>
```

**Status Indicators:**
- `Active` - This version is currently deployed
- `-` - Version is in history but not active

---

### `cloudtm rollback`

Rollback to a previous version or manage active rollbacks.

**Usage:**
```bash
cloudtm rollback                # Show current rollback status
cloudtm rollback --to vN        # Rollback to version N
cloudtm rollback --del          # Delete active rollback
cloudtm rollback --delete       # Delete active rollback (alias)
```

**Flags:**
- `--to vN` - Rollback to specific version
- `--del` / `--delete` - Delete active rollback

**Prerequisites for Rollback:**
1. All resources must be destroyed first (`cloudtm destroy`)
2. No active rollback in progress (`rollback.json` must be empty)

**What it does (rollback mode):**
1. Validates prerequisites
2. Creates `rollback/` directory
3. Copies version files to `rollback/`
4. Runs `terraform init` in rollback directory
5. Runs `terraform apply --auto-approve` in rollback directory
6. Updates `rollback.json` with version

**What it does (delete mode):**
1. Checks for active rollback
2. Runs `terraform destroy --auto-approve` in `rollback/`
3. Deletes `rollback/` directory
4. Resets `rollback.json`

**Example (View Status):**
```bash
$ cloudtm rollback

🔄 Current Rollback Status
──────────────────────────────────────────────────────────────
ℹ️  No active rollback

Usage:
  cloudtm rollback --to vN        # Rollback to version
  cloudtm rollback --del          # Delete active rollback
```

**Example (Rollback to v2):**
```bash
$ cloudtm destroy  # Must destroy first

$ cloudtm rollback --to v2

🔍 Checking terraform.tfstate...
✅ Terraform state is empty
🔍 Checking rollback status...
✅ No active rollback in progress
✅ Found version 'v2'
✅ Created rollback directory
✅ Copied configs from 'v2' to rollback directory
✅ Copied metadata 'v2.json' to rollback directory

🚀 Running 'terraform init' in rollback directory...
✅ Terraform initialized successfully

🚀 Running 'terraform apply --auto-approve' in rollback directory...
Apply complete! Resources: 3 added, 0 changed, 0 destroyed.

✅ Updated rollback.json to version: v2

🎉 Rollback completed successfully!
✅ Infrastructure rolled back to version: v2
📁 Rollback configs available in: .cloudtm/rollback/
```

**Example (Delete Rollback):**
```bash
$ cloudtm rollback --del

🔍 Checking rollback status...
✅ Found active rollback: v2

🚀 Running 'terraform destroy --auto-approve' in rollback directory...
Destroy complete! Resources: 3 destroyed.

✅ Rollback resources destroyed successfully
✅ Deleted rollback directory
✅ Reset rollback.json

🎉 Rollback cleanup completed!
```

---

### `cloudtm version`

Display the CloudTimeMachine CLI version.

**Usage:**
```bash
cloudtm version
```

**Example:**
```bash
$ cloudtm version
cloudtm CLI version v0.1.0
```

---

## Workflow Examples

### Scenario 1: Setting Up a New Project

```bash
# 1. Navigate to Terraform project
cd my-aws-infrastructure

# 2. Initialize CloudTM
cloudtm init

# 3. Apply initial infrastructure
cloudtm apply
# Creates version v1

# 4. Make changes to main.tf
vim main.tf

# 5. Apply changes
cloudtm apply
# Creates version v2

# 6. View versions
cloudtm list
```

### Scenario 2: Recovering from a Bad Deploy

```bash
# You just deployed v3 and realized it's breaking production

# 1. Check current state
cloudtm list
# Shows v3 is Active

# 2. Destroy current infrastructure
cloudtm destroy

# 3. Rollback to last known good version (v2)
cloudtm rollback --to v2

# 4. Verify rollback worked
cloudtm list

# 5. When ready to move forward, clean up rollback
cloudtm rollback --del
```

### Scenario 3: Testing Changes Safely

```bash
# You want to test a major change without losing v3

# 1. Note current version
cloudtm list
# Current: v3 (Active)

# 2. Make experimental changes
vim main.tf

# 3. Apply changes
cloudtm apply
# Creates v4

# 4. Test the new infrastructure
# If it doesn't work well...

# 5. Rollback to v3
cloudtm destroy
cloudtm rollback --to v3

# 6. Clean up when satisfied
cloudtm rollback --del
```

### Scenario 4: Team Collaboration

```bash
# Team member A makes changes
cloudtm apply  # Creates v5

# Team member B needs to sync
git pull
cloudtm list   # See what changed

# Team member B makes additional changes
cloudtm apply  # Creates v6

# If team member B needs to check team member A's config
cloudtm rollback --to v5
# Review the files in .cloudtm/rollback/
cloudtm rollback --del  # Clean up
```
---

## Best Practices

### 1. Always Initialize CloudTM First

```bash
# Do this immediately in new projects
cloudtm init
```

### 2. Use Meaningful Git Commits

CloudTM versions correspond to infrastructure states. Use git commits to track what changed:

```bash
git add .
git commit -m "feat: add load balancer for web tier"
cloudtm apply  # Creates versioned snapshot
```

### 3. Test Rollbacks in Non-Production

Before using rollback in production, test the process:

```bash
# In development environment
cloudtm apply           # Deploy changes
cloudtm destroy         # Tear down
cloudtm rollback --to v1  # Test rollback
cloudtm rollback --del  # Clean up
```

### 4. Regular Cleanup

Old versions accumulate. Periodically review and archive:

```bash
cloudtm list
# Manually remove old versions from .cloudtm/versions/ if needed
```

### 5. Add .cloudtm to .gitignore

CloudTM snapshots are local backups. Don't commit them:

```bash
echo ".cloudtm/" >> .gitignore
```

### 6. Document Major Changes

Add comments in metadata by tracking in git:

```bash
# After major change
cloudtm apply
git add .
git commit -m "v3: migrated to new VPC architecture"
```

### 7. Use --auto-approve Carefully

Interactive mode is safer for production:

```bash
# Production
cloudtm apply  # Review changes

# CI/CD pipeline
cloudtm apply --auto-approve
```

### 8. Always Destroy Before Rollback

CloudTM enforces this, but understand why:

```bash
# WRONG
cloudtm rollback --to v2  # Error: resources exist

# RIGHT
cloudtm destroy
cloudtm rollback --to v2
```

---

## Troubleshooting

### Issue: "CloudTimeMachine not initialized"

**Problem:** Running commands without `cloudtm init`

**Solution:**
```bash
cloudtm init
```

### Issue: "Resources still exist in terraform.tfstate"

**Problem:** Trying to rollback without destroying first

**Solution:**
```bash
cloudtm destroy
cloudtm rollback --to v2
```

### Issue: "Rollback already in progress"

**Problem:** Trying to rollback when one is active

**Solution:**
```bash
cloudtm rollback  # Check status
cloudtm rollback --del  # Clean up existing rollback
cloudtm rollback --to v3  # Try again
```

### Issue: "No versions found"

**Problem:** No snapshots created yet

**Solution:**
```bash
cloudtm apply  # Make first snapshot
```

### Issue: Snapshot not created after apply

**Problem:** No resource changes detected

**Check:**
```bash
# CloudTM only snapshots when changes occur
# Verify changes in terraform plan
terraform plan
```

### Issue: Permission errors on init

**Problem:** Cannot create .cloudtm directory

**Solution:**
```bash
# Check directory permissions
ls -la
chmod u+w .
cloudtm init
```

---

## Comparison

### CloudTM vs Native Terraform

| Feature | CloudTM | Native Terraform |
|---------|---------|------------------|
| State versioning | ✅ Automatic | ❌ Manual |
| Configuration snapshots | ✅ Full copies | ❌ None |
| Rollback capability | ✅ Built-in | ❌ Manual |
| Metadata tracking | ✅ Timestamps + changes | ❌ None |
| Safety checks | ✅ Enforced | ❌ User responsibility |
| Learning curve | ✅ Minimal | ✅ Familiar |
| Remote state | ⚠️ Local only | ✅ S3, etc. |

### CloudTM vs Terraform Cloud

| Feature | CloudTM | Terraform Cloud |
|---------|---------|-----------------|
| Cost | ✅ Free | 💰 Paid (free tier limited) |
| Local control | ✅ Complete | ❌ Cloud-based |
| State versioning | ✅ Local snapshots | ✅ Remote versioning |
| Setup complexity | ✅ Simple | ⚠️ Moderate |
| Team features | ❌ None | ✅ Collaboration |
| Audit trail | ✅ Basic | ✅ Comprehensive |

### When to Use CloudTM

**Best for:**
- Individual developers
- Small teams
- Local development
- Quick projects
- Learning Terraform
- Cost-conscious users

**Not ideal for:**
- Large teams (consider Terraform Cloud)
- Multi-region deployments with state sharing
- Compliance requirements (remote state audit)
- Projects needing sophisticated RBAC

---

## Advanced Topics

### Custom Snapshot Locations

While CloudTM stores snapshots locally in `.cloudtm/`, you can sync them to remote storage:

```bash
# Sync to S3 (manual)
aws s3 sync .cloudtm/ s3://my-bucket/cloudtm-backups/

# Sync to Git LFS
git lfs track ".cloudtm/versions/**/*.tfstate"
```

### Integration with CI/CD

```yaml
# GitHub Actions example
- name: Apply with CloudTM
  run: |
    cloudtm apply --auto-approve
    
- name: Upload snapshots
  uses: actions/upload-artifact@v3
  with:
    name: cloudtm-snapshots
    path: .cloudtm/
```

### Scripting with CloudTM

```bash
#!/bin/bash
# Deploy with automatic rollback on failure

cloudtm apply --auto-approve
if [ $? -ne 0 ]; then
    echo "Apply failed, rolling back"
    cloudtm destroy --auto-approve
    cloudtm rollback --to v$((CURRENT_VERSION - 1))
fi
```

---

## Future Roadmap

- [ ] Remote state backend support (S3, GCS)
- [ ] Encrypted snapshots
- [ ] Snapshot compression
- [ ] Automated cleanup policies
- [ ] Diff between versions
- [ ] Import existing Terraform projects
- [ ] Web UI for browsing versions
- [ ] Terraform Cloud integration

---

## Contributing

We welcome contributions! See the main README for guidelines.

---

## License

Apache License 2.0 - See LICENSE file for details.

---

## Support

- **Issues:** [GitHub Issues](https://github.com/raxkumar/cloudtm/issues)
- **Discussions:** [GitHub Discussions](https://github.com/raxkumar/cloudtm/discussions)
- **Documentation:** [README.md](README.md)

---

**Made with ❤️ for the Terraform community**

