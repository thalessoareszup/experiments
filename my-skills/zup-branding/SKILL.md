---
name: zup-color-branding
description: Provides Zup's institutional color palette and design guidelines for web applications. Use when building web interfaces, selecting colors, or ensuring brand consistency in digital products.
---

# Zup Color Branding Guide

Apply Zup's color palette to ensure brand consistency and professional design in all web applications.

## Primary Palette (Institutional Colors)

These colors represent Zup's core identity and must be prioritized. Use them for primary backgrounds, dominant design elements, and visual hierarchy.

### 01. Burgundy Excellence
- **Hex**: `#260A12`
- **RGB**: 38, 10, 18
- **CMYK**: 12, 100, 50, 74
- **Purpose**: Primary base for backgrounds, emphasizing maturity and professionalism
- **Use Cases**: Primary backgrounds, dominant sections, brand anchors

### 02. Gray Impact
- **Hex**: `#DED4D4`
- **RGB**: 222, 212, 212
- **CMYK**: 11, 14, 11, 0
- **Purpose**: Adds contrast, conveying sophistication and modernity
- **Use Cases**: Secondary backgrounds, text on dark backgrounds, neutral spaces

### 03. Connection Gradient
- **Hex Range**: `#5C1C27` (light) to `#260A12` (dark)
- **RGB**: 92, 28, 39 to 38, 10, 18
- **CMYK**: Analogous to Burgundy Excellence
- **Purpose**: Creates visual movement and depth
- **Use Cases**: Gradient overlays, background transitions, dynamic elements
- **Direction**: Lighter at bottom-left and top-right, darker toward center, creating upward motion

## Complementary Palette

Use exclusively for small detail elements such as text highlights, borders, bullets, and icons. Never use as predominant colors.

### 01. Burgundy Sensibility
- **Hex**: `#852838`
- **RGB**: 133, 40, 56
- **CMYK**: 30, 92, 68, 30
- **Use Cases**: Text highlights, icon accents, subtle borders

### 02. Gray Knowledge
- **Hex**: `#AD9797`
- **RGB**: 173, 151, 151
- **CMYK**: 0, 13, 11, 40
- **Use Cases**: Secondary icons, subtle dividers, muted text

### 03. Coral Boldness
- **Hex**: `#CC7958`
- **RGB**: 204, 121, 88
- **CMYK**: 16, 58, 65, 5
- **Use Cases**: Call-to-action accents, warning elements, emphasis

## Color Application Guidelines

**Primary Elements**
- Page backgrounds: Burgundy Excellence (#260A12)
- Surface backgrounds: Gray Impact (#DED4D4)
- Text on light: Burgundy Excellence (#260A12)
- Text on dark: Gray Impact (#DED4D4)

**Secondary Elements**
- Highlight text: Burgundy Sensibility (#852838) or Coral Boldness (#CC7958)
- Icon tints: Gray Knowledge (#AD9797)
- Subtle borders: Gray Knowledge (#AD9797)
- Action buttons: Coral Boldness (#CC7958) with Burgundy Excellence text

**Gradients**
- Hero sections: Connection Gradient (#5C1C27 → #260A12)
- Overlays: Connection Gradient with opacity
- Transitions: Smooth gradients between burgundy tones

## Example Hex Palette for CSS/Design Tools

```
Primary:
  - #260A12 (Burgundy Excellence)
  - #DED4D4 (Gray Impact)
  - #5C1C27 (Connection Gradient Light)

Complementary:
  - #852838 (Burgundy Sensibility)
  - #AD9797 (Gray Knowledge)
  - #CC7958 (Coral Boldness)
```

## Best Practices

1. **Hierarchy**: Primary palette for 80% of design, complementary for 20% accents
2. **Contrast**: Always verify text readability using Web Content Accessibility Guidelines
3. **Consistency**: Use exact hex values across all platforms (web, mobile, design tools)
4. **Accessibility**: Burgundy Excellence (#260A12) on Gray Impact (#DED4D4) has strong contrast (WCAG AAA compliant)
5. **Realism**: Test colors on actual displays and devices before final implementation
