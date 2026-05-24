"use client";
import React from "react";
import { motion, useInView } from "framer-motion";
import { ShieldCheck, GitMerge, Cpu, Code2, TrendingUp, Layers } from "lucide-react";

const features = [
  { icon: ShieldCheck, title: "Guard Pipeline", desc: "Stateless circuit breaker verifying strict policy constraints with zero-latency enforcement." },
  { icon: GitMerge, title: "MVCC Ledger", desc: "Multi-version concurrency control for conflict-free file writes across parallel agents." },
  { icon: Cpu, title: "Go Runtime", desc: "Single statically compiled binary with goroutine-based concurrency and sub-50ms cold starts." },
  { icon: Code2, title: "Semantic Diff", desc: "AST-based code analysis providing human-readable summaries of structural code mutations." },
  { icon: TrendingUp, title: "Benchmarking", desc: "Evaluate agent performance across LLMs natively — test tokens, accuracy, and speed." },
  { icon: Layers, title: "Plugin System", desc: "Extensible architecture with tools, skills, hooks, memory layers, and prompt additions." },
];

export default function FeaturesSection() {
  const ref = React.useRef(null);
  const inView = useInView(ref, { once: true, margin: "-80px" });

  return (
    <section className="py-24 md:py-32 px-6 relative" ref={ref}>
      <div className="max-w-5xl mx-auto">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={inView ? { opacity: 1, y: 0 } : {}}
          transition={{ duration: 0.6 }}
          className="text-center mb-16"
        >
          <h2 className="text-3xl md:text-5xl font-bold mb-4">
            Everything You Need to Build<br />
            <span className="font-cursive text-accent italic text-4xl md:text-5xl">with Agents</span>
          </h2>
          <p className="text-text-muted text-sm md:text-base max-w-lg mx-auto">
            A complete runtime designed for production agent workflows — not just demos.
          </p>
        </motion.div>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {features.map((f, i) => (
            <motion.div
              key={f.title}
              initial={{ opacity: 0, y: 20 }}
              animate={inView ? { opacity: 1, y: 0 } : {}}
              transition={{ delay: 0.1 + i * 0.08, duration: 0.5 }}
              className="dark-card p-6 group cursor-default"
            >
              <div className="w-11 h-11 rounded-xl bg-surface-2 border border-border-2 flex items-center justify-center mb-5 group-hover:border-accent/30 group-hover:bg-accent/10 transition-all duration-300">
                <f.icon className="w-5 h-5 text-text-muted group-hover:text-accent transition-colors" />
              </div>
              <h3 className="text-base font-semibold text-white mb-2" style={{ fontFamily: 'Inter, sans-serif' }}>{f.title}</h3>
              <p className="text-sm text-text-muted leading-relaxed">{f.desc}</p>
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  );
}
