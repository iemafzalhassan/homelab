FROM oven/bun:1 AS builder
WORKDIR /app

# Copy root workspace configurations
COPY package.json bun.lock ./
COPY turbo.json ./

# Copy workspace packages
COPY packages ./packages
COPY apps/web ./apps/web

# Install dependencies (respecting workspaces)
RUN bun install

# Build the web app
WORKDIR /app/apps/web
RUN bun run build

# Runner stage
FROM node:20-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production
ENV PORT=3333
EXPOSE 3333

# We copy the standalone build output (requires output: "standalone" in next.config.ts)
COPY --from=builder /app/apps/web/public ./apps/web/public
COPY --from=builder /app/apps/web/.next/standalone ./
COPY --from=builder /app/apps/web/.next/static ./apps/web/.next/static

CMD ["node", "apps/web/server.js"]
