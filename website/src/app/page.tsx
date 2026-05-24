"use client";

import React, { useEffect } from "react";
import { motion, useAnimation } from "framer-motion";
import { ArrowRight, ShieldAlert, Cpu, GitMerge, FileCode2, LineChart, FileTerminal, ArrowUpRight, TerminalSquare } from "lucide-react";
import clsx from "clsx";
import { twMerge } from "tailwind-merge";

function cn(...inputs: (string | undefined | null | false)[]) {
  return twMerge(clsx(inputs));
}

const HeroSection = () => {
  return (
    <section className="relative min-h-[90vh] flex flex-col justify-center px-4 md:px-16 overflow-hidden border-b-2 border-gray-800">
      <div className="absolute top-10 left-10 text-xs font-mono text-gray-600 uppercase tracking-widest hidden md:block">
        SYS.INIT // CORE.LOAD // V.0.2.0-GO
      </div>
      <div className="absolute top-10 right-10 text-xs font-mono text-primary flex items-center gap-2">
        <div className="w-2 h-2 bg-primary animate-pulse"></div> SYSTEM ONLINE
      </div>

      <motion.div
        initial={{ opacity: 0, x: -50 }}
        animate={{ opacity: 1, x: 0 }}
        transition={{ duration: 0.6, ease: "easeOut" }}
        className="z-10 max-w-5xl"
      >
        <div className="bg-primary text-black text-sm font-bold uppercase px-3 py-1 inline-block mb-6 tracking-wider">
          CRITICAL ARCHITECTURE UPDATE
        </div>
        <h1 className="text-6xl md:text-8xl font-black text-white mb-6 leading-[0.9]">
          THE AGENT RUNTIME,<br />REIMAGINED IN <span className="text-primary border-b-4 border-primary pb-2">GO</span>.
        </h1>
        <motion.p 
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.3, duration: 0.8 }}
          className="text-xl md:text-2xl text-gray-400 mb-10 max-w-2xl border-l-4 border-accent pl-6 py-2 bg-gradient-to-r from-accent/10 to-transparent"
        >
          Conflict-free MVCC Ledger. Zero-latency security pipeline. Sub-50ms cold starts.
        </motion.p>
        
        <motion.div 
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.6, duration: 0.5 }}
          className="flex flex-col sm:flex-row items-start gap-6 mt-12"
        >
          <a href="https://github.com/Creat1ve-shubh/gitagent" target="_blank" rel="noreferrer" className="flex items-center gap-2 px-8 py-4 bg-white text-black font-bold uppercase tracking-wider hover:bg-gray-200 transition-colors brutal-border">
            View Source Code <ArrowUpRight className="w-5 h-5" />
          </a>
          <a href="#engineering" className="flex items-center gap-2 px-8 py-4 bg-black border-2 border-gray-600 text-white font-bold uppercase tracking-wider hover:border-white transition-colors brutal-border">
            Engineering Submission <ArrowRight className="w-5 h-5" />
          </a>
        </motion.div>
      </motion.div>
      
      <div className="absolute bottom-0 right-10 opacity-10 pointer-events-none">
        <Cpu size={400} strokeWidth={0.5} />
      </div>
    </section>
  );
};

const PerformanceShowdown = () => {
  return (
    <section id="engineering" className="py-24 px-4 md:px-16 bg-background relative border-b-2 border-gray-800">
      <div className="max-w-7xl mx-auto flex flex-col lg:flex-row gap-16 items-center">
        <div className="lg:w-1/2">
          <h2 className="text-4xl md:text-6xl font-bold mb-6 text-white uppercase"><span className="text-accent">01.</span> Performance Showdown</h2>
          <p className="text-gray-400 text-lg mb-8 leading-relaxed">
            The V8 engine boot times were bottlenecking agent workflows. In a CI/CD pipeline, starting up Node.js for every agent task was unacceptable. 
            We replaced it with a single, statically compiled Go binary. The result? Near-instantaneous execution.
          </p>
          <div className="grid grid-cols-2 gap-6">
            <div className="border-l-2 border-gray-700 pl-4">
              <div className="text-sm text-gray-500 uppercase tracking-widest mb-1">Previous (Node.js)</div>
              <div className="text-3xl text-gray-300 font-mono">&gt;800ms</div>
            </div>
            <div className="border-l-2 border-accent pl-4">
              <div className="text-sm text-accent uppercase tracking-widest mb-1">Current (Go)</div>
              <div className="text-3xl text-white font-mono">&lt;50ms</div>
            </div>
          </div>
        </div>
        
        <div className="lg:w-1/2 w-full flex items-end justify-center gap-8 h-[400px] p-8 brutal-border bg-[#0a0a0a] relative">
          <div className="absolute top-4 left-4 text-xs text-gray-600 font-mono">BENCHMARK_COLD_START</div>
          <div className="flex flex-col items-center gap-4 w-1/3">
            <motion.div 
              initial={{ height: 0 }}
              whileInView={{ height: "100%" }}
              transition={{ duration: 1.5, ease: "easeOut" }}
              className="w-full max-w-[120px] bg-gray-800 border-2 border-gray-600 relative"
            >
              <div className="absolute -top-10 w-full text-center text-gray-400 font-mono font-bold">850ms</div>
            </motion.div>
          </div>
          
          <div className="flex flex-col items-center gap-4 w-1/3">
            <motion.div 
              initial={{ height: 0 }}
              whileInView={{ height: "10%" }}
              transition={{ duration: 0.2, ease: "easeOut", delay: 0.5 }}
              className="w-full max-w-[120px] bg-accent border-2 border-white relative shadow-[0_0_30px_rgba(0,255,65,0.4)]"
            >
              <div className="absolute -top-10 w-full text-center text-white font-mono font-bold text-xl">42ms</div>
            </motion.div>
          </div>
        </div>
      </div>
    </section>
  );
};

const ArchitectureDiagram = () => {
  return (
    <section className="py-24 px-4 md:px-16 bg-[#030303] relative border-b-2 border-gray-800">
      <div className="max-w-7xl mx-auto">
        <div className="mb-16">
          <h2 className="text-4xl md:text-6xl font-bold mb-6 text-white uppercase"><span className="text-primary">02.</span> Multi-Threaded Architecture</h2>
          <p className="text-gray-400 text-lg max-w-2xl">
            We abandoned the fragile async event loop. The new pipeline in Go provides zero-latency security checks and safe concurrent file modifications.
          </p>
        </div>

        <div className="p-8 brutal-border bg-black relative overflow-hidden">
          <div className="absolute inset-0 opacity-10 bg-[radial-gradient(circle_at_center,rgba(255,255,255,0.2)_1px,transparent_1px)] bg-[size:20px_20px]"></div>
          
          <div className="flex flex-col md:flex-row items-stretch justify-between gap-8 relative z-10 min-h-[300px]">
            
            {/* Step 1 */}
            <motion.div 
              whileHover={{ scale: 1.02 }} 
              className="flex-1 flex flex-col bg-[#0a0a0a] border border-gray-700 p-6 relative group"
            >
              <div className="absolute top-0 right-0 bg-gray-700 text-black text-xs font-bold px-2 py-1">STAGE 01</div>
              <FileTerminal className="w-12 h-12 text-white mb-6 group-hover:text-primary transition-colors" />
              <h3 className="text-xl font-bold text-white mb-2 uppercase">CLI Input</h3>
              <p className="text-gray-500 text-sm">User requests execution via prompt or IDE context.</p>
            </motion.div>
            
            <div className="hidden md:flex items-center justify-center">
              <ArrowRight className="text-gray-600 w-8 h-8" />
            </div>
            
            {/* Step 2 */}
            <motion.div 
              whileHover={{ scale: 1.02 }} 
              className="flex-1 flex flex-col bg-[#0a0a0a] border border-red-900 p-6 relative group shadow-[0_0_20px_rgba(255,42,42,0.1)]"
            >
              <div className="absolute top-0 right-0 bg-primary text-black text-xs font-bold px-2 py-1">STAGE 02</div>
              <ShieldAlert className="w-12 h-12 text-primary mb-6" />
              <h3 className="text-xl font-bold text-white mb-2 uppercase">Guard Pipeline</h3>
              <p className="text-gray-500 text-sm">Stateless circuit breaker. Intercepts all tool calls instantly to block unauthorized access.</p>
            </motion.div>
            
            <div className="hidden md:flex items-center justify-center">
              <ArrowRight className="text-gray-600 w-8 h-8" />
            </div>
            
            {/* Step 3 */}
            <motion.div 
              whileHover={{ scale: 1.02 }} 
              className="flex-1 flex flex-col bg-[#0a0a0a] border border-accent/40 p-6 relative group shadow-[0_0_20px_rgba(0,255,65,0.1)]"
            >
              <div className="absolute top-0 right-0 bg-accent text-black text-xs font-bold px-2 py-1">STAGE 03</div>
              <GitMerge className="w-12 h-12 text-accent mb-6" />
              <h3 className="text-xl font-bold text-white mb-2 uppercase">MVCC Ledger</h3>
              <p className="text-gray-500 text-sm">Conflict-free file management. Multiple agents can write concurrently without race conditions.</p>
            </motion.div>

          </div>
        </div>
      </div>
    </section>
  );
};

const ToolingSection = () => {
  return (
    <section className="py-24 px-4 md:px-16 bg-background relative border-b-2 border-gray-800">
      <div className="max-w-7xl mx-auto">
        <div className="mb-16">
          <h2 className="text-4xl md:text-6xl font-bold mb-6 text-white uppercase"><span className="text-white">03.</span> Next-Gen Tooling</h2>
          <p className="text-gray-400 text-lg max-w-2xl">Semantic Diff and Benchmarking suites built directly into the core runtime.</p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12">
          {/* Semantic Diff */}
          <div className="p-8 bg-[#0a0a0a] border border-gray-800 brutal-border group hover:border-gray-500">
            <div className="flex items-center gap-4 mb-8 pb-4 border-b border-gray-800">
              <div className="bg-white text-black p-3"><FileCode2 className="w-6 h-6" /></div>
              <h3 className="text-2xl font-bold uppercase tracking-wide">Semantic Diff</h3>
            </div>
            <p className="text-gray-400 mb-8 font-mono text-sm">git diff is noisy. gitclaw diff parses the AST for human-readable summaries.</p>
            <div className="bg-black p-6 border-l-4 border-accent font-mono text-sm leading-relaxed relative">
              <div className="absolute top-2 right-2 flex gap-2">
                <div className="w-3 h-3 rounded-full bg-red-500"></div>
                <div className="w-3 h-3 rounded-full bg-yellow-500"></div>
                <div className="w-3 h-3 rounded-full bg-green-500"></div>
              </div>
              <div className="text-gray-500 mb-4 mt-2">$ gitclaw diff --semantic</div>
              <div className="text-accent font-bold mb-4">&gt; ANALYZING AST...</div>
              <ul className="text-gray-300 space-y-3 pl-4 border-l-2 border-gray-800">
                <li className="before:content-['-'] before:mr-2 before:text-accent">Refactored <span className="text-white bg-gray-800 px-1">parse</span> function to ES6 arrow.</li>
                <li className="before:content-['-'] before:mr-2 before:text-accent">Upgraded variable <span className="text-white bg-gray-800 px-1">x</span> to const for immutability.</li>
                <li className="before:content-['-'] before:mr-2 before:text-accent">Logic remains completely identical.</li>
              </ul>
            </div>
          </div>

          {/* Benchmark */}
          <div className="p-8 bg-[#0a0a0a] border border-gray-800 brutal-border group hover:border-gray-500">
            <div className="flex items-center gap-4 mb-8 pb-4 border-b border-gray-800">
              <div className="bg-white text-black p-3"><LineChart className="w-6 h-6" /></div>
              <h3 className="text-2xl font-bold uppercase tracking-wide">Benchmarking</h3>
            </div>
            <p className="text-gray-400 mb-8 font-mono text-sm">Evaluate agent performance between LLMs natively.</p>
            
            <div className="overflow-x-auto border border-gray-800">
              <table className="w-full text-sm text-left text-gray-300 font-mono">
                <thead className="text-xs uppercase bg-black text-gray-400">
                  <tr>
                    <th className="px-6 py-4 border-b border-r border-gray-800">Metric</th>
                    <th className="px-6 py-4 border-b border-r border-gray-800">V1 (GPT-4o)</th>
                    <th className="px-6 py-4 border-b border-gray-800 text-white bg-white/5">V2 (Claude-3.5)</th>
                  </tr>
                </thead>
                <tbody className="bg-transparent">
                  <tr className="border-b border-gray-800">
                    <td className="px-6 py-4 border-r border-gray-800 font-bold text-white">Tokens</td>
                    <td className="px-6 py-4 border-r border-gray-800">12,450</td>
                    <td className="px-6 py-4 text-accent bg-white/5">8,100 (-34%)</td>
                  </tr>
                  <tr className="border-b border-gray-800">
                    <td className="px-6 py-4 border-r border-gray-800 font-bold text-white">Time</td>
                    <td className="px-6 py-4 border-r border-gray-800">45s</td>
                    <td className="px-6 py-4 text-accent bg-white/5">18s (-60%)</td>
                  </tr>
                  <tr>
                    <td className="px-6 py-4 border-r border-gray-800 font-bold text-white">Pass Rate</td>
                    <td className="px-6 py-4 border-r border-gray-800">85%</td>
                    <td className="px-6 py-4 text-accent bg-white/5">100% (+15%)</td>
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
    { phase: "PHASE 01: FOUNDATION", desc: "Initial release v0.1.0. Building the Node.js scaffold, adding the local repo sandbox, and integrating gitmachine." },
    { phase: "PHASE 02: INTERFACE", desc: "Introduction of the OpenAI Realtime voice adapter, IDE-style Monaco editor, and mobile-responsive UI." },
    { phase: "PHASE 03: THE BRAIN", desc: "Rolling out the plugin system, chat branching, background memory saving, and skill learning." },
    { phase: "PHASE 04: OBSERVABILITY", desc: "Adding OpenTelemetry instrumentation and the unified Logs tab for debugging." },
    { phase: "PHASE 05: THE OVERHAUL", desc: "Our Engineering Submission: Migrating from TS to Go, introducing the MVCC ledger, Stateless Circuit Breaker, semantic diff, and benchmarking.", isFinal: true },
  ];

  return (
    <section className="py-24 px-4 md:px-16 bg-[#030303] relative border-b-2 border-gray-800 overflow-hidden">
      <div className="absolute top-0 right-0 w-1/3 h-full bg-gradient-to-l from-primary/5 to-transparent pointer-events-none"></div>
      
      <div className="max-w-7xl mx-auto flex flex-col lg:flex-row gap-16">
        <div className="lg:w-1/3">
          <h2 className="text-4xl md:text-6xl font-bold mb-6 text-white uppercase leading-[0.9]"><span className="text-gray-600 block text-2xl mb-2">04.</span> Evolution<br/>Log</h2>
          <p className="text-gray-400 text-lg">A full commit scan timeline demonstrating our mastery of the repository.</p>
        </div>

        <div className="lg:w-2/3 relative border-l-4 border-gray-800 pl-8 md:pl-12">
          {events.map((evt, i) => (
            <motion.div 
              key={i}
              initial={{ opacity: 0, x: -20 }}
              whileInView={{ opacity: 1, x: 0 }}
              viewport={{ once: true, margin: "-100px" }}
              transition={{ delay: i * 0.15 }}
              className="mb-16 relative"
            >
              <div className={cn(
                "absolute -left-[42px] md:-left-[58px] top-1 w-6 h-6 border-4 border-[#030303]",
                evt.isFinal ? "bg-primary rounded-none shadow-[0_0_15px_rgba(255,42,42,0.6)]" : "bg-gray-600 rounded-none"
              )}></div>
              <h3 className={cn("text-2xl font-bold mb-3 uppercase tracking-wide", evt.isFinal ? "text-primary" : "text-white")}>{evt.phase}</h3>
              <p className="text-gray-400 font-mono text-sm max-w-xl leading-relaxed">{evt.desc}</p>
              
              {evt.isFinal && (
                <div className="mt-6 inline-flex items-center gap-2 bg-primary text-black px-4 py-2 font-bold uppercase text-sm">
                  <TerminalSquare className="w-4 h-4" /> Final Submission Complete
                </div>
              )}
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
      <div className="noise-overlay"></div>
      <div className="scanline"></div>
      
      <nav className="flex items-center justify-between p-6 border-b-2 border-gray-800 relative z-40 bg-background/80 backdrop-blur-sm">
        <div className="text-2xl font-black uppercase tracking-widest flex items-center gap-3">
          <div className="w-8 h-8 bg-white text-black flex items-center justify-center">G</div>
          GITCLAW<span className="text-primary font-mono text-sm font-bold ml-2 border border-primary px-1">V2</span>
        </div>
        <div className="hidden md:flex gap-8 text-sm font-bold uppercase tracking-wider text-gray-400">
          <a href="#engineering" className="hover:text-white transition-colors">Engineering</a>
          <a href="#" className="hover:text-white transition-colors">Architecture</a>
          <a href="#" className="text-primary hover:text-white transition-colors">Hire Us</a>
        </div>
      </nav>

      <HeroSection />
      <PerformanceShowdown />
      <ArchitectureDiagram />
      <ToolingSection />
      <TimelineSection />
      
      <footer className="py-12 px-16 bg-black border-t-2 border-primary text-center">
        <div className="text-3xl font-black uppercase tracking-widest text-white mb-4">GITCLAW</div>
        <p className="text-gray-600 font-mono text-xs uppercase tracking-widest">Built for the Senior Engineering Submission &copy; 2026 // END OF LINE.</p>
      </footer>
    </main>
  );
}
