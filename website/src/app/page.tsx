"use client";

import React from "react";
import { motion } from "framer-motion";
import { ArrowRight, ShieldAlert, Cpu, GitMerge, FileCode2, LineChart, FileTerminal, ArrowUpRight } from "lucide-react";
import clsx from "clsx";
import { twMerge } from "tailwind-merge";

function cn(...inputs: (string | undefined | null | false)[]) {
  return twMerge(clsx(inputs));
}

const HeroSection = () => {
  return (
    <section className="relative min-h-screen flex flex-col items-center justify-center overflow-hidden px-4">
      <div className="absolute inset-0 hero-glow opacity-30 pointer-events-none rounded-full transform -translate-y-1/2 scale-150"></div>
      
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.8, ease: "easeOut" }}
        className="z-10 text-center max-w-4xl"
      >
        <h1 className="text-5xl md:text-7xl font-extrabold tracking-tight text-transparent bg-clip-text bg-gradient-to-r from-white to-gray-400 mb-6 drop-shadow-lg">
          The Agent Runtime,<br />Reimagined in <span className="text-primary">Go</span>.
        </h1>
        <motion.p 
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.3, duration: 0.8 }}
          className="text-xl md:text-2xl text-gray-300 mb-10 font-medium"
        >
          Conflict-free. Zero-latency security. Sub-50ms cold starts.
        </motion.p>
        
        <motion.div 
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.6, duration: 0.5 }}
          className="flex flex-col sm:flex-row items-center justify-center gap-4"
        >
          <a href="https://github.com/Creat1ve-shubh/gitagent" target="_blank" rel="noreferrer" className="flex items-center gap-2 px-8 py-4 rounded-full bg-white text-black font-semibold hover:bg-gray-200 transition-colors">
            View on GitHub <ArrowUpRight className="w-5 h-5" />
          </a>
          <a href="#engineering" className="flex items-center gap-2 px-8 py-4 rounded-full bg-surface border border-gray-700 hover:border-gray-500 text-white font-semibold transition-colors backdrop-blur-md">
            View Our Engineering Submission
          </a>
        </motion.div>
      </motion.div>
    </section>
  );
};

const PerformanceShowdown = () => {
  return (
    <section id="engineering" className="py-24 px-4 bg-background relative border-t border-surface">
      <div className="max-w-6xl mx-auto">
        <div className="mb-16 text-center">
          <h2 className="text-3xl md:text-5xl font-bold mb-6">The Performance Showdown</h2>
          <p className="text-gray-400 text-lg max-w-2xl mx-auto">Why we rewrote Gitclaw: V8 engine boot times were bottlenecking agent workflows. The compiled Go binary solves this instantly.</p>
        </div>
        
        <div className="flex flex-col md:flex-row items-end justify-center gap-12 h-96 mt-12 bg-surface p-8 rounded-2xl border border-gray-800">
          <div className="flex flex-col items-center gap-4 w-1/3">
            <span className="text-4xl font-mono text-gray-500 font-bold">&gt;800ms</span>
            <motion.div 
              initial={{ height: 0 }}
              whileInView={{ height: "100%" }}
              transition={{ duration: 1.5, ease: "easeOut" }}
              className="w-24 bg-gray-700 rounded-t-lg relative"
            >
              <div className="absolute bottom-0 w-full text-center py-2 text-xs text-gray-300 border-t border-gray-600">Node.js (V8)</div>
            </motion.div>
          </div>
          
          <div className="flex flex-col items-center gap-4 w-1/3">
            <span className="text-4xl font-mono text-primary font-bold">&lt;50ms</span>
            <motion.div 
              initial={{ height: 0 }}
              whileInView={{ height: "15%" }}
              transition={{ duration: 0.2, ease: "easeOut" }}
              className="w-24 bg-primary rounded-t-lg relative shadow-[0_0_20px_rgba(0,173,216,0.5)]"
            >
              <div className="absolute bottom-0 w-full text-center py-2 text-xs font-bold text-black bg-white/20">Go Binary</div>
            </motion.div>
          </div>
        </div>
      </div>
    </section>
  );
};

const ArchitectureDiagram = () => {
  return (
    <section className="py-24 px-4 bg-[#050505] relative border-t border-surface">
      <div className="max-w-6xl mx-auto">
        <div className="mb-16 text-center">
          <h2 className="text-3xl md:text-5xl font-bold mb-6">Multi-Threaded Architecture</h2>
          <p className="text-gray-400 text-lg max-w-2xl mx-auto">We replaced the fragile async event loop with a robust pipeline in Go.</p>
        </div>

        <div className="relative p-8 rounded-2xl border border-gray-800 bg-surface backdrop-blur-md overflow-hidden">
          <div className="flex flex-col md:flex-row items-center justify-between gap-6 relative z-10">
            <motion.div whileHover={{ scale: 1.05 }} className="flex flex-col items-center p-6 bg-black rounded-xl border border-gray-700 shadow-lg">
              <FileTerminal className="w-10 h-10 text-white mb-2" />
              <span className="font-semibold text-sm">CLI / Request</span>
            </motion.div>
            
            <ArrowRight className="text-gray-600 w-8 h-8 hidden md:block" />
            
            <motion.div whileHover={{ scale: 1.05 }} className="flex flex-col items-center p-6 bg-black rounded-xl border border-red-900 shadow-[0_0_15px_rgba(220,38,38,0.2)]">
              <ShieldAlert className="w-10 h-10 text-red-500 mb-2" />
              <span className="font-semibold text-sm">Guard Pipeline</span>
            </motion.div>
            
            <ArrowRight className="text-gray-600 w-8 h-8 hidden md:block" />
            
            <motion.div whileHover={{ scale: 1.05 }} className="flex flex-col items-center p-6 bg-black rounded-xl border border-primary shadow-[0_0_15px_rgba(0,173,216,0.2)]">
              <Cpu className="w-10 h-10 text-primary mb-2" />
              <span className="font-semibold text-sm">Agent Engine</span>
            </motion.div>
            
            <ArrowRight className="text-gray-600 w-8 h-8 hidden md:block" />
            
            <motion.div whileHover={{ scale: 1.05 }} className="flex flex-col items-center p-6 bg-black rounded-xl border border-accent shadow-[0_0_15px_rgba(121,40,202,0.2)]">
              <GitMerge className="w-10 h-10 text-accent mb-2" />
              <span className="font-semibold text-sm">MVCC Ledger</span>
            </motion.div>
          </div>
        </div>
      </div>
    </section>
  );
};

const ToolingSection = () => {
  return (
    <section className="py-24 px-4 bg-background relative border-t border-surface">
      <div className="max-w-6xl mx-auto">
        <div className="mb-16 text-center">
          <h2 className="text-3xl md:text-5xl font-bold mb-6">Next-Gen Tooling</h2>
          <p className="text-gray-400 text-lg max-w-2xl mx-auto">Semantic Diff and Benchmarking suites built directly into the CLI.</p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
          {/* Semantic Diff */}
          <div className="p-8 rounded-2xl bg-surface border border-gray-800 backdrop-blur-md">
            <div className="flex items-center gap-3 mb-6">
              <FileCode2 className="w-8 h-8 text-primary" />
              <h3 className="text-2xl font-bold">Semantic Diff</h3>
            </div>
            <p className="text-gray-400 mb-6 text-sm">Standard diffs are noisy. Gitclaw diff parses the AST to provide a human-readable summary.</p>
            <div className="bg-black p-4 rounded-xl border border-gray-800 font-mono text-sm leading-relaxed">
              <div className="text-gray-500 mb-2">// gitclaw diff</div>
              <div className="text-primary font-bold mb-2">🤖 Semantic Diff:</div>
              <ul className="list-disc list-inside text-gray-300 space-y-2">
                <li>Refactored <span className="text-white bg-gray-800 px-1 rounded">parse</span> function to ES6 arrow function.</li>
                <li>Upgraded variable <span className="text-white bg-gray-800 px-1 rounded">x</span> to const for strict immutability.</li>
                <li>Logic remains completely identical.</li>
              </ul>
            </div>
          </div>

          {/* Benchmark */}
          <div className="p-8 rounded-2xl bg-surface border border-gray-800 backdrop-blur-md">
            <div className="flex items-center gap-3 mb-6">
              <LineChart className="w-8 h-8 text-accent" />
              <h3 className="text-2xl font-bold">Agent Benchmarking</h3>
            </div>
            <p className="text-gray-400 mb-6 text-sm">Evaluate agent performance between models with the native benchmark suite.</p>
            
            <div className="overflow-hidden rounded-xl border border-gray-800">
              <table className="w-full text-sm text-left text-gray-300">
                <thead className="text-xs uppercase bg-black border-b border-gray-800">
                  <tr>
                    <th className="px-4 py-3">Metric</th>
                    <th className="px-4 py-3">V1 (GPT-4o)</th>
                    <th className="px-4 py-3 text-accent">V2 (Claude-3.5)</th>
                  </tr>
                </thead>
                <tbody className="bg-[#0c0c0e]">
                  <tr className="border-b border-gray-800">
                    <td className="px-4 py-3">Tokens</td>
                    <td className="px-4 py-3">12,450</td>
                    <td className="px-4 py-3 text-white font-bold">8,100 (-34%)</td>
                  </tr>
                  <tr className="border-b border-gray-800">
                    <td className="px-4 py-3">Time</td>
                    <td className="px-4 py-3">45s</td>
                    <td className="px-4 py-3 text-white font-bold">18s (-60%)</td>
                  </tr>
                  <tr>
                    <td className="px-4 py-3">Pass Rate</td>
                    <td className="px-4 py-3">85%</td>
                    <td className="px-4 py-3 text-white font-bold">100%</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
};

const TimelineSection = () => {
  const events = [
    { phase: "Phase 1: The Foundation", desc: "Initial release v0.1.0. Building the Node.js scaffold, adding the local repo sandbox, and integrating gitmachine." },
    { phase: "Phase 2: UI & Voice", desc: "Introduction of the OpenAI Realtime voice adapter, IDE-style Monaco editor, and mobile-responsive UI." },
    { phase: "Phase 3: The Agent Brain", desc: "Rolling out the plugin system, chat branching, background memory saving, and skill learning." },
    { phase: "Phase 4: Observability", desc: "Adding OpenTelemetry instrumentation and the unified Logs tab for debugging." },
    { phase: "Phase 5: The Grand Overhaul", desc: "Our Engineering Submission: Migrating from TS to Go, introducing the MVCC ledger, Stateless Circuit Breaker, semantic diff, and benchmarking." },
  ];

  return (
    <section className="py-24 px-4 bg-[#050505] relative border-t border-surface">
      <div className="max-w-4xl mx-auto">
        <div className="mb-16 text-center">
          <h2 className="text-3xl md:text-5xl font-bold mb-6">Evolution of Gitclaw</h2>
          <p className="text-gray-400 text-lg max-w-2xl mx-auto">A full commit scan timeline demonstrating our mastery of the repository.</p>
        </div>

        <div className="relative border-l border-gray-800 ml-4 md:ml-0">
          {events.map((evt, i) => (
            <motion.div 
              key={i}
              initial={{ opacity: 0, x: -20 }}
              whileInView={{ opacity: 1, x: 0 }}
              viewport={{ once: true }}
              transition={{ delay: i * 0.2 }}
              className="mb-10 ml-8 relative"
            >
              <div className={cn(
                "absolute -left-10 w-4 h-4 rounded-full border-2 border-background",
                i === events.length - 1 ? "bg-accent shadow-[0_0_10px_#7928ca]" : "bg-gray-600"
              )}></div>
              <h3 className={cn("text-xl font-bold mb-2", i === events.length - 1 ? "text-accent" : "text-white")}>{evt.phase}</h3>
              <p className="text-gray-400">{evt.desc}</p>
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  );
};

export default function Home() {
  return (
    <main className="bg-background text-white min-h-screen font-sans selection:bg-primary selection:text-black">
      <HeroSection />
      <PerformanceShowdown />
      <ArchitectureDiagram />
      <ToolingSection />
      <TimelineSection />
      
      <footer className="py-8 text-center border-t border-surface text-gray-500 text-sm">
        <p>Built for the Senior Engineering Submission. &copy; 2026</p>
      </footer>
    </main>
  );
}
