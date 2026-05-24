"use client";
import React from "react";
import { motion, useInView } from "framer-motion";
import { ShieldCheck, GitMerge, Cpu, ArrowRight, Terminal, Database, Zap, Lock } from "lucide-react";

export default function ArchitectureSection() {
  const ref = React.useRef(null);
  const inView = useInView(ref, { once: true, margin: "-80px" });

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
            initial={{ opacity: 0, y: 20 }}
            animate={inView ? { opacity: 1, y: 0 } : {}}
            transition={{ delay: 0.1, duration: 0.5 }}
            className="dark-card p-6 md:p-8 bento-wide group"
          >
            <div className="flex items-start justify-between mb-6">
              <div>
                <h3 className="text-lg font-semibold text-white mb-1" style={{ fontFamily: 'Inter, sans-serif' }}>Guard Pipeline</h3>
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
            initial={{ opacity: 0, y: 20 }}
            animate={inView ? { opacity: 1, y: 0 } : {}}
            transition={{ delay: 0.2, duration: 0.5 }}
            className="dark-card p-6 md:p-8 group"
          >
            <div className="w-10 h-10 rounded-xl bg-cyan-500/10 border border-cyan-500/20 flex items-center justify-center mb-5">
              <Cpu className="w-5 h-5 text-cyan-400" />
            </div>
            <h3 className="text-lg font-semibold text-white mb-1" style={{ fontFamily: 'Inter, sans-serif' }}>Go Router</h3>
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
            initial={{ opacity: 0, y: 20 }}
            animate={inView ? { opacity: 1, y: 0 } : {}}
            transition={{ delay: 0.25, duration: 0.5 }}
            className="dark-card p-6 md:p-8 group"
          >
            <div className="w-10 h-10 rounded-xl bg-amber-500/10 border border-amber-500/20 flex items-center justify-center mb-5">
              <GitMerge className="w-5 h-5 text-amber-400" />
            </div>
            <h3 className="text-lg font-semibold text-white mb-1" style={{ fontFamily: 'Inter, sans-serif' }}>MVCC Write Ledger</h3>
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
            initial={{ opacity: 0, y: 20 }}
            animate={inView ? { opacity: 1, y: 0 } : {}}
            transition={{ delay: 0.3, duration: 0.5 }}
            className="dark-card p-6 md:p-8 bento-wide group"
          >
            <div className="flex items-start justify-between mb-6">
              <div>
                <h3 className="text-lg font-semibold text-white mb-1" style={{ fontFamily: 'Inter, sans-serif' }}>Semantic Diff Engine</h3>
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
            initial={{ opacity: 0, y: 20 }}
            animate={inView ? { opacity: 1, y: 0 } : {}}
            transition={{ delay: 0.35, duration: 0.5 }}
            className="dark-card p-6 md:p-8 group"
          >
            <div className="w-10 h-10 rounded-xl bg-emerald-500/10 border border-emerald-500/20 flex items-center justify-center mb-5">
              <Zap className="w-5 h-5 text-emerald-400" />
            </div>
            <h3 className="text-lg font-semibold text-white mb-4" style={{ fontFamily: 'Inter, sans-serif' }}>Performance Gains</h3>
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
      </div>
    </section>
  );
}
