"use client";

import React from "react";
import Lenis from "lenis";
import HeroSection from "./components/HeroSection";
import FeaturesSection from "./components/FeaturesSection";
import ArchitectureSection from "./components/ArchitectureSection";
import PerformanceSection from "./components/PerformanceSection";
import CandidateSection from "./components/CandidateSection";

export default function Home() {
  React.useEffect(() => {
    const lenis = new Lenis({
      autoRaf: true,
      duration: 1.4,
      easing: (t) => Math.min(1, 1.001 - Math.pow(2, -10 * t)),
      smoothWheel: true,
    });
    return () => { lenis.destroy(); };
  }, []);

  return (
    <main className="bg-background text-text min-h-screen">
      {/* Navigation */}
      <nav className="fixed top-0 w-full z-50 px-6 md:px-16 py-4 bg-background/80 backdrop-blur-xl border-b border-border">
        <div className="flex items-center justify-between max-w-6xl mx-auto">
          <div className="flex items-center gap-2.5 group">
            <div className="w-7 h-7 bg-accent rounded-lg text-white flex items-center justify-center text-xs font-bold group-hover:scale-110 transition-transform">
              G
            </div>
            <span className="text-sm font-bold tracking-tight text-white">
              Gitclaw <span className="text-accent text-[10px] font-semibold ml-0.5">V2</span>
            </span>
          </div>

          <div className="hidden md:flex items-center gap-8">
            {["Features", "Architecture", "Performance"].map((item) => (
              <a key={item} href={`#${item.toLowerCase()}`}
                className="text-xs font-medium text-text-muted hover:text-white transition-colors">
                {item}
              </a>
            ))}
            <a href="https://github.com/Creat1ve-shubh/gitagent" target="_blank" rel="noreferrer"
              className="text-xs font-medium text-text-muted hover:text-white transition-colors">
              GitHub
            </a>
            <a href="#architecture" className="btn-accent text-xs px-4 py-2">
              Get Started
            </a>
          </div>
        </div>
      </nav>

      <HeroSection />

      <div id="features">
        <FeaturesSection />
      </div>

      <ArchitectureSection />
      <PerformanceSection />
      <CandidateSection />

      {/* Footer */}
      <footer className="py-8 border-t border-border text-center">
        <p className="text-text-muted text-xs">
          &copy; 2026 Shubh Shrivastava · Engineering Submission for Lyzr
        </p>
      </footer>
    </main>
  );
}
