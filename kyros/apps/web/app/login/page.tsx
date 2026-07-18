import { signIn } from "@/auth";

export default function LoginPage() {
  return (
    <main className="min-h-screen w-full flex items-center justify-center bg-[hsl(240,10%,3%)] overflow-hidden relative selection:bg-cyan-500/30 text-white">
      {/* Background gradients */}
      <div className="absolute top-[-20%] left-[-10%] w-[500px] h-[500px] rounded-full bg-cyan-500/10 blur-[120px]" />
      <div className="absolute bottom-[-20%] right-[-10%] w-[600px] h-[600px] rounded-full bg-blue-500/10 blur-[150px]" />

      <div className="z-10 w-full max-w-md p-8 rounded-2xl glass border border-white/5 shadow-2xl relative overflow-hidden group">
        <div className="absolute inset-0 bg-gradient-to-br from-white/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500" />
        
        <div className="relative z-10 flex flex-col items-center">
          <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-cyan-400 to-blue-600 flex items-center justify-center mb-6 shadow-lg shadow-cyan-500/20">
            <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="text-white">
              <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
            </svg>
          </div>
          
          <h1 className="text-3xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-white to-gray-400 mb-2 font-['Outfit']">
            Welcome to Kyros
          </h1>
          <p className="text-gray-400 text-center mb-8">
            The Trusted Software Supply Chain Platform. Sign in to access your organization's dashboard.
          </p>

          <form
            action={async () => {
              "use server"
              await signIn("keycloak", { redirectTo: "/dashboard" })
            }}
            className="w-full"
          >
            <button
              type="submit"
              className="w-full py-3 px-4 bg-white/5 hover:bg-white/10 border border-white/10 rounded-xl transition-all duration-300 font-medium flex items-center justify-center space-x-2 group/btn relative overflow-hidden cursor-pointer"
            >
              <div className="absolute inset-0 bg-gradient-to-r from-cyan-500/10 to-blue-500/10 opacity-0 group-hover/btn:opacity-100 transition-opacity duration-300" />
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="text-cyan-400 relative z-10">
                <path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"/>
                <polyline points="10 17 15 12 10 7"/>
                <line x1="15" y1="12" x2="3" y2="12"/>
              </svg>
              <span className="relative z-10">Continue with SSO</span>
            </button>
          </form>

          <div className="mt-8 pt-6 border-t border-white/5 w-full text-center">
            <p className="text-sm text-gray-500">
              Authentication powered by <span className="text-gray-300">Keycloak</span>
            </p>
          </div>
        </div>
      </div>
    </main>
  );
}
