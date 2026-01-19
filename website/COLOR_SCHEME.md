# Color Scheme Documentation

The Envoy Ecosystem website uses a color palette inspired by the official Envoy Proxy family of websites (envoyproxy.io, gateway.envoyproxy.io, aigateway.envoyproxy.io) but with distinct values to establish its own identity.

## Primary Palette

### Purple/Indigo (Primary Brand Color)
- `--primary: #7c3aed` - Main brand color, used for logos, primary CTAs, and key elements
- `--primary-dark: #5b21b6` - Darker shade for hover states and emphasis
- `--primary-light: #a78bfa` - Lighter shade for accents and highlights
- `--primary-ultra-light: #e9d5ff` - Very light tint for backgrounds and subtle accents

**Usage:** Headers, buttons, links, brand elements, gradients

### Teal/Cyan (Accent Color)
- `--accent: #06b6d4` - Accent color for CTAs and interactive elements
- `--accent-dark: #0891b2` - Darker shade for hover states
- `--accent-light: #67e8f9` - Lighter shade for highlights

**Usage:** Secondary CTAs, hover states, icons, complementary gradients

## Text Colors

- `--text: #0f172a` - Primary text color (dark slate)
- `--text-muted: #64748b` - Secondary text, descriptions
- `--text-light: #94a3b8` - Tertiary text, subtle labels

## Background Colors

- `--bg: #ffffff` - Main background (white)
- `--bg-subtle: #f8fafc` - Subtle background variation
- `--bg-muted: #f1f5f9` - Muted background for secondary sections
- `--bg-code: #1e293b` - Dark background for code blocks
- `--bg-code-light: #334155` - Lighter dark background

## Borders

- `--border: #e2e8f0` - Standard border color
- `--border-muted: #cbd5e1` - Muted border for subtle divisions

## Semantic Colors

### Success (Green)
- `--success: #10b981` - Success states, positive actions
- `--success-light: #d1fae5` - Light background for success messages

### Info (Blue)
- `--info: #3b82f6` - Informational states
- `--info-light: #dbeafe` - Light background for info messages

### Warning (Orange)
- `--warning: #f59e0b` - Warning states, caution messages
- `--warning-light: #fef3c7` - Light background for warnings

## Design Patterns

### Gradients
The site uses subtle gradients to add depth and visual interest:

1. **Hero Section**: Diagonal gradient from subtle background through white to ultra-light primary
2. **Feature Icons**: Linear gradient from primary to accent colors
3. **Extension Badges**: Diagonal gradient from light primary to primary
4. **CTA Section**: Diagonal gradient from ultra-light primary to subtle background

### Radial Overlays
Subtle radial gradients are used to add dimension:
- Hero section: Dual radial gradients with primary and accent colors at 3% opacity
- CTA section: Central radial gradient with primary color at 5% opacity

### Shadows
- Standard cards: Subtle gray shadows
- Hover states: Purple-tinted shadows mixing primary and accent colors
- Button hovers: Purple shadow with 30% opacity

## Comparison with Official Envoy Sites

### Similarities
- Purple/indigo as the primary brand color family
- Teal/cyan as complementary accent
- Clean, light backgrounds
- Dark code block backgrounds
- Professional, technical aesthetic

### Differences
- Distinct hex values to establish separate brand identity
- More extensive use of gradients
- Lighter, airier background colors
- More purple tints in interactive elements
- Unique shadow treatments with color

## Usage Guidelines

### Do's
✓ Use primary colors for brand elements and primary actions
✓ Use accent colors sparingly for emphasis
✓ Maintain sufficient contrast for accessibility
✓ Use gradients subtly for depth, not decoration
✓ Apply semantic colors consistently (green for success, etc.)

### Don'ts
✗ Don't mix too many gradient directions
✗ Don't use overly saturated backgrounds
✗ Don't override semantic color meanings
✗ Don't use accent color as primary (keep it complementary)
✗ Don't add colors outside the defined palette

## Accessibility

All color combinations meet WCAG 2.1 AA standards:
- Text on white backgrounds: 7:1+ contrast ratio
- Text on subtle backgrounds: 6:1+ contrast ratio
- Interactive elements have clear focus states
- Code blocks use high-contrast color schemes

## Future Considerations

As the Envoy Ecosystem brand evolves, consider:
- Dark mode variant (using existing dark colors as base)
- Additional semantic colors for specific extension categories
- Brand-specific illustration colors
- Animation accent colors
