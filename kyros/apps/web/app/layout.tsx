import type { Metadata } from "next";
import { Inter, Outfit } from "next/font/google";
import { ThemeProvider } from "next-themes";
import { ThemeToggle } from "../components/theme-toggle";
import "./globals.css";

const inter = Inter({ subsets: ["latin"], variable: "--font-inter" });
const outfit = Outfit({ subsets: ["latin"], variable: "--font-outfit" });

export const metadata: Metadata = {
  title: "Kyros | Secure by Default",
  description: "The Trusted Software Supply Chain Platform. Hardened OCI images, built from scratch.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning className={`${inter.variable} ${outfit.variable}`}>
      <body className={`font-sans antialiased min-h-screen flex flex-col bg-background text-foreground`}>
        <ThemeProvider
          attribute="class"
          defaultTheme="dark"
          enableSystem
          disableTransitionOnChange
        >
          {/* Glassmorphism Header */}
          <header className="sticky top-0 z-50 w-full border-b border-white/5 bg-background/60 backdrop-blur-xl supports-[backdrop-filter]:bg-background/40">
            <div className="container mx-auto flex h-16 max-w-screen-2xl items-center px-6 md:px-12">
              <div className="mr-4 flex w-full justify-between md:w-auto md:justify-start">
                <a className="mr-8 flex items-center space-x-2 group" href="/">
                  <div className="w-8 h-8 rounded-lg bg-gradient-to-tr from-primary to-accent flex items-center justify-center shadow-lg shadow-primary/20 group-hover:shadow-primary/40 transition-all duration-300">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" className="text-white">
                      <path d="M12 2L2 7L12 12L22 7L12 2Z" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                      <path d="M2 17L12 22L22 17" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                      <path d="M2 12L12 17L22 12" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                    </svg>
                  </div>
                  <span className="font-heading font-bold text-xl tracking-tight text-foreground">
                    Kyros
                  </span>
                </a>
                <nav className="hidden md:flex items-center gap-8 text-sm font-medium">
                  <a className="transition-colors hover:text-foreground text-foreground/70" href="/search">Search</a>
                  <a className="transition-colors hover:text-foreground text-foreground/70" href="/orgs">Publishers</a>
                  <a className="transition-colors hover:text-foreground text-foreground/70" href="/docs">Docs</a>
                </nav>
              </div>
              <div className="flex flex-1 items-center justify-end space-x-4">
                <nav className="flex items-center space-x-2">
                  <a href="/login" className="hidden sm:inline-flex items-center justify-center rounded-full text-sm font-medium transition-colors hover:text-foreground text-foreground/70 h-9 px-4">
                    Sign In
                  </a>
                  <a href="/signup" className="hidden sm:inline-flex items-center justify-center rounded-full text-sm font-medium bg-primary text-primary-foreground hover:bg-primary/90 h-9 px-4 shadow-[0_0_15px_rgba(var(--color-primary),0.3)]">
                    Get Started
                  </a>
                  <ThemeToggle />
                </nav>
              </div>
            </div>
          </header>
          <main className="flex-1 flex flex-col w-full">
            {children}
          </main>
        </ThemeProvider>
      </body>
    </html>
  );
}
