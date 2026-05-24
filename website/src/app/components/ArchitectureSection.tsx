"use client";
import React, { useState } from "react";
import { motion, AnimatePresence, useInView } from "framer-motion";
import { ShieldCheck, GitMerge, Cpu, ArrowRight, Terminal, Database, Zap, X } from "lucide-react";

type ExpandedId = "guard" | "router" | "mvcc" | "diff" | "perf" | null;

export default function ArchitectureSection() {
  const ref = React.useRef(null);
  const inView = useInView(ref, { once: true, margin: "-80px" });
  const [expandedId, setExpandedId] = useState<ExpandedId>(null);

  const expandedContent = {
    guard: {
      title: "Guard Pipeline",
      icon: ShieldCheck,
      color: "text-red-400",
      bg: "bg-red-500/10 border-red-500/20",
      details: (
        <div className="space-y-4">
          <p className="text-text-muted text-sm leading-relaxed">
            The Guard Pipeline acts as a stateless circuit breaker for every agent action. It enforces strict policy constraints, prevents prompt injections, and guarantees zero-latency enforcement before any system change is committed.
          </p>
          <div className="bg-surface-2 p-4 rounded-xl border border-border">
            <h4 className="text-white text-sm font-semibold mb-3">Security Features:</h4>
            <ul className="list-disc list-inside text-sm text-text-muted space-y-2">
              <li>Regex-based and semantic prompt scanning</li>
              <li>Resource limits and quota enforcement</li>
              <li>Real-time threat telemetry</li>
              <li>RBAC integration with Enterprise SSO</li>
            </ul>
          </div>
        </div>
      )
    },
    router: {
      title: "Go Router",
      icon: Cpu,
      color: "text-cyan-400",
      bg: "bg-cyan-500/10 border-cyan-500/20",
      details: (
        <div className="space-y-4">
          <p className="text-text-muted text-sm leading-relaxed">
            Our multi-threaded dispatcher relies on lightweight goroutines to handle agent requests. This architecture replaces the fragile Node.js event loop, unlocking massive concurrency.
          </p>
          <div className="bg-surface-2 p-4 rounded-xl border border-border font-mono text-xs text-text-muted space-y-2">
            <div className="text-cyan-400">// Handling 10k+ concurrent connections</div>
            <div>func handleAgent(req *Request) {"{"}</div>
            <div className="ml-4">ctx := context.WithTimeout(req.Context(), 5*time.Second)</div>
            <div className="ml-4">go processTask(ctx, req.Task)</div>
            <div>{"}"}</div>
          </div>
        </div>
      )
    },
    mvcc: {
      title: "MVCC Write Ledger",
      icon: GitMerge,
      color: "text-amber-400",
      bg: "bg-amber-500/10 border-amber-500/20",
      details: (
        <div className="space-y-4">
          <p className="text-text-muted text-sm leading-relaxed">
            Multi-Version Concurrency Control (MVCC) ensures that when multiple agents attempt to modify the same file or workspace simultaneously, their changes are intelligently versioned and merged without race conditions.
          </p>
          <div className="bg-surface-2 p-4 rounded-xl border border-border">
            <div className="flex items-center justify-between text-xs text-text-muted mb-2">
              <span>Agent A (Write)</span>
              <ArrowRight className="w-3 h-3 text-accent" />
              <span className="text-green-400">Version 2.1</span>
            </div>
            <div className="flex items-center justify-between text-xs text-text-muted mb-2">
              <span>Agent B (Write)</span>
              <ArrowRight className="w-3 h-3 text-accent" />
              <span className="text-green-400">Version 2.2</span>
            </div>
            <div className="flex items-center justify-between text-xs font-semibold text-white mt-4 border-t border-border pt-2">
              <span>Ledger Resolution</span>
              <span className="text-accent">Auto-Merge Strategy</span>
            </div>
          </div>
        </div>
      )
    },
    diff: {
      title: "Semantic Diff Engine",
      icon: Terminal,
      color: "text-violet-400",
      bg: "bg-violet-500/10 border-violet-500/20",
      details: (
        <div className="space-y-4">
          <p className="text-text-muted text-sm leading-relaxed">
            Standard unified diffs focus on lines and characters. Our Semantic Diff Engine parses the Abstract Syntax Tree (AST) to explain changes in plain English, understanding exactly what code logic was modified.
          </p>
          <div className="bg-surface-2 p-4 rounded-xl border border-border">
            <h4 className="text-white text-sm font-semibold mb-2">Capabilities:</h4>
            <div className="grid grid-cols-2 gap-2 text-xs text-text-muted">
              <div className="p-2 border border-border/50 rounded-lg">Variable Renames</div>
              <div className="p-2 border border-border/50 rounded-lg">Function Extractions</div>
              <div className="p-2 border border-border/50 rounded-lg">Type Changes</div>
              <div className="p-2 border border-border/50 rounded-lg">Logic Reordering</div>
            </div>
          </div>
        </div>
      )
    },
    perf: {
      title: "Performance Gains",
      icon: Zap,
      color: "text-emerald-400",
      bg: "bg-emerald-500/10 border-emerald-500/20",
      details: (
        <div className="space-y-4">
          <p className="text-text-muted text-sm leading-relaxed">
            Moving to Go eliminated the V8 engine boot time and drastically reduced the memory footprint, enabling agents to run instantly even in constrained serverless environments.
          </p>
          <div className="space-y-3 bg-surface-2 p-4 rounded-xl border border-border">
            <div>
              <div className="flex justify-between text-xs mb-1"><span className="text-white">Cold Start (Go)</span><span className="text-emerald-400">42ms</span></div>
              <div className="w-full bg-surface-3 rounded-full h-1.5"><div className="bg-emerald-400 h-1.5 rounded-full" style={{ width: '5%' }}></div></div>
            </div>
            <div>
              <div className="flex justify-between text-xs mb-1"><span className="text-text-muted">Cold Start (Node.js)</span><span className="text-text-muted">850ms</span></div>
              <div className="w-full bg-surface-3 rounded-full h-1.5"><div className="bg-text-muted h-1.5 rounded-full" style={{ width: '85%' }}></div></div>
            </div>
          </div>
        </div>
      )
    }
  };

  return (
    <section id="architecture" className="py-24 md:py-32 px-6 relative" ref={ref}>
      <div className="absolute inset-0 dot-pattern opacity-40" />
      <div className="max-w-6xl mx-auto relative z-10">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={inView ? { opacity: 1, y: 0 } : {}}
          transition={{ duration: 0.6 }}
          className="mb-14"
        >
          <p className="text-accent text-sm font-semibold tracking-wide uppercase mb-3">Architecture</p>
          <h2 className="text-3xl md:text-5xl font-bold leading-tight">
            Designed for Production<br />
            <span className="font-cursive text-accent italic">Agent Workflows</span>
          </h2>
        </motion.div>

        {/* Bento Grid */}
        <div className="bento-grid">
          {/* Card 1: Guard Pipeline - wide */}
          <motion.div
            layoutId="card-guard"
            onClick={() => setExpandedId("guard")}
            initial={{ opacity: 0, y: 20 }}
            animate={inView ? { opacity: 1, y: 0 } : {}}
            transition={{ delay: 0.1, duration: 0.5 }}
            className="dark-card p-6 md:p-8 bento-wide group cursor-pointer"
          >
            <div className="flex items-start justify-between mb-6">
              <div>
                <h3 className="text-lg font-semibold text-white mb-1 flex items-center gap-2" style={{ fontFamily: 'Inter, sans-serif' }}>
                  Guard Pipeline
                  <ArrowRight className="w-4 h-4 opacity-0 -ml-2 group-hover:opacity-100 group-hover:ml-0 transition-all text-accent" />
                </h3>
                <p className="text-sm text-text-muted">Stateless circuit breaker enforcing policy constraints with 0.3ms latency</p>
              </div>
              <div className="w-10 h-10 rounded-xl bg-red-500/10 border border-red-500/20 flex items-center justify-center shrink-0">
                <ShieldCheck className="w-5 h-5 text-red-400" />
              </div>
            </div>
            {/* Pipeline visual */}
            <div className="bg-surface-2 rounded-xl p-4 border border-border">
              <div className="flex items-center gap-3 flex-wrap">
                {["Request", "Auth Check", "Rate Limit", "Policy Gate", "✓ Pass"].map((step, i) => (
                  <React.Fragment key={step}>
                    <div className={`px-3 py-1.5 rounded-lg text-xs font-medium ${i === 4 ? 'bg-green-500/10 text-green-400 border border-green-500/20' : 'bg-surface-3 text-text-muted border border-border'}`}>
                      {step}
                    </div>
                    {i < 4 && <ArrowRight className="w-3 h-3 text-text-muted hidden sm:block" />}
                  </React.Fragment>
                ))}
              </div>
            </div>
          </motion.div>

          {/* Card 2: Go Runtime */}
          <motion.div
            layoutId="card-router"
            onClick={() => setExpandedId("router")}
            initial={{ opacity: 0, y: 20 }}
            animate={inView ? { opacity: 1, y: 0 } : {}}
            transition={{ delay: 0.2, duration: 0.5 }}
            className="dark-card p-6 md:p-8 group cursor-pointer"
          >
            <div className="w-10 h-10 rounded-xl bg-cyan-500/10 border border-cyan-500/20 flex items-center justify-center mb-5">
              <Cpu className="w-5 h-5 text-cyan-400" />
            </div>
            <h3 className="text-lg font-semibold text-white mb-1 flex items-center gap-2" style={{ fontFamily: 'Inter, sans-serif' }}>
              Go Router
              <ArrowRight className="w-4 h-4 opacity-0 -ml-2 group-hover:opacity-100 group-hover:ml-0 transition-all text-accent" />
            </h3>
            <p className="text-sm text-text-muted mb-5">Multi-threaded dispatcher with goroutine pooling</p>
            <div className="bg-surface-2 rounded-xl p-4 border border-border font-mono text-xs">
              <div className="text-cyan-400">func main() {"{"}</div>
              <div className="text-text-muted ml-4">router := NewRouter()</div>
              <div className="text-text-muted ml-4">router.Use(GuardMiddleware)</div>
              <div className="text-accent ml-4">go router.Serve(":8080")</div>
              <div className="text-cyan-400">{"}"}</div>
            </div>
          </motion.div>

          {/* Card 3: MVCC Ledger */}
          <motion.div
            layoutId="card-mvcc"
            onClick={() => setExpandedId("mvcc")}
            initial={{ opacity: 0, y: 20 }}
            animate={inView ? { opacity: 1, y: 0 } : {}}
            transition={{ delay: 0.25, duration: 0.5 }}
            className="dark-card p-6 md:p-8 group cursor-pointer"
          >
            <div className="w-10 h-10 rounded-xl bg-amber-500/10 border border-amber-500/20 flex items-center justify-center mb-5">
              <GitMerge className="w-5 h-5 text-amber-400" />
            </div>
            <h3 className="text-lg font-semibold text-white mb-1 flex items-center gap-2" style={{ fontFamily: 'Inter, sans-serif' }}>
              MVCC Write Ledger
              <ArrowRight className="w-4 h-4 opacity-0 -ml-2 group-hover:opacity-100 group-hover:ml-0 transition-all text-accent" />
            </h3>
            <p className="text-sm text-text-muted mb-5">Conflict-free concurrent writes across multiple agents</p>
            {/* Version visualization */}
            <div className="space-y-2">
              {["Agent A → v1", "Agent B → v2", "Agent C → v3"].map((v, i) => (
                <motion.div
                  key={v}
                  initial={{ opacity: 0, x: -10 }}
                  animate={inView ? { opacity: 1, x: 0 } : {}}
                  transition={{ delay: 0.6 + i * 0.1, duration: 0.4 }}
                  className="flex items-center gap-2 text-xs"
                >
                  <div className={`w-2 h-2 rounded-full ${i === 0 ? 'bg-blue-400' : i === 1 ? 'bg-purple-400' : 'bg-green-400'}`} />
                  <span className="text-text-muted font-mono">{v}</span>
                  <span className="text-green-400 ml-auto">✓</span>
                </motion.div>
              ))}
            </div>
          </motion.div>

          {/* Card 4: Semantic Diff - wide */}
          <motion.div
            layoutId="card-diff"
            onClick={() => setExpandedId("diff")}
            initial={{ opacity: 0, y: 20 }}
            animate={inView ? { opacity: 1, y: 0 } : {}}
            transition={{ delay: 0.3, duration: 0.5 }}
            className="dark-card p-6 md:p-8 bento-wide group cursor-pointer"
          >
            <div className="flex items-start justify-between mb-6">
              <div>
                <h3 className="text-lg font-semibold text-white mb-1 flex items-center gap-2" style={{ fontFamily: 'Inter, sans-serif' }}>
                  Semantic Diff Engine
                  <ArrowRight className="w-4 h-4 opacity-0 -ml-2 group-hover:opacity-100 group-hover:ml-0 transition-all text-accent" />
                </h3>
                <p className="text-sm text-text-muted">AST-based analysis providing human-readable mutation summaries</p>
              </div>
              <div className="w-10 h-10 rounded-xl bg-violet-500/10 border border-violet-500/20 flex items-center justify-center shrink-0">
                <Terminal className="w-5 h-5 text-violet-400" />
              </div>
            </div>
            <div className="bg-surface-2 rounded-xl p-4 border border-border font-mono text-xs space-y-2">
              <div className="text-text-muted">$ gitclaw diff --semantic</div>
              <div className="text-accent">→ Refactored parse function to ES6 arrow.</div>
              <div className="text-accent">→ Upgraded variable x to const.</div>
              <div className="text-green-400">→ Logic remains completely identical.</div>
            </div>
          </motion.div>

          {/* Card 5: Performance */}
          <motion.div
            layoutId="card-perf"
            onClick={() => setExpandedId("perf")}
            initial={{ opacity: 0, y: 20 }}
            animate={inView ? { opacity: 1, y: 0 } : {}}
            transition={{ delay: 0.35, duration: 0.5 }}
            className="dark-card p-6 md:p-8 group cursor-pointer"
          >
            <div className="w-10 h-10 rounded-xl bg-emerald-500/10 border border-emerald-500/20 flex items-center justify-center mb-5">
              <Zap className="w-5 h-5 text-emerald-400" />
            </div>
            <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2" style={{ fontFamily: 'Inter, sans-serif' }}>
              Performance Gains
              <ArrowRight className="w-4 h-4 opacity-0 -ml-2 group-hover:opacity-100 group-hover:ml-0 transition-all text-accent" />
            </h3>
            <div className="space-y-3">
              {[
                { label: "Cold Start", from: "850ms", to: "42ms" },
                { label: "Binary", from: "180MB", to: "9.2MB" },
                { label: "Guard", from: "12ms", to: "0.3ms" },
              ].map((m) => (
                <div key={m.label} className="flex items-center justify-between text-xs">
                  <span className="text-text-muted">{m.label}</span>
                  <div className="flex items-center gap-2">
                    <span className="text-text-muted line-through">{m.from}</span>
                    <ArrowRight className="w-3 h-3 text-accent" />
                    <span className="text-white font-semibold">{m.to}</span>
                  </div>
                </div>
              ))}
            </div>
          </motion.div>
        </div>

        {/* Expanded Modal */}
        <AnimatePresence>
          {expandedId && (
            <>
              <motion.div
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                exit={{ opacity: 0 }}
                onClick={() => setExpandedId(null)}
                className="fixed inset-0 bg-background/80 backdrop-blur-sm z-50 cursor-pointer"
              />
              <motion.div
                layoutId={`card-${expandedId}`}
                className="fixed inset-x-4 top-[10%] md:inset-x-auto md:left-1/2 md:-translate-x-1/2 md:w-[600px] z-50 bg-surface border border-border rounded-2xl shadow-2xl overflow-hidden"
              >
                <div className="p-6 md:p-8 relative">
                  <button
                    onClick={() => setExpandedId(null)}
                    className="absolute top-4 right-4 p-2 rounded-full bg-surface-2 text-text-muted hover:text-white hover:bg-surface-3 transition-colors"
                  >
                    <X className="w-5 h-5" />
                  </button>
                  <div className={`w-12 h-12 rounded-xl ${expandedContent[expandedId].bg} flex items-center justify-center mb-6`}>
                    {React.createElement(expandedContent[expandedId].icon, { className: `w-6 h-6 ${expandedContent[expandedId].color}` })}
                  </div>
                  <h2 className="text-2xl font-bold text-white mb-6 font-serif">
                    {expandedContent[expandedId].title}
                  </h2>
                  <div className="text-base text-text-muted">
                    {expandedContent[expandedId].details}
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
