import type { Config } from 'tailwindcss'

// Palette/type scale placeholders only — the real design-token
// system (colors, fonts, radii, spacing) gets defined in STEP 2
// (Design System). Kept minimal here on purpose.
export default {
  content: ['./index.html', './src/**/*.{vue,ts}'],
  theme: {
    extend: {},
  },
  plugins: [],
} satisfies Config
