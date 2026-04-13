# Automatic 3rd party KiCad Library Repository

![GitHub Last Commit](https://img.shields.io/github/last-commit/CDFER/cd_fer-kicad-repository)
![License](https://img.shields.io/badge/License-MIT-blue)
![KiCad v8.0+](https://img.shields.io/badge/KiCad-8.0%2B-red)
![KiCad v9.0+](https://img.shields.io/badge/KiCad-9.0%2B-red)

Automatically maintained 3rd Party KiCad library repository, updated daily with the latest releases.

## Features

- **Daily Updates**: Automated GitHub Actions workflow keeps library metadata current
- **Version Tracking**: Clear release history with timestamped updates
- **One-Click Installation**: Simple repository integration in KiCad

## Installation

1. Open KiCad
2. Go to **Plugin and Content Manager** → **Manage**
3. Add this repository URL:

   ```
   https://raw.githubusercontent.com/CDFER/cd_fer-kicad-repository/main/repository.json
   ```

KiCad will automatically alert your when to update your library.

## Automated Maintenance

This repository uses a GitHub Actions workflow that:

- 🕒 Runs daily at 06:15 UTC
- 🔄 Checks for new library releases
- 📦 Updates package metadata automatically
- 📅 Maintains accurate timestamps in repository.json
- ✅ Verifies JSON structure integrity

## Repository Structure

```
.
├── .github/
│   └── workflows/
│       └── update-repo.yml     # GitHub Actions workflow
├── update_metadata.py          # Maintenance script
├── packages.json               # Package version metadata
├── repository.json             # KiCad repository configuration
└── README.md                   # This documentation
```

## Contributing

While this repository is primarily automated, we welcome:

- Bug reports for outdated components
- Suggestions for improved automation
- Documentation improvements

Before contributing:

1. Open an issue to discuss proposed changes
2. Ensure updates follow KiCad library standards
3. Maintain backward compatibility with KiCad 8.0+

## License

This repository is licensed under [MIT License](LICENSE) - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- KiCad development team for the awesome EDA tool
- GitHub for providing the Actions automation platform
