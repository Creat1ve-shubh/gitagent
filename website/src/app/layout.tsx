import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Gitclaw V2 | Agent Runtime Reimagined in Go",
  description:
    "A masterclass in modern systems engineering. MVCC Ledger, zero-latency security pipeline, and sub-50ms cold starts.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="h-full antialiased">
      <head>
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link rel="preconnect" href="https://fonts.gstatic.com" crossOrigin="anonymous" />
        <link
          href="https://fonts.googleapis.com/css2?family=Dancing+Script:wght@500;700&family=Inter:wght@300;400;500;600;700;800;900&family=Playfair+Display:ital,wght@0,400;0,600;0,800;0,900;1,400;1,600&family=JetBrains+Mono:wght@400;500;600&display=swap"
          rel="stylesheet"
        />
      </head>
      <body className="min-h-full flex flex-col bg-background">{children}</body>
    </html>
  );
}
