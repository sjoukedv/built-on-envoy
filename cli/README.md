# Envoy Ecosystem CLI (`ee`)

Command-line tool for working with Envoy Ecosystem extensions.

## Status

🚧 **Coming Soon** - Implementation in progress

## Planned Features

### Core Commands

```bash
# Discovery
ee plugin list                    # List all available extensions

# Running Extensions
ee run --plugin <name>           # Run Envoy with extensions
ee run --plugin ./path           # Run with local extension
ee run --config envoy.yaml       # Run with custom config

# Development
ee create-plugin <name>          # Scaffold new extension
ee plugin publish <path>         # Publish extension to marketplace

# Configuration Generation
ee gen --plugin <name>           # Generate Envoy config
ee gen --plugin <name> --full    # Generate full config
```

### Implementation Details

**Planned Tech Stack:**
- Language: Go or Rust (for single-binary distribution)
- Distribution: Homebrew, direct download, curl script
- Features:
  - Auto-download of extensions
  - Envoy config generation
  - Built-in test upstream (HTTPBin-like)
  - Example request generation
  - GitHub PR automation

**Key Behaviors:**
- Zero-config by default
- Auto-generate Envoy configuration when not provided
- Include test upstream for immediate experimentation
- Support both community and local extensions
- Generate copy-pasteable test commands

## Development

To be implemented. Expected structure:

```
cli/
├── cmd/                 # Command implementations
├── pkg/
│   ├── config/         # Config generation
│   ├── plugin/         # Plugin management
│   ├── envoy/          # Envoy integration
│   └── github/         # GitHub API integration
├── go.mod              # Go dependencies
└── main.go
```

## Contributing

If you're interested in contributing to the CLI implementation, please open an issue or discussion.

## Distribution

Planned distribution methods:
- Homebrew: `brew install envoy-ecosystem`
- Direct download: Binary releases on GitHub
- Curl script: `curl -sL https://envoy-ecosystem.io/install.sh | sh`
