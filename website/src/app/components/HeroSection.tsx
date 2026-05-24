"use client";
import React from "react";
import { motion } from "framer-motion";
import { ArrowRight, Play } from "lucide-react";

export default function HeroSection() {
  return (
    <section className="relative min-h-screen flex flex-col items-center justify-center px-6 overflow-hidden">
      {/* Hero glow */}
      <div className="hero-glow" />

      {/* Background grid */}
      <div className="absolute inset-0 grid-bg opacity-60" />

      {/* Beam lines */}
      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-px h-40 bg-gradient-to-b from-transparent via-accent/30 to-transparent" />

      {/* Content */}
      <div className="relative z-10 text-center max-w-4xl mx-auto">
        {/* Badge */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6 }}
          className="inline-flex items-center gap-2 px-4 py-1.5 rounded-full border border-border-2 bg-surface/60 backdrop-blur-sm mb-10"
        >
          <span className="w-2 h-2 rounded-full bg-green-400 animate-pulse" />
          <span className="text-xs font-medium text-text-muted tracking-wide">v2.0 — Rewritten in Go</span>
        </motion.div>

        {/* Main title */}
        <motion.h1
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 0.1 }}
          className="text-6xl md:text-8xl lg:text-9xl font-black text-white leading-[0.9] tracking-tight mb-4"
        >
          gitagent
        </motion.h1>

        <motion.span
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 0.2 }}
          className="block font-cursive text-5xl md:text-7xl lg:text-8xl text-accent italic mb-8"
        >
          go
        </motion.span>

        {/* Subtitle */}
        <motion.p
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.4, duration: 0.8 }}
          className="text-base md:text-lg text-text-muted max-w-xl mx-auto leading-relaxed mb-10 font-light"
        >
          Turn your codebase into a self-operating system with powerful AI. Conflict-free concurrency, zero-latency security, and sub-50ms cold starts.
        </motion.p>

        {/* Terminal widget */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.5, duration: 0.7 }}
          className="max-w-lg mx-auto mb-12"
        >
          <div className="dark-card p-4 text-left">
            <div className="flex items-center gap-2 mb-3">
              <div className="w-3 h-3 rounded-full bg-[#ff5f57]" />
              <div className="w-3 h-3 rounded-full bg-[#febc2e]" />
              <div className="w-3 h-3 rounded-full bg-[#28c840]" />
              <span className="text-[10px] text-text-muted ml-2 font-mono">terminal</span>
            </div>
            <div className="font-mono text-sm text-secondary space-y-1">
              <div className="flex items-center gap-2">
                <span className="text-accent">$</span>
                <span className="text-white">gitclaw --model ollama:kimi2.5 --voice --dir ~/project</span>
              </div>
              <div className="text-text-muted text-xs mt-2 flex items-center gap-2">
                <Play className="w-3 h-3 text-green-400" />
                <span>Agent runtime started in <span className="text-accent font-semibold">42ms</span></span>
              </div>
            </div>
          </div>
        </motion.div>

        {/* CTA Buttons */}
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.7, duration: 0.6 }}
          className="flex flex-col sm:flex-row items-center justify-center gap-4"
        >
          <a href="#architecture" className="btn-accent flex items-center gap-2 px-6 py-3">
            Explore Architecture <ArrowRight className="w-4 h-4" />
          </a>
          <a href="https://github.com/Creat1ve-shubh/gitagent" target="_blank" rel="noreferrer"
            className="flex items-center gap-2 px-6 py-3 rounded-[10px] border border-border-2 text-sm font-semibold text-secondary hover:text-white hover:border-accent/40 transition-all">
            View Source Code
          </a>
        </motion.div>
      </div>

      {/* Scroll indicator */}
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 1.5 }}
        className="absolute bottom-8"
      >
        <motion.div
          animate={{ y: [0, 8, 0] }}
          transition={{ duration: 2, repeat: Infinity, ease: "easeInOut" }}
          className="w-5 h-8 rounded-full border border-border-2 flex justify-center pt-1.5"
        >
          <div className="w-1 h-2 rounded-full bg-text-muted" />
        </motion.div>
      </motion.div>
    </section>
  );
}
