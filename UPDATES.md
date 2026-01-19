# Website Updates

## Color Scheme Refresh (January 2026)

### Overview
Updated the website color palette to be inspired by the official Envoy Proxy family of websites while maintaining a distinct identity for the Envoy Ecosystem project.

### Changes Made

#### Color Palette Update

**Before:**
- Primary: `#6500E0` (Bright purple)
- Accent: `#00D9C9` (Bright teal)
- Simple, single-color backgrounds

**After:**
- Primary: `#7c3aed` (Purple/Indigo with better contrast)
- Accent: `#06b6d4` (Teal/Cyan, more professional)
- Expanded palette with multiple shades and tints
- Added semantic colors with light variants

#### Visual Enhancements

1. **Hero Section**
   - Added diagonal gradient background
   - Implemented dual radial gradient overlays (3% opacity)
   - Creates depth without overwhelming content

2. **Feature Cards**
   - Added gradient top border on hover
   - Icons now use gradient fill (primary → accent)
   - Improved hover states with colored shadows
   - Subtle purple-tinted box shadows

3. **Extension Cards**
   - Extension badges now use gradient backgrounds
   - Enhanced hover effects with dual-color shadows
   - Border color transitions on hover

4. **CTA Section**
   - Added gradient background
   - Radial overlay for subtle emphasis
   - Improved button hover states with colored shadows

5. **Background Variations**
   - `--bg-subtle: #f8fafc` for alternating sections
   - `--bg-muted: #f1f5f9` for secondary elements
   - Better visual hierarchy through background tints

#### Technical Improvements

- Expanded CSS custom properties from 11 to 21 variables
- Added gradient patterns for consistent application
- Improved color contrast for accessibility
- All combinations meet WCAG 2.1 AA standards

### Design Inspiration Sources

Analyzed three official Envoy websites:
- https://www.envoyproxy.io - Main Envoy Proxy site
- https://gateway.envoyproxy.io - Envoy Gateway documentation
- https://aigateway.envoyproxy.io - AI Gateway documentation

**Key takeaways:**
- Purple/indigo as brand foundation
- Teal/cyan as complementary accent
- Clean, technical aesthetic
- Professional gradient usage
- Light, airy backgrounds

### Files Modified

- `/website/src/pages/index.astro` - Main page with updated colors and styles
- `/website/README.md` - Updated color documentation
- `/website/COLOR_SCHEME.md` - New comprehensive color guide

### Result

The updated design:
- Maintains brand consistency with Envoy family
- Establishes distinct Envoy Ecosystem identity
- Improves visual hierarchy and depth
- Enhances user experience with subtle animations
- Maintains excellent accessibility standards
- Preserves the zero-friction, developer-first UX

### Next Steps

Consider:
- Dark mode implementation using existing dark palette
- Additional brand assets (logos, icons)
- Illustration color palette for future graphics
- Animation enhancements for interactive elements
