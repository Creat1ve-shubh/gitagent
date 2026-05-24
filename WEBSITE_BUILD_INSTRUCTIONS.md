# Next.js Website Build Instructions: Gitclaw V2 Launch (Team Submission)

**Author:** Senior Staff Engineer Candidate (Our Team Submission)
**Review Date:** May 24, 2026
**Objective:** Build a premium, stunning Next.js landing page to announce and showcase the Gitclaw architectural migration from TypeScript to Go. This website serves as our engineering portfolio submission to demonstrate our mastery of the repository, highlighting massive improvements in speed, safety, tooling, and concurrency.

## 1. Tech Stack & Design Philosophy

As a senior engineering team applying for this role, our primary directive is: **This cannot look like a standard, boring developer tool documentation site.** We just executed a massive systems engineering feat. The website needs to feel like a high-end, premium product that proves our caliber.

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

---

## 3. Component Architecture & Storytelling

### A. The Hero Section (The Hook)
- **Visual:** A massive, glowing, animated text effect: "The Agent Runtime, Reimagined in Go."
- **Subtitle:** "Conflict-free. Zero-latency security. Sub-50ms cold starts."
- **Action:** Two glowing buttons: "View on GitHub" and "View Our Engineering Submission".

### B. The Performance Showdown (Why We Rewrote It)
- **Concept:** A side-by-side comparison section.
- **Visual:** Two vertical bars animating upwards to represent cold start times. 
  - *Left Bar (Node.js/TS):* Animates slowly, stops at a high height (representing ~800ms+ start time). Color: Faded Gray.
  - *Right Bar (Go):* Snaps instantly to a very low height (representing <50ms). Color: Electric Blue.

### C. The New Architecture Diagram
- **Concept:** Visualize our new Go-based routing and pipeline system.
- **Implementation:** Use a sleek, styled Mermaid.js component or a custom Framer Motion SVG animation that maps out the flow:
  `User -> Go Router -> Circuit Breaker -> Agent Engine -> MVCC Write Ledger -> File System`.
- **Copy:** Explain how we replaced the fragile Node.js async event loop with a robust, multi-threaded Go architecture.

### D. The MVCC Write Ledger (Concurrency Visualized)
- **Concept:** Explaining the Multi-Version Concurrency Control.
- **Visual:** Create a sleek visual of three "Agent Avatars" shooting beams (writes) at a central "File" icon. Instead of crashing, the file icon splits into versions (v1, v2, v3) stacking neatly.

### E. Next-Gen Tooling: Semantic Diff & Benchmark
- **Concept:** Highlight our newly added developer experience tools.
- **Visual - Semantic Diff:** A split-pane code window. On the left, a standard messy `git diff` with red and green lines. On the right, our clean `gitclaw diff` output explaining the AST-level changes in plain English.
- **Visual - Benchmarking:** A beautiful data table showing `gitclaw bench` output, comparing "Agent V1" vs "Agent V2" on metrics like Token Usage, Tool Calls, and Task Completion Time.

---

## 4. Full Development Timeline (Our Commit Scan)

To prove our deep understanding of the entire repository's history, we want a dedicated **"Evolution of Gitclaw"** timeline section on the website. This should be an interactive vertical timeline component:

1. **Phase 1: The Foundation:** `Initial release v0.1.0`. Building the Node.js scaffold, adding the local repo sandbox, and integrating gitmachine.
2. **Phase 2: UI & Voice:** Introduction of the OpenAI Realtime voice adapter, IDE-style Monaco editor, and mobile-responsive UI.
3. **Phase 3: The Agent Brain:** Rolling out the plugin system, chat branching, background memory saving, and skill learning.
4. **Phase 4: Observability:** Adding OpenTelemetry instrumentation and the unified Logs tab for debugging.
5. **Phase 5: The Grand Overhaul (Our Submission):** The massive `overhaul` commits. Migrating from TS to Go, introducing the MVCC ledger, the Stateless Circuit Breaker, `gitclaw diff`, and `gitclaw bench`.

---

## 5. Senior Engineer Review Notes & Best Practices

1. **Performance is a Feature:** The website announcing our performance improvements *must itself be performant*. Ensure all images are optimized, fonts are preloaded, and Framer Motion animations only trigger `whileInView` to save CPU cycles.
2. **Semantic HTML & SEO:** Use proper `<section>`, `<article>`, and `<header>` tags. The page title should be `Gitclaw V2 | Engineering Submission`.
3. **Responsive Design:** The terminal mockups and performance graphs must degrade gracefully to mobile screens. Stack the graphs vertically on small screens.

**Final Verdict:** Execute this with high fidelity. The UI must represent our technical capabilities. This isn't just a website; this is our final submission to prove we are the Senior Engineers for the job. Make it feel alive, fast, and secure.
