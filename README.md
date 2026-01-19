# Envoy Ecosystem

A community-driven marketplace for Envoy Proxy extensions. Discover, run, and build custom filters with zero friction.

## Project Overview

**Envoy Ecosystem** is designed to make extending Envoy Proxy as simple as possible. It consists of:

1. **Marketplace Repository**: A GitHub repository where each folder contains an extension
2. **CLI Tool (`ee`)**: A command-line tool for discovering, running, and building extensions
3. **Website**: A developer-focused website for browsing extensions and getting started

## Repository Structure

```
envoy-ecosystem/
├── website/              # Official website (Astro-based)
│   ├── src/
│   │   └── pages/
│   │       └── index.astro
│   ├── netlify.toml     # Netlify deployment config
│   └── README.md        # Website documentation
│
├── extensions/          # (Coming soon) Extension marketplace
│   ├── rate-limiter/
│   ├── auth-jwt/
│   ├── cors/
│   └── ...
│
└── cli/                 # (Coming soon) CLI tool
```

## Quick Start

### Try an Extension

```bash
# Install the CLI
curl -sL https://envoy-ecosystem.io/install.sh | sh

# Run an extension
ee run --plugin rate-limiter --plugin auth-jwt
```

### Create Your Own Extension

```bash
# Scaffold a new extension
ee create-plugin my-filter

# Test it locally
ee run --plugin ./my-filter
```

### Publish an Extension

```bash
# Publish to the ecosystem
ee plugin publish ./my-filter
```

## CLI Commands

- `ee plugin list` - Browse available extensions
- `ee run --plugin <name>` - Start Envoy with extensions
- `ee run --plugin ./path` - Test local extensions
- `ee gen --plugin <name>` - Generate Envoy configuration
- `ee create-plugin <name>` - Scaffold new extension
- `ee plugin publish <path>` - Publish to marketplace

## Website Development

See [website/README.md](website/README.md) for website development instructions.

The website is built with [Astro](https://astro.build/) and deployed on Netlify.

```bash
cd website
npm install
npm run dev
```

## Design Philosophy

**Zero-friction developer experience**:
- One-liner installation
- Auto-generated configurations
- Built-in test environments
- Simple contribution process

## Contributing

We welcome contributions! To publish an extension:

1. Fork this repository
2. Add your extension in a new folder under `extensions/`
3. Include a `README.md` and manifest file
4. Open a pull request

Or use the CLI: `ee plugin publish ./your-extension`

## Tech Stack

- **Website**: Astro, Vanilla CSS
- **CLI**: (Implementation pending - likely Go or Rust)
- **Deployment**: Netlify (website), GitHub (marketplace)

## Resources

- [Website](https://envoy-ecosystem.io)
- [Envoy Proxy](https://www.envoyproxy.io/)
- [Tetrate](https://www.tetrate.io/)

## License

Copyright © 2026 Tetrate

Envoy® is a registered trademark of the Linux Foundation.
