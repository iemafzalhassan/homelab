"use client";

import { useEffect, useState } from "react";
import { Terminal, ShieldCheck, Fingerprint, PackageCheck } from "lucide-react";

const steps = [
  { text: "FROM cgr.dev/chainguard/static:latest", icon: Terminal, color: "text-muted-foreground" },
  { text: "Building image kyros-registry:v1.0...", icon: PackageCheck, color: "text-primary" },
  { text: "Scanning for vulnerabilities...", icon: ShieldCheck, color: "text-accent" },
  { text: "Result: 0 CVEs detected. Image is secure.", icon: ShieldCheck, color: "text-green-400" },
  { text: "Generating high-fidelity SBOM...", icon: PackageCheck, color: "text-muted-foreground" },
  { text: "Signing image with Cosign...", icon: Fingerprint, color: "text-accent" },
  { text: "Successfully pushed to registry.kyros.io", icon: Terminal, color: "text-primary" },
];

export function AnimatedTerminal() {
  const [currentStep, setCurrentStep] = useState(0);

  useEffect(() => {
    const timer = setInterval(() => {
      setCurrentStep((prev) => (prev < steps.length - 1 ? prev + 1 : 0));
    }, 2000); // 2s per step
    return () => clearInterval(timer);
  }, []);

  return (
    <div className="w-full max-w-lg mx-auto relative group">
      {/* Glow effect */}
      <div className="absolute -inset-0.5 bg-gradient-to-r from-primary to-accent rounded-xl blur opacity-30 group-hover:opacity-70 transition duration-1000 group-hover:duration-200 animate-pulse-glow"></div>
      
      {/* Terminal window */}
      <div className="relative rounded-xl bg-[#0d1117] border border-white/10 overflow-hidden shadow-2xl flex flex-col h-[320px]">
        
        {/* Terminal Header */}
        <div className="flex items-center px-4 py-3 bg-[#161b22] border-b border-white/5">
          <div className="flex space-x-2">
            <div className="w-3 h-3 rounded-full bg-red-500/80"></div>
            <div className="w-3 h-3 rounded-full bg-yellow-500/80"></div>
            <div className="w-3 h-3 rounded-full bg-green-500/80"></div>
          </div>
          <div className="mx-auto text-xs font-mono text-muted-foreground">kyros-builder ~ /app</div>
        </div>
        
        {/* Terminal Body */}
        <div className="p-5 font-mono text-sm space-y-4 overflow-y-auto flex-1">
          {steps.map((step, index) => {
            const Icon = step.icon;
            const isVisible = index <= currentStep;
            const isCurrent = index === currentStep;
            
            return (
              <div 
                key={index} 
                className={`flex items-start space-x-3 transition-all duration-300 ${isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4 hidden'}`}
              >
                <div className={`mt-0.5 ${step.color} ${isCurrent ? 'animate-pulse' : 'opacity-70'}`}>
                  {index === 0 ? (
                    <span className="text-primary font-bold">❯</span>
                  ) : (
                    <Icon className="w-4 h-4" />
                  )}
                </div>
                <div className={`${step.color} ${isCurrent ? '' : 'opacity-70'} font-medium`}>
                  {step.text}
                </div>
              </div>
            );
          })}
          {currentStep < steps.length - 1 && (
            <div className="flex items-start space-x-3 opacity-50 animate-pulse">
              <span className="text-primary font-bold mt-0.5">❯</span>
              <span className="w-2 h-4 bg-primary inline-block"></span>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
