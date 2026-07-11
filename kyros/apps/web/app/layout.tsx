import type { Metadata } from "next";
import { Inter } from "next/font/google";
import { ThemeProvider } from "next-themes";
import "./globals.css";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Kyros — The Trusted Software Supply Chain Platform",
  description: "Hardened OCI images, built from scratch.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={`${inter.className} antialiased min-h-screen flex flex-col`}>
        <ThemeProvider
          attribute="class"
          defaultTheme="dark"
          enableSystem
          disableTransitionOnChange
        >
          <header className="border-b px-6 py-4 flex items-center justify-between">
            <h1 className="text-xl font-bold tracking-tight">Kyros</h1>
            <nav className="flex gap-4">
              <a href="/search" className="text-sm text-muted-foreground hover:text-foreground transition-colors">Search</a>
              <a href="/orgs" className="text-sm text-muted-foreground hover:text-foreground transition-colors">Publishers</a>
            </nav>
          </header>
          <main className="flex-1">
            {children}
          </main>
        </ThemeProvider>
      </body>
    </html>
  );
}
