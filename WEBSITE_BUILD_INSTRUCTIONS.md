# Next.js Website Build Instructions: Gitclaw V2 Launch

**Author:** Senior Staff Engineer
**Review Date:** May 24, 2026
**Objective:** Build a premium, stunning Next.js landing page to announce and showcase the Gitclaw architectural migration from TypeScript to Go, highlighting the massive improvements in speed, safety, and concurrency.

## 1. Tech Stack & Design Philosophy

As a senior engineer reviewing the plan for this launch site, my primary directive is: **This cannot look like a standard, boring developer tool documentation site.** We just executed a massive systems engineering feat. The website needs to feel like a high-end, premium product.

### Core Stack
- **Framework:** Next.js 14+ (App Router)
- **Styling:** Tailwind CSS (configured for a bespoke design system, not default colors)
- **Animations:** Framer Motion (crucial for micro-interactions and scroll-based storytelling)
- **Icons:** Lucide React
- **Typography:** Inter (sans-serif, UI) and JetBrains Mono (for code snippets and technical metrics)

### UI/UX Directives
1. **Premium Dark Mode:** Use deep, rich off-blacks (e.g., `#0A0A0B`) combined with neon/electric accents (Electric Blue `#0070F3`, Go Cyan `#00ADD8`, and Neon Purple `#7928CA`).
2. **Glassmorphism:** Use subtle translucent panels with backdrop blur to give depth and layering.
3. **Dynamic Interactivity:** Elements should react to the user. Hover effects on cards, subtle parallax scrolling on background meshes, and numbers that count up when scrolled into view.
4. **"Show, Don't Tell":** Don't just write "sub-50ms cold start". Show a visual comparison bar chart animating side-by-side with the old Node.js runtime.

---

## 2. Step-by-Step Implementation Guide

### Step 1: Initialization
Initialize a new Next.js project with Tailwind:
```bash
npx create-next-app@latest gitclaw-v2-launch
# Select: TypeScript, ESLint, Tailwind CSS, App Router
```

Install animation and styling dependencies:
```bash
npm install framer-motion lucide-react clsx tailwind-merge
```

### Step 2: Global Styles and Typography
Modify `tailwind.config.ts` to extend the color palette:
```typescript
module.exports = {
  theme: {
    extend: {
      colors: {
        background: '#09090b',
        surface: 'rgba(255, 255, 255, 0.05)',
        primary: '#00ADD8', // Go Cyan
        accent: '#7928CA',
      },
      backgroundImage: {
        'gradient-radial': 'radial-gradient(var(--tw-gradient-stops))',
        'hero-glow': 'conic-gradient(from 180deg at 50% 50%, #00ADD855 0deg, #7928CA55 180deg, #00ADD855 360deg)',
      }
    }
  }
}
```

---

## 3. Component Architecture & Storytelling

### A. The Hero Section (The Hook)
- **Visual:** A massive, glowing, animated text effect: "The Agent Runtime, Reimagined in Go."
- **Subtitle:** "Conflict-free. Zero-latency security. Sub-50ms cold starts."
- **Action:** Two glowing buttons: "View on GitHub" and "Read the Architecture Document".
- **Animation:** Use Framer Motion to stagger the entrance of the title, subtitle, and buttons. Have a subtle moving gradient orb in the background.

### B. The Performance Showdown (Why We Rewrote It)
- **Concept:** A side-by-side comparison section.
- **Visual:** Two vertical bars animating upwards. 
  - *Left Bar (Node.js/TS):* Animates slowly, stops at a high height (representing ~800ms+ start time). Color: Faded Gray.
  - *Right Bar (Go):* Snaps instantly to a very low height (representing <50ms). Color: Electric Blue.
- **Copy:** Explain why V8 engine boot times were bottlenecking agent workflows in CI/CD and pre-commit hooks, and how the Go binary solved it.

### C. The MVCC Write Ledger (Concurrency Visualized)
- **Concept:** Explaining the Multi-Version Concurrency Control.
- **Visual:** Create a sleek visual of three "Agent Avatars" shooting beams (writes) at a central "File" icon. Instead of crashing, the file icon splits into versions (v1, v2, v3) stacking neatly.
- **Copy:** Explain how this prevents race conditions and corrupted files when multiple agents operate in the same workspace.

### D. The Stateless Guard Pipeline (Security)
- **Concept:** The zero-latency circuit breaker.
- **Visual:** A terminal window mock-up. It types out `agent > rm -rf /`. Instantly, a red glowing shield icon slams down, and the terminal flashes red: `BLOCKED: Policy Violation`.
- **Copy:** Explain how the Go standard library allows us to intercept tool calls in memory instantly, preventing runaway agents before they touch the OS.

---

## 4. Senior Engineer Review Notes & Best Practices

1. **Performance is a Feature:** The website announcing our performance improvements *must itself be performant*. Ensure all images are optimized, fonts are preloaded, and Framer Motion animations only trigger `whileInView` to save CPU cycles.
2. **Semantic HTML & SEO:** Use proper `<section>`, `<article>`, and `<header>` tags. The page title should be `Gitclaw V2 | The High-Performance Go Agent Runtime`. Meta tags should describe the TS to Go migration.
3. **Responsive Design:** The terminal mockups and performance graphs must degrade gracefully to mobile screens. Stack the graphs vertically on small screens.
4. **Code Quality:** Ensure strict TypeScript typing for all custom components. Use a utility like `clsx` or `tailwind-merge` for conditionally combining class names on animated elements.

**Final Verdict:** Execute this with high fidelity. The UI must match the engineering excellence of the backend rewrite. Make it feel alive, fast, and secure.
