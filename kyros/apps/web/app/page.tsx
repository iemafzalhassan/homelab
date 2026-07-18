import { Shield, Search, Zap, ArrowRight, Layers, Lock, Cpu } from "lucide-react"
import { AnimatedTerminal } from "../components/AnimatedTerminal"

export default function Page() {
  return (
    <div className="flex flex-col items-center min-h-[calc(100vh-64px)] overflow-hidden">
      
      {/* Background Effects */}
      <div className="absolute inset-0 z-[-1] bg-grid-pattern opacity-20"></div>
      <div className="absolute top-0 right-1/4 w-96 h-96 bg-primary/20 rounded-full blur-[128px] -z-10 mix-blend-screen pointer-events-none"></div>
      <div className="absolute bottom-0 left-1/4 w-[30rem] h-[30rem] bg-accent/10 rounded-full blur-[128px] -z-10 mix-blend-screen pointer-events-none"></div>

      {/* Hero Section */}
      <section className="w-full flex flex-col lg:flex-row items-center justify-between pt-24 pb-32 px-6 lg:px-12 max-w-screen-2xl gap-16 relative">
        <div className="w-full lg:w-1/2 space-y-8 animate-in fade-in slide-in-from-bottom-8 duration-1000">
          <div className="inline-flex items-center rounded-full border border-primary/30 px-3 py-1 text-sm font-medium transition-colors bg-primary/10 text-primary-foreground backdrop-blur-sm shadow-[0_0_15px_rgba(var(--color-primary),0.1)]">
            <span className="flex h-2 w-2 rounded-full bg-primary mr-2 animate-pulse"></span>
            Kyros Developer Preview
          </div>
          
          <h1 className="text-5xl sm:text-6xl md:text-7xl font-extrabold tracking-tight lg:text-8xl text-foreground">
            Secure by <br className="hidden md:block"/>
            <span className="text-transparent bg-clip-text bg-gradient-to-r from-primary via-blue-400 to-accent animate-pulse-glow inline-block">default.</span>
          </h1>
          
          <p className="max-w-xl text-xl text-muted-foreground leading-relaxed">
            The Trusted Software Supply Chain Platform. Hardened OCI images, built from scratch, continuously scanned, and cryptographically signed.
          </p>
          
          <div className="flex flex-col sm:flex-row gap-4 pt-4">
            <a href="/search" className="inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium transition-all bg-primary text-primary-foreground hover:bg-primary/90 h-12 px-8 shadow-[0_0_20px_rgba(var(--color-primary),0.4)] hover:shadow-[0_0_30px_rgba(var(--color-primary),0.6)] group">
              Explore Images
              <ArrowRight className="ml-2 w-4 h-4 group-hover:translate-x-1 transition-transform" />
            </a>
            <a href="/docs" className="inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium transition-colors border border-white/10 bg-white/5 hover:bg-white/10 backdrop-blur-sm h-12 px-8 text-foreground">
              Documentation
            </a>
          </div>
        </div>

        <div className="w-full lg:w-1/2 relative animate-in fade-in slide-in-from-right-8 duration-1000 delay-200">
          <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[120%] h-[120%] bg-gradient-to-tr from-primary/20 to-accent/20 blur-[100px] -z-10 rounded-full"></div>
          <div className="animate-float">
            <AnimatedTerminal />
          </div>
        </div>
      </section>

      {/* Feature Grid */}
      <section className="w-full max-w-screen-2xl px-6 lg:px-12 py-24 relative z-10 border-t border-white/5 bg-background/50 backdrop-blur-sm">
        <div className="text-center mb-16 space-y-4">
          <h2 className="text-3xl md:text-5xl font-bold font-heading">Enterprise-grade security, <br/> built for developers.</h2>
          <p className="text-muted-foreground text-lg max-w-2xl mx-auto">Kyros eliminates the friction between shipping fast and staying secure by integrating security into the absolute base of your containers.</p>
        </div>

        <div className="grid gap-8 md:grid-cols-2 lg:grid-cols-3">
          <div className="group glass-card p-8 rounded-2xl transition-all duration-300 hover:-translate-y-2 hover:shadow-[0_0_40px_rgba(var(--color-primary),0.15)] hover:border-primary/30 relative overflow-hidden">
            <div className="absolute -right-10 -top-10 w-40 h-40 bg-primary/10 rounded-full blur-3xl group-hover:bg-primary/20 transition-all duration-500"></div>
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-primary/20 border border-primary/30 mb-6 group-hover:scale-110 transition-transform duration-300">
              <Shield className="h-6 w-6 text-primary" />
            </div>
            <h3 className="text-2xl font-bold mb-3">Zero CVE Base Images</h3>
            <p className="text-muted-foreground leading-relaxed">Our minimalist distroless base images drastically reduce attack surface, ensuring you start with exactly zero known vulnerabilities.</p>
          </div>

          <div className="group glass-card p-8 rounded-2xl transition-all duration-300 hover:-translate-y-2 hover:shadow-[0_0_40px_rgba(var(--color-accent),0.15)] hover:border-accent/30 relative overflow-hidden">
            <div className="absolute -right-10 -top-10 w-40 h-40 bg-accent/10 rounded-full blur-3xl group-hover:bg-accent/20 transition-all duration-500"></div>
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-accent/20 border border-accent/30 mb-6 group-hover:scale-110 transition-transform duration-300">
              <Search className="h-6 w-6 text-accent" />
            </div>
            <h3 className="text-2xl font-bold mb-3">Continuous Scanning</h3>
            <p className="text-muted-foreground leading-relaxed">Every image is actively monitored against the latest threat intelligence and rebuilt instantly when upstream vulnerabilities are found.</p>
          </div>

          <div className="group glass-card p-8 rounded-2xl transition-all duration-300 hover:-translate-y-2 hover:shadow-[0_0_40px_rgba(var(--color-green-500),0.15)] hover:border-green-500/30 relative overflow-hidden md:col-span-2 lg:col-span-1">
            <div className="absolute -right-10 -top-10 w-40 h-40 bg-green-500/10 rounded-full blur-3xl group-hover:bg-green-500/20 transition-all duration-500"></div>
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-green-500/20 border border-green-500/30 mb-6 group-hover:scale-110 transition-transform duration-300">
              <Layers className="h-6 w-6 text-green-400" />
            </div>
            <h3 className="text-2xl font-bold mb-3">Native SBOMs</h3>
            <p className="text-muted-foreground leading-relaxed">High fidelity Software Bill of Materials (SBOM) and cryptographically verifiable signatures natively attached to every OCI index.</p>
          </div>
        </div>
      </section>
      
    </div>
  );
}
