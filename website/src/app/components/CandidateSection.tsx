"use client";
import React from "react";
import { motion, useInView } from "framer-motion";
import { ExternalLink } from "lucide-react";

export default function CandidateSection() {
  const ref = React.useRef(null);
  const inView = useInView(ref, { once: true, margin: "-80px" });

  return (
    <section className="py-24 md:py-32 px-6 relative" ref={ref}>
      <div className="absolute inset-0 spotlight" />
      <div className="max-w-3xl mx-auto text-center relative z-10">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={inView ? { opacity: 1, y: 0 } : {}}
          transition={{ duration: 0.6 }}
        >
          <div className="w-16 h-16 mx-auto rounded-2xl bg-accent/10 border border-accent/20 flex items-center justify-center mb-8 animate-pulse-glow">
            <span className="text-2xl font-black text-accent">G</span>
          </div>

          <h2 className="text-3xl md:text-5xl font-bold mb-3">Built by Shubh Shrivastava</h2>
          <p className="font-cursive text-2xl md:text-3xl text-accent italic mb-8">for Lyzr</p>

          <p className="text-text-muted text-sm md:text-base leading-relaxed max-w-lg mx-auto mb-10">
            This entire architecture migration, the new features, and this presentation website represent my final engineering submission to join the Lyzr team.
          </p>

          <a
            href="https://github.com/Creat1ve-shubh/gitagent"
            target="_blank"
            rel="noreferrer"
            className="btn-accent inline-flex items-center gap-2 px-8 py-4 text-base"
          >
            <ExternalLink className="w-4 h-4" /> Review My Code
          </a>
        </motion.div>
      </div>
    </section>
  );
}
