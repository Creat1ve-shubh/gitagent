"use client";

import React from "react";
import { motion } from "framer-motion";
import { ArrowRight, ShieldCheck, Cpu, GitMerge, Code2, TrendingUp, Sparkles, User, Box } from "lucide-react";
import clsx from "clsx";
import { twMerge } from "tailwind-merge";

function cn(...inputs: (string | undefined | null | false)[]) {
  return twMerge(clsx(inputs));
}

const HeroSection = () => {
  return (
    <section className="relative min-h-[90vh] flex flex-col justify-center px-6 md:px-16 overflow-hidden">
      <div className="absolute top-1/4 -right-20 w-96 h-96 bg-secondary rounded-full mix-blend-multiply filter blur-3xl opacity-70 animate-blob"></div>
      <div className="absolute top-1/3 -left-20 w-72 h-72 bg-[#e8e4d9] rounded-full mix-blend-multiply filter blur-3xl opacity-70 animate-blob animation-delay-2000"></div>

      <motion.div
        initial={{ opacity: 0, y: 30 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 1, ease: "easeOut" }}
        className="z-10 max-w-5xl mx-auto text-center mt-20"
      >
        <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-white border border-gray-200 shadow-sm mb-8">
          <Sparkles className="w-4 h-4 text-accent" />
          <span className="text-sm font-medium tracking-wide text-gray-600 uppercase">Submission to Lyzr</span>
        </div>
        
        <h1 className="text-6xl md:text-8xl font-black text-primary mb-6 leading-[1.1] tracking-tight">
          The Agent Runtime,<br />
          <span className="font-cursive text-7xl md:text-9xl text-accent font-normal italic pr-4">Reimagined</span> in Go.
        </h1>
        
        <motion.p 
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.4, duration: 1 }}
          className="text-xl md:text-2xl text-gray-500 mb-12 max-w-3xl mx-auto font-light leading-relaxed"
        >
          A masterclass in modern systems engineering. Featuring a conflict-free MVCC Ledger, zero-latency security pipeline, and sub-50ms cold starts.
        </motion.p>
        
        <motion.div 
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.7, duration: 0.8 }}
          className="flex flex-col sm:flex-row items-center justify-center gap-6"
        >
          <a href="#architecture" className="flex items-center gap-2 px-8 py-4 bg-primary text-white rounded-full font-medium hover:bg-gray-800 transition-all shadow-lg hover:shadow-xl transform hover:-translate-y-1">
            Explore Architecture <ArrowRight className="w-5 h-5" />
          </a>
          <a href="https://github.com/Creat1ve-shubh/gitagent" target="_blank" rel="noreferrer" className="flex items-center gap-2 px-8 py-4 bg-white text-primary rounded-full font-medium border border-gray-200 hover:border-gray-400 transition-all shadow-sm">
            View Source Code
          </a>
        </motion.div>
      </motion.div>
    </section>
  );
};

const ArchitectureDiagram = () => {
  return (
    <section id="architecture" className="py-32 px-6 md:px-16 bg-white relative">
      <div className="max-w-6xl mx-auto">
        <div className="text-center mb-24">
          <span className="font-cursive text-3xl text-accent mb-4 block">01. Flow & Safety</span>
          <h2 className="text-4xl md:text-6xl font-bold text-primary">Multi-Threaded <br/> Architecture</h2>
        </div>

        <div className="relative glass-panel p-8 md:p-16">
          <div className="flex flex-col md:flex-row items-center justify-between gap-6 relative z-10">
            
            {/* Step 1 */}
            <motion.div 
              whileHover={{ y: -5 }} 
              className="flex-1 w-full bg-[#fdfbf7] border border-gray-100 rounded-3xl p-8 shadow-sm flex flex-col items-center text-center"
            >
              <div className="w-16 h-16 bg-white rounded-full shadow-sm flex items-center justify-center mb-6 border border-gray-100">
                <User className="w-8 h-8 text-primary" />
              </div>
              <h3 className="text-2xl font-bold text-primary mb-3">User Input</h3>
              <p className="text-gray-500 text-sm leading-relaxed">Direct interaction via IDE or CLI triggering the agent runtime.</p>
            </motion.div>
            
            <ArrowRight className="text-gray-300 w-10 h-10 hidden md:block" />
            
            {/* Step 2 */}
            <motion.div 
              whileHover={{ y: -5 }} 
              className="flex-1 w-full bg-[#fdfbf7] border border-gray-100 rounded-3xl p-8 shadow-sm flex flex-col items-center text-center relative"
            >
              <div className="absolute top-4 right-4 w-3 h-3 bg-red-400 rounded-full animate-pulse"></div>
              <div className="w-16 h-16 bg-white rounded-full shadow-sm flex items-center justify-center mb-6 border border-gray-100">
                <ShieldCheck className="w-8 h-8 text-primary" />
              </div>
              <h3 className="text-2xl font-bold text-primary mb-3">Guard Pipeline</h3>
              <p className="text-gray-500 text-sm leading-relaxed">Stateless circuit breaker verifying strict policy constraints instantly.</p>
            </motion.div>
            
            <ArrowRight className="text-gray-300 w-10 h-10 hidden md:block" />
            
            {/* Step 3 */}
            <motion.div 
              whileHover={{ y: -5 }} 
              className="flex-1 w-full bg-[#fdfbf7] border border-gray-100 rounded-3xl p-8 shadow-sm flex flex-col items-center text-center"
            >
              <div className="w-16 h-16 bg-white rounded-full shadow-sm flex items-center justify-center mb-6 border border-gray-100">
                <GitMerge className="w-8 h-8 text-primary" />
              </div>
              <h3 className="text-2xl font-bold text-primary mb-3">MVCC Ledger</h3>
              <p className="text-gray-500 text-sm leading-relaxed">Safe concurrent writes allowing multiple agents without race conditions.</p>
            </motion.div>

          </div>
        </div>
      </div>
    </section>
  );
};

const PerformanceShowdown = () => {
  return (
    <section className="py-32 px-6 md:px-16 bg-secondary relative">
      <div className="max-w-7xl mx-auto flex flex-col lg:flex-row gap-20 items-center">
        
        <div className="lg:w-1/2">
          <span className="font-cursive text-3xl text-accent mb-4 block">02. Engineering Feat</span>
          <h2 className="text-4xl md:text-6xl font-bold text-primary mb-8 leading-tight">The Performance Showdown</h2>
          <p className="text-gray-600 text-lg mb-8 leading-relaxed font-light">
            V8 engine boot times were bottlenecking agent workflows. In a CI/CD pipeline, starting up Node.js for every agent task was simply unacceptable. 
            We replaced it with a single, statically compiled Go binary. The result is pure, unadulterated speed.
          </p>
          <div className="grid grid-cols-2 gap-12 mt-12">
            <div>
              <div className="text-xs font-bold text-gray-500 uppercase tracking-widest mb-2">Previous (Node.js)</div>
              <div className="text-5xl font-black text-gray-400 font-serif opacity-50">&gt;800<span className="text-2xl">ms</span></div>
            </div>
            <div>
              <div className="text-xs font-bold text-primary uppercase tracking-widest mb-2">Current (Go)</div>
              <div className="text-5xl font-black text-primary font-serif">&lt;50<span className="text-2xl">ms</span></div>
            </div>
          </div>
        </div>
        
        <div className="lg:w-1/2 w-full glass-panel p-10 h-[450px] flex items-end justify-center gap-12 relative overflow-hidden">
          <div className="absolute inset-0 bg-[radial-gradient(#e5e5e5_1px,transparent_1px)] bg-[size:20px_20px] opacity-50"></div>
          
          <div className="flex flex-col items-center gap-4 w-1/3 relative z-10">
            <motion.div 
              initial={{ height: 0 }}
              whileInView={{ height: "100%" }}
              transition={{ duration: 1.5, ease: "easeOut" }}
              className="w-full max-w-[140px] bg-gray-200 rounded-t-xl relative shadow-inner"
            >
              <div className="absolute -top-12 w-full text-center text-gray-500 font-medium text-lg">850ms</div>
            </motion.div>
          </div>
          
          <div className="flex flex-col items-center gap-4 w-1/3 relative z-10">
            <motion.div 
              initial={{ height: 0 }}
              whileInView={{ height: "15%" }}
              transition={{ duration: 0.5, ease: "easeOut", delay: 0.5 }}
              className="w-full max-w-[140px] bg-primary rounded-t-xl relative shadow-xl"
            >
              <div className="absolute -top-12 w-full text-center text-primary font-bold text-xl">42ms</div>
            </motion.div>
          </div>
        </div>
        
      </div>
    </section>
  );
};

const ToolingSection = () => {
  return (
    <section className="py-32 px-6 md:px-16 bg-background relative">
      <div className="max-w-6xl mx-auto">
        <div className="text-center mb-24">
          <span className="font-cursive text-3xl text-accent mb-4 block">03. Developer Experience</span>
          <h2 className="text-4xl md:text-6xl font-bold text-primary">Next-Gen Tooling</h2>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12">
          {/* Semantic Diff */}
          <motion.div whileHover={{ y: -5 }} className="glass-panel p-10 group">
            <div className="w-14 h-14 rounded-2xl bg-white shadow-sm flex items-center justify-center mb-8 border border-gray-100 group-hover:scale-110 transition-transform">
              <Code2 className="w-6 h-6 text-primary" />
            </div>
            <h3 className="text-3xl font-bold text-primary mb-4">Semantic Diff</h3>
            <p className="text-gray-500 mb-8 font-light leading-relaxed">
              Standard git diffs are noisy. Gitclaw diff parses the Abstract Syntax Tree to provide clean, human-readable summaries of code mutations.
            </p>
            <div className="bg-[#fcfaf7] p-6 rounded-2xl border border-gray-200 font-mono text-sm shadow-inner relative overflow-hidden">
              <div className="absolute left-0 top-0 bottom-0 w-1 bg-accent"></div>
              <div className="text-gray-400 mb-4">$ gitclaw diff --semantic</div>
              <ul className="text-gray-700 space-y-3">
                <li className="flex gap-3"><span className="text-accent">→</span> Refactored parse function to ES6 arrow.</li>
                <li className="flex gap-3"><span className="text-accent">→</span> Upgraded variable x to const.</li>
                <li className="flex gap-3"><span className="text-accent">→</span> Logic remains completely identical.</li>
              </ul>
            </div>
          </motion.div>

          {/* Benchmark */}
          <motion.div whileHover={{ y: -5 }} className="glass-panel p-10 group">
            <div className="w-14 h-14 rounded-2xl bg-white shadow-sm flex items-center justify-center mb-8 border border-gray-100 group-hover:scale-110 transition-transform">
              <TrendingUp className="w-6 h-6 text-primary" />
            </div>
            <h3 className="text-3xl font-bold text-primary mb-4">Benchmarking</h3>
            <p className="text-gray-500 mb-8 font-light leading-relaxed">
              Evaluate agent performance across different LLMs natively. Test execution time, token usage, and accuracy entirely within the CLI.
            </p>
            
            <div className="overflow-hidden rounded-2xl border border-gray-200 bg-white">
              <table className="w-full text-sm text-left">
                <thead className="bg-secondary text-gray-600 font-medium">
                  <tr>
                    <th className="px-6 py-4 border-b border-r border-gray-100">Metric</th>
                    <th className="px-6 py-4 border-b border-gray-100">Agent V1</th>
                    <th className="px-6 py-4 border-b bg-[#fdfbf7] text-primary font-bold">Agent V2</th>
                  </tr>
                </thead>
                <tbody className="text-gray-600">
                  <tr className="border-b border-gray-100">
                    <td className="px-6 py-4 border-r border-gray-100">Tokens</td>
                    <td className="px-6 py-4 border-r border-gray-100">12,450</td>
                    <td className="px-6 py-4 font-semibold text-primary bg-[#fdfbf7]">8,100 (-34%)</td>
                  </tr>
                  <tr>
                    <td className="px-6 py-4 border-r border-gray-100">Pass Rate</td>
                    <td className="px-6 py-4 border-r border-gray-100">85%</td>
                    <td className="px-6 py-4 font-semibold text-primary bg-[#fdfbf7]">100% (+15%)</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </motion.div>
        </div>
      </div>
    </section>
  );
};

const CandidateSection = () => {
  return (
    <section className="py-32 px-6 md:px-16 bg-white border-t border-gray-100 text-center relative overflow-hidden">
      <div className="max-w-4xl mx-auto relative z-10">
        <div className="w-24 h-24 mx-auto bg-primary rounded-full flex items-center justify-center mb-8 shadow-2xl">
          <Box className="w-10 h-10 text-white" />
        </div>
        <h2 className="text-5xl font-black text-primary mb-6">Built by Shubh Shrivastava</h2>
        <p className="font-cursive text-3xl text-accent mb-12">for Lyzr</p>
        
        <p className="text-xl text-gray-500 font-light max-w-2xl mx-auto leading-relaxed mb-12">
          This entire architecture migration, the new features, and this presentation website represent my final engineering submission to join the Lyzr team. 
        </p>
        
        <a href="https://github.com/Creat1ve-shubh/gitagent" target="_blank" rel="noreferrer" className="inline-flex items-center gap-2 px-10 py-5 bg-primary text-white rounded-full font-medium hover:bg-gray-800 transition-all shadow-xl hover:shadow-2xl transform hover:-translate-y-1 text-lg">
          Review My Code
        </a>
      </div>
    </section>
  );
};

export default function Home() {
  return (
    <main className="bg-background text-text min-h-screen font-sans selection:bg-accent selection:text-white">
      <nav className="flex items-center justify-between p-6 md:px-16 absolute top-0 w-full z-50">
        <div className="text-2xl font-black tracking-tighter flex items-center gap-2 text-primary">
          <div className="w-8 h-8 bg-primary rounded-lg text-white flex items-center justify-center text-lg">G</div>
          GITCLAW<span className="text-accent text-sm font-bold align-top mt-1">V2</span>
        </div>
        <div className="hidden md:flex gap-10 text-sm font-semibold tracking-wide text-gray-500">
          <a href="#architecture" className="elegant-border hover:text-primary transition-colors py-1">Architecture</a>
          <a href="https://github.com/Creat1ve-shubh/gitagent" target="_blank" rel="noreferrer" className="elegant-border hover:text-primary transition-colors py-1">GitHub</a>
        </div>
      </nav>

      <HeroSection />
      <ArchitectureDiagram />
      <PerformanceShowdown />
      <ToolingSection />
      <CandidateSection />
      
      <footer className="py-12 bg-[#f9f7f1] text-center border-t border-gray-200">
        <p className="text-gray-400 font-medium text-sm">
          &copy; 2026 Shubh Shrivastava. Engineering Submission for Lyzr.
        </p>
      </footer>
    </main>
  );
}
