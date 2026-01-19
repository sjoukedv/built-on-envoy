# Envoy Ecosystem CLI (`ec`)

Command-line tool for working with Envoy Ecosystem extensions.

## Status

🚧 **Coming Soon** - Implementation in progress

## Planned Features

### Core Commands

```bash
# Discovery
ec plugin list                    # List all available extensions

# Running Extensions
ec run --plugin <name>           # Run Envoy with extensions
ec run --plugin ./path           # Run with local extension
ec run --config envoy.yaml       # Run with custom config

# Development
ec create-plugin <name>          # Scaffold new extension
ec plugin publish <path>         # Publish extension to marketplace

# Configuration Generation
ec gen --plugin <name>           # Generate Envoy config
ec gen --plugin <name> --full    # Generate full config
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

The CLI design is based on the [Envoy Composer UX document](../Envoy%20Composer%20UX.pdf).

If you're interested in contributing to the CLI implementation, please open an issue or discussion.

## Distribution

Planned distribution methods:
- Homebrew: `brew install envoy-ecosystem`
- Direct download: Binary releases on GitHub
- Curl script: `curl -sL https://envoy-ecosystem.io/install.sh | sh`
