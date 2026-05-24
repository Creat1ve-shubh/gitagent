"use client";
import React from "react";
import { motion, useInView } from "framer-motion";

const benchmarks = [
  { metric: "Cold Start", v1: "850ms", v2: "42ms", change: "-95%", highlight: true },
  { metric: "Tokens Used", v1: "12,450", v2: "8,100", change: "-34%", highlight: false },
  { metric: "Pass Rate", v1: "85%", v2: "100%", change: "+15%", highlight: true },
  { metric: "Binary Size", v1: "180MB", v2: "9.2MB", change: "-95%", highlight: false },
  { metric: "Guard Latency", v1: "12ms", v2: "0.3ms", change: "-97%", highlight: true },
];

export default function PerformanceSection() {
  const ref = React.useRef(null);
  const inView = useInView(ref, { once: true, margin: "-80px" });

  return (
    <section className="py-24 md:py-32 px-6 relative" ref={ref}>
      <div className="max-w-5xl mx-auto">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
          {/* Left: Copy */}
          <motion.div
            initial={{ opacity: 0, x: -20 }}
            animate={inView ? { opacity: 1, x: 0 } : {}}
            transition={{ duration: 0.6 }}
          >
            <p className="text-accent text-sm font-semibold tracking-wide uppercase mb-3">Performance</p>
            <h2 className="text-3xl md:text-5xl font-bold leading-tight mb-6">
              The Speed<br />
              <span className="font-cursive text-accent italic">Showdown</span>
            </h2>
            <p className="text-text-muted text-sm leading-relaxed mb-8">
              V8 engine boot times bottlenecked agent workflows. We replaced it with a single, statically compiled Go binary. The result speaks for itself.
            </p>

            {/* Big numbers */}
            <div className="grid grid-cols-2 gap-6">
              <div className="dark-card p-5">
                <div className="text-[10px] font-semibold text-text-muted uppercase tracking-widest mb-2">Node.js</div>
                <div className="text-3xl font-black text-text-muted/40 font-mono">&gt;800<span className="text-lg">ms</span></div>
              </div>
              <div className="gradient-border p-5">
                <div className="text-[10px] font-semibold text-accent uppercase tracking-widest mb-2">Go Runtime</div>
                <motion.div
                  initial={{ opacity: 0, scale: 0.8 }}
                  animate={inView ? { opacity: 1, scale: 1 } : {}}
                  transition={{ delay: 0.3, duration: 0.5, type: "spring" }}
                  className="text-3xl font-black text-white font-mono"
                >
                  &lt;50<span className="text-lg">ms</span>
                </motion.div>
              </div>
            </div>
          </motion.div>

          {/* Right: Benchmark Table */}
          <motion.div
            initial={{ opacity: 0, x: 20 }}
            animate={inView ? { opacity: 1, x: 0 } : {}}
            transition={{ delay: 0.2, duration: 0.6 }}
          >
            <div className="dark-card overflow-hidden">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-border">
                    <th className="px-5 py-4 text-left text-xs font-semibold text-text-muted uppercase tracking-wider">Metric</th>
                    <th className="px-5 py-4 text-left text-xs font-semibold text-text-muted uppercase tracking-wider">V1</th>
                    <th className="px-5 py-4 text-left text-xs font-semibold text-accent uppercase tracking-wider">V2 (Go)</th>
                  </tr>
                </thead>
                <tbody>
                  {benchmarks.map((row, i) => (
                    <motion.tr
                      key={row.metric}
                      initial={{ opacity: 0 }}
                      animate={inView ? { opacity: 1 } : {}}
                      transition={{ delay: 0.4 + i * 0.08, duration: 0.4 }}
                      className="border-b border-border last:border-0 hover:bg-surface-2/50 transition-colors"
                    >
                      <td className="px-5 py-3.5 text-text-muted font-medium">{row.metric}</td>
                      <td className="px-5 py-3.5 text-text-muted/60 font-mono">{row.v1}</td>
                      <td className="px-5 py-3.5 font-mono">
                        <span className="text-white font-semibold">{row.v2}</span>
                        <span className={`ml-2 text-xs ${row.highlight ? 'text-green-400' : 'text-accent'}`}>({row.change})</span>
                      </td>
                    </motion.tr>
                  ))}
                </tbody>
              </table>
            </div>
          </motion.div>
        </div>
      </div>
    </section>
  );
}
