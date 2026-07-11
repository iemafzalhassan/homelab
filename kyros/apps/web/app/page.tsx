export default function Page() {
  return (
    <div className="flex flex-col items-center justify-center min-h-[calc(100vh-73px)] p-8">
      <div className="max-w-3xl text-center space-y-6">
        <h2 className="text-4xl md:text-6xl font-bold tracking-tighter">
          The Trusted Software <br className="hidden sm:inline" /> Supply Chain Platform
        </h2>
        <p className="text-lg text-muted-foreground md:text-xl">
          Kyros provides hardened OCI images, built from scratch, continuously scanned and signed.
        </p>
        <div className="flex justify-center gap-4 pt-4">
          <a href="/search" className="inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground hover:bg-primary/90 h-10 px-4 py-2">
            Search Images
          </a>
          <a href="/docs" className="inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 border border-input bg-background hover:bg-accent hover:text-accent-foreground h-10 px-4 py-2">
            Read Docs
          </a>
        </div>
      </div>
    </div>
  );
}
