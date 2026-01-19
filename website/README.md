# Envoy Ecosystem Website

The official website for the Envoy Ecosystem project by Tetrate.

## Development

### Prerequisites

- Node.js 18+ and npm

### Getting Started

Install dependencies:

```bash
npm install
```

Start the development server:

```bash
npm run dev
```

The site will be available at `http://localhost:4321`

### Building for Production

Build the site:

```bash
npm run build
```

Preview the production build:

```bash
npm run preview
```

## Deployment

This site is configured for deployment on Netlify. Simply connect your repository to Netlify and it will automatically deploy on every push to the main branch.

### Netlify Configuration

The `netlify.toml` file is already configured with:
- Build command: `npm run build`
- Publish directory: `dist`

### Manual Deployment

You can also deploy manually:

```bash
npm run build
netlify deploy --prod --dir=dist
```

## Tech Stack

- **Framework**: [Astro](https://astro.build/) - Fast, modern static site generator
- **Styling**: Vanilla CSS with CSS custom properties
- **Deployment**: Netlify

## Project Structure

```
/
├── public/          # Static assets (favicon, images, etc.)
├── src/
│   └── pages/       # Page components
│       └── index.astro  # Homepage
├── astro.config.mjs # Astro configuration
├── netlify.toml     # Netlify deployment configuration
└── package.json
```

## Design Philosophy

The website follows these key principles:

- **Zero-friction UX**: Code blocks as primary CTAs
- **Developer-first**: Minimal marketing copy, maximum code examples
- **Immediate value**: Users can start experimenting without scrolling
- **Performance**: Ships zero JavaScript by default (except for copy buttons)

## Updating Content

The website is a single-page application with all content in `src/pages/index.astro`. To update:

1. Edit the content directly in the Astro component
2. The site uses CSS custom properties for theming (defined in the `:root` selector)
3. Featured extensions can be added/modified in the extensions section

## Color Scheme

The site uses Envoy-inspired colors:

- Primary: `#6500E0` (Envoy purple)
- Accent: `#00D9C9` (Teal)
- Success: `#00C853` (Green)
- Info: `#2196F3` (Blue)

## License

Copyright © 2026 Tetrate
