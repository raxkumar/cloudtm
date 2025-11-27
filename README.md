<div align="center">
  <img src="cloudtm.png" alt="CloudTimeMachine Logo" width="400"/>
  
  # CloudTimeMachine (cloudtm)

  [![Go Version](https://img.shields.io/github/go-mod/go-version/raxkumar/cloudtm)](https://github.com/raxkumar/cloudtm)
  [![Release](https://img.shields.io/github/v/release/raxkumar/cloudtm)](https://github.com/raxkumar/cloudtm/releases)
  [![License](https://img.shields.io/github/license/raxkumar/cloudtm)](LICENSE)
</div>

A lightweight Terraform wrapper CLI that automatically snapshots, versions, and manages state files for effortless rollbacks.

## ✨ Features

- 🔄 **Automatic Snapshots** - Every `apply` creates a versioned snapshot
- 📦 **Version Management** - Track all infrastructure changes
- ⏪ **Safe Rollbacks** - Restore to any previous version
- 🗂️ **Complete History** - Never lose your infrastructure state
- 🚀 **Simple & Fast** - Wraps Terraform commands seamlessly

## 📋 Prerequisites

- **Terraform** 1.0+ installed ([Download](https://developer.hashicorp.com/terraform/downloads))
- **Git** initialized in your project (optional but recommended)

## 🚀 Installation

### Homebrew (macOS/Linux)

```bash
brew tap raxkumar/cloudtm
brew install cloudtm
```

### Binary Download

Download the latest release for your platform from [GitHub Releases](https://github.com/raxkumar/cloudtm/releases).

```bash
# macOS/Linux
curl -LO https://github.com/raxkumar/cloudtm/releases/latest/download/cloudtm_$(uname -s)_$(uname -m).tar.gz
tar -xzf cloudtm_*.tar.gz
sudo mv cloudtm /usr/local/bin/
```

## 📖 Quick Start

```bash
# 1. Initialize CloudTM in your Terraform project
cloudtm init

# 2. Apply changes (automatically creates snapshots)
cloudtm apply

# 3. List all versions
cloudtm list

# 4. Rollback to a previous version
cloudtm rollback --to v2

# 5. Delete rollback when done
cloudtm rollback --del
```

## 🔧 Available Commands

| Command | Description | Flags |
|---------|-------------|-------|
| `init` | Initialize CloudTM in current project | - |
| `apply` | Apply infrastructure changes | `--auto-approve` |
| `destroy` | Destroy infrastructure resources | `--auto-approve` |
| `list` | Show all snapshot versions | - |
| `rollback` | Rollback to a version or view/delete active rollback | `--to vN`, `--del`, `--delete` |
| `version` | Show CLI version | - |

## 📚 Usage Example

```bash
# Initialize in your Terraform directory
cd my-terraform-project
cloudtm init

# Make changes to your .tf files, then apply
cloudtm apply

# View all snapshots
cloudtm list

# Output:
# 📦 CloudTimeMachine Versions
# ──────────────────────────────────────────────
# Current: v2 (Active)
# 
# Version   Timestamp              Added  Changed  Destroyed
# v2 *      2025-11-26 17:36:35Z   1      0        0
# v1        2025-11-26 15:02:22Z   5      0        0

# Destroy infrastructure before rollback
cloudtm destroy

# Rollback to version 1
cloudtm rollback --to v1
```

## 🗂️ Directory Structure

CloudTM creates a `.cloudtm/` directory in your project:

```
.cloudtm/
├── versions/          # Versioned snapshots
│   ├── v1/
│   │   └── tf_configs/
│   ├── v2/
│   └── v3/
├── meta/              # Version metadata
│   ├── v1.json
│   ├── v2.json
│   └── v3.json
├── rollback/          # Active rollback directory
├── current.json       # Current version tracker
└── rollback.json      # Rollback status
```

## 📘 Documentation

For detailed documentation, architecture, and advanced usage, see [OVERVIEW.md](OVERVIEW.md).

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📝 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- [GitHub Repository](https://github.com/raxkumar/cloudtm)
- [Issue Tracker](https://github.com/raxkumar/cloudtm/issues)
- [Releases](https://github.com/raxkumar/cloudtm/releases)

---

Made with ❤️ by [raxkumar](https://github.com/raxkumar)
