"use client";
import React, { useState } from "react";
import { motion, AnimatePresence, useInView } from "framer-motion";
import { ShieldCheck, GitMerge, Cpu, Code2, TrendingUp, Layers, X, ArrowRight } from "lucide-react";

type FeatureId = "guard" | "mvcc" | "runtime" | "diff" | "bench" | "plugins" | null;

const features = [
  { 
    id: "guard",
    icon: ShieldCheck, 
    title: "Guard Pipeline", 
    desc: "Stateless circuit breaker verifying strict policy constraints with zero-latency enforcement.",
    color: "text-red-400",
    bg: "bg-red-500/10 border-red-500/20",
    expandedContent: (
      <div className="space-y-4 mt-4">
        <p className="text-text-muted text-sm leading-relaxed">
          The Guard Pipeline intercepts every agent action before execution. It prevents unauthorized code execution, limits API usage, and blocks malicious prompts dynamically.
        </p>
        <div className="bg-surface-2 p-4 rounded-xl border border-border">
          <h4 className="text-white text-sm font-semibold mb-2">Enforcement Rules</h4>
          <ul className="list-disc list-inside text-xs text-text-muted space-y-1.5">
            <li>No external network calls allowed in sandboxed tasks</li>
            <li>Maximum 5 file modifications per sub-agent</li>
            <li>PII & secrets detection using Regex guards</li>
          </ul>
        </div>
      </div>
    )
  },
  { 
    id: "mvcc",
    icon: GitMerge, 
    title: "MVCC Ledger", 
    desc: "Multi-version concurrency control for conflict-free file writes across parallel agents.",
    color: "text-amber-400",
    bg: "bg-amber-500/10 border-amber-500/20",
    expandedContent: (
      <div className="space-y-4 mt-4">
        <p className="text-text-muted text-sm leading-relaxed">
          When multiple agents collaborate, the MVCC Ledger handles file state versioning automatically. It prevents race conditions and gracefully merges non-overlapping changes.
        </p>
        <div className="bg-surface-2 p-3 rounded-xl border border-border flex items-center justify-between text-xs text-text-muted">
          <div className="text-center"><span className="block text-white mb-1">Agent 1</span>Editing src/app.ts</div>
          <ArrowRight className="w-4 h-4 text-accent" />
          <div className="text-center"><span className="block text-green-400 mb-1">Ledger</span>Auto-merged (v2.1)</div>
          <ArrowRight className="w-4 h-4 text-accent" />
          <div className="text-center"><span className="block text-white mb-1">Agent 2</span>Editing src/app.ts</div>
        </div>
      </div>
    )
  },
  { 
    id: "runtime",
    icon: Cpu, 
    title: "Go Runtime", 
    desc: "Single statically compiled binary with goroutine-based concurrency and sub-50ms cold starts.",
    color: "text-cyan-400",
    bg: "bg-cyan-500/10 border-cyan-500/20",
    expandedContent: (
      <div className="space-y-4 mt-4">
        <p className="text-text-muted text-sm leading-relaxed">
          We ditched the heavy V8 engine for a single, statically compiled Go binary. The result is a runtime that boots instantly, consumes a fraction of the memory, and scales effortlessly.
        </p>
        <div className="grid grid-cols-2 gap-3">
          <div className="bg-surface-2 p-3 rounded-xl border border-border text-center">
            <div className="text-[10px] uppercase text-text-muted tracking-wider mb-1">Cold Start</div>
            <div className="text-lg font-bold text-cyan-400 font-mono">42ms</div>
          </div>
          <div className="bg-surface-2 p-3 rounded-xl border border-border text-center">
            <div className="text-[10px] uppercase text-text-muted tracking-wider mb-1">Memory (Idle)</div>
            <div className="text-lg font-bold text-cyan-400 font-mono">14MB</div>
          </div>
        </div>
      </div>
    )
  },
  { 
    id: "diff",
    icon: Code2, 
    title: "Semantic Diff", 
    desc: "AST-based code analysis providing human-readable summaries of structural code mutations.",
    color: "text-violet-400",
    bg: "bg-violet-500/10 border-violet-500/20",
    expandedContent: (
      <div className="space-y-4 mt-4">
        <p className="text-text-muted text-sm leading-relaxed">
          Instead of showing you raw +/- lines, Gitclaw's Semantic Diff parses the code structure to tell you exactly what the agent did in plain English.
        </p>
        <div className="bg-surface-2 p-3 rounded-xl border border-border font-mono text-xs text-text-muted space-y-1.5">
          <div className="text-violet-400">Changed: src/auth.ts</div>
          <div><span className="text-green-400">+</span> Extracted function <span className="text-white">validateToken</span></div>
          <div><span className="text-green-400">+</span> Changed variable <span className="text-white">secret</span> to const</div>
          <div><span className="text-green-400">+</span> No logic changes detected</div>
        </div>
      </div>
    )
  },
  { 
    id: "bench",
    icon: TrendingUp, 
    title: "Benchmarking", 
    desc: "Evaluate agent performance across LLMs natively — test tokens, accuracy, and speed.",
    color: "text-emerald-400",
    bg: "bg-emerald-500/10 border-emerald-500/20",
    expandedContent: (
      <div className="space-y-4 mt-4">
        <p className="text-text-muted text-sm leading-relaxed">
          Run automated benchmarks to compare different LLMs (OpenAI, Anthropic, Ollama) against your specific repository tasks to find the most cost-effective and accurate model.
        </p>
        <div className="space-y-2 bg-surface-2 p-3 rounded-xl border border-border">
          <div className="flex justify-between text-xs"><span className="text-white">GPT-4o</span><span className="text-emerald-400">98% Pass (12s)</span></div>
          <div className="flex justify-between text-xs"><span className="text-white">Claude 3.5</span><span className="text-emerald-400">96% Pass (14s)</span></div>
          <div className="flex justify-between text-xs"><span className="text-white">Llama 3 (Local)</span><span className="text-amber-400">82% Pass (8s)</span></div>
        </div>
      </div>
    )
  },
  { 
    id: "plugins",
    icon: Layers, 
    title: "Plugin System", 
    desc: "Extensible architecture with tools, skills, hooks, memory layers, and prompt additions.",
    color: "text-blue-400",
    bg: "bg-blue-500/10 border-blue-500/20",
    expandedContent: (
      <div className="space-y-4 mt-4">
        <p className="text-text-muted text-sm leading-relaxed">
          Gitclaw's modular architecture lets you hook into the agent lifecycle. Add custom tools, provide enterprise context, or intercept lifecycle events with ease.
        </p>
        <div className="flex flex-wrap gap-2 text-xs">
          <span className="px-2 py-1 rounded bg-surface-3 border border-border text-text-muted">Tools</span>
          <span className="px-2 py-1 rounded bg-surface-3 border border-border text-text-muted">Memory Backends</span>
          <span className="px-2 py-1 rounded bg-surface-3 border border-border text-text-muted">Custom LLMs</span>
          <span className="px-2 py-1 rounded bg-surface-3 border border-border text-text-muted">Lifecycle Hooks</span>
        </div>
      </div>
    )
  },
];

export default function FeaturesSection() {
  const ref = React.useRef(null);
  const inView = useInView(ref, { once: true, margin: "-80px" });
  const [expandedId, setExpandedId] = useState<FeatureId>(null);

  const expandedFeature = features.find(f => f.id === expandedId);

  return (
    <section className="py-24 md:py-32 px-6 relative" ref={ref}>
      <div className="max-w-5xl mx-auto">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={inView ? { opacity: 1, y: 0 } : {}}
          transition={{ duration: 0.6 }}
          className="text-center mb-16 relative z-10"
        >
          <h2 className="text-3xl md:text-5xl font-bold mb-4">
            Everything You Need to Build<br />
            <span className="font-cursive text-accent italic text-4xl md:text-5xl">with Agents</span>
          </h2>
          <p className="text-text-muted text-sm md:text-base max-w-lg mx-auto">
            A complete runtime designed for production agent workflows — not just demos.
          </p>
        </motion.div>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 relative z-10">
          {features.map((f, i) => (
            <motion.div
              layoutId={`feature-card-${f.id}`}
              key={f.id}
              onClick={() => setExpandedId(f.id as FeatureId)}
              initial={{ opacity: 0, y: 20 }}
              animate={inView ? { opacity: 1, y: 0 } : {}}
              transition={{ delay: 0.1 + i * 0.08, duration: 0.5 }}
              className="dark-card p-6 group cursor-pointer"
            >
              <div className="w-11 h-11 rounded-xl bg-surface-2 border border-border-2 flex items-center justify-center mb-5 group-hover:border-accent/30 group-hover:bg-accent/10 transition-all duration-300">
                <f.icon className="w-5 h-5 text-text-muted group-hover:text-accent transition-colors" />
              </div>
              <h3 className="text-base font-semibold text-white mb-2 flex items-center justify-between" style={{ fontFamily: 'Inter, sans-serif' }}>
                {f.title}
                <ArrowRight className="w-4 h-4 opacity-0 -mr-4 group-hover:opacity-100 group-hover:mr-0 transition-all text-accent" />
              </h3>
              <p className="text-sm text-text-muted leading-relaxed line-clamp-2">{f.desc}</p>
            </motion.div>
          ))}
        </div>

        {/* Expanded Modal */}
        <AnimatePresence>
          {expandedFeature && (
            <>
              <motion.div
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                exit={{ opacity: 0 }}
                onClick={() => setExpandedId(null)}
                className="fixed inset-0 bg-background/80 backdrop-blur-sm z-50 cursor-pointer"
              />
              <motion.div
                layoutId={`feature-card-${expandedFeature.id}`}
                className="fixed inset-x-4 top-[15%] md:inset-x-auto md:left-1/2 md:-translate-x-1/2 md:w-[500px] z-50 bg-surface border border-border rounded-2xl shadow-2xl overflow-hidden"
              >
                <div className="p-6 md:p-8 relative">
                  <button
                    onClick={() => setExpandedId(null)}
                    className="absolute top-4 right-4 p-2 rounded-full bg-surface-2 text-text-muted hover:text-white hover:bg-surface-3 transition-colors"
                  >
                    <X className="w-5 h-5" />
                  </button>
                  <div className={`w-12 h-12 rounded-xl ${expandedFeature.bg} flex items-center justify-center mb-6`}>
                    <expandedFeature.icon className={`w-6 h-6 ${expandedFeature.color}`} />
                  </div>
                  <h2 className="text-2xl font-bold text-white mb-2 font-serif">
                    {expandedFeature.title}
                  </h2>
                  <p className="text-base text-white/90">
                    {expandedFeature.desc}
                  </p>
                  
                  <div className="mt-6 border-t border-border pt-4">
                    {expandedFeature.expandedContent}
                  </div>
                </div>
              </motion.div>
            </>
          )}
        </AnimatePresence>

      </div>
    </section>
  );
}
