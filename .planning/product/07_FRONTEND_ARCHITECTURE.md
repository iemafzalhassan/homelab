# Kyros Frontend Architecture

## Overview
This document details the frontend architecture of the Kyros platform, covering routing, layouts, state management, data fetching, caching, real-time updates, permissions, and feature flags. The frontend is built with Next.js 15, React 18, TypeScript, and Tailwind CSS, following modern best practices for scalability, maintainability, and performance.

## Technology Stack
- **Framework**: Next.js 15 (App Router)
- **Language**: TypeScript
- **UI Library**: Custom component library based on Radix UI primitives
- **Styling**: Tailwind CSS
- **State Management**: 
  - Server state: React Query (TanStack Query)
  - Client state: Zustand (for lightweight global state) and React Context (for theme, auth, etc.)
- **Data Fetching**: React Query with SWR (stale-while-revalidate) caching
- **Real-time Updates**: WebSocket connections via a custom hook
- **API Communication**: REST API with JSON payloads
- **Authentication**: NextAuth.js with custom providers for Kyros-specific flows
- **Testing**: Jest and React Testing Library for unit tests, Playwright for end-to-end tests
- **Code Quality**: ESLint, Prettier, TypeScript strict mode

## Routing

### App Router (Next.js 15)
Kyros uses the Next.js 15 App Router for file-system based routing with support for layouts, nested routes, and route groups.

#### Route Structure
```
/app
  /(app)                 # Route group for authenticated sections
    /dashboard/page.tsx  # /dashboard
    /repositories
      /page.tsx          # /repositories
      /[namespace]/[name]/page.tsx  # /repositories/{namespace}/{name}
    /images
      /page.tsx          # /images
      /[digest]/page.tsx # /images/{digest}
    /trust-score/page.tsx # /trust-score
    /security/page.tsx   # /security
    /users
      /page.tsx          # /users
      /[userId]/page.tsx # /users/{userId}
    /settings/page.tsx   # /settings
  /(auth)                # Route group for public/authenticated sections
    /login/page.tsx      # /login
    /register/page.tsx   # /register
    /reset-password/page.tsx # /reset-password
  /api                   # API routes (proxy to backend or custom endpoints)
    /trpc/[...].ts       # tRPC endpoints (if used)
    /auth/[...].ts       # NextAuth endpoints
  /layout.tsx            # Root layout
  /loading.tsx           # Global loading boundary
  /error.tsx             # Global error boundary
  /not-found.tsx         # 404 page
```

### Route Groups
- **(app)**: Authenticated routes requiring user session
- **(auth)**: Public routes (login, register) and auth-related routes
- **api**: Backend API routes and custom endpoint handlers

### Navigation
- **Client-side Navigation**: Using `next/link` for SPA-like transitions
- **Programmatic Navigation**: Using `useRouter()` hook from `next/navigation`
- **External Links**: Using `<a>` tags with `target="_blank"` and `rel="noopener noreferrer"`
- **Route Prefetching**: Automatic prefetching of linked routes in viewport

### Route Protection
- **Middleware**: Custom middleware to protect routes in `(app)` group
- **Route-level Protection**: Using `getServerSession` from NextAuth in layout or page components
- **Redirects**: Unauthenticated users redirected to login with `callbackUrl`

### URL Structure
- **RESTful**: Resource-based URLs (e.g., `/repositories/{namespace}/{name}`)
- **Query Parameters**: Used for filtering, sorting, pagination (e.g., `?limit=20&cursor=abc123`)
- **Hash Fragments**: Used for client-side state (e.g., `#tab=tags`)
- **Canonical URLs**: Preferred URL for each resource to prevent duplication
- **Redirects**: Proper 301/302 redirects for moved or renamed resources

## Layouts

### Root Layout
Defines the basic HTML structure, metadata, and global providers.

```tsx
// app/layout.tsx
import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import { KyrosProvider } from '@/providers/KyrosProvider';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { WebsocketProvider } from '@/providers/WebsocketProvider';

const inter = Inter({ subsets: ['latin'] });

export const metadata: Metadata = {
  title: 'Kyros',
  description: 'Open-source cloud-native software supply chain platform',
};

export const defaultMetadata: Metadata = {
  title: 'Kyros - Dashboard',
  description: 'Overview of system status and activity',
};

export const viewport = {
  width: 'device-width',
  initialScale: 1,
  maximumScale: 5,
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: 1000 * 60 * 5, // 5 minutes
        cacheTime: 1000 * 60 * 60, // 1 hour
        refetchOnWindowFocus: false,
        retry: 1,
      },
    },
  });

  return (
    <html lang="en" suppressHydrationWarning>
      <body className={inter.className}>
        <KyrosProvider>
          <QueryClientProvider client={queryClient}>
            <WebsocketProvider>
              {children}
            </WebsocketProvider>
          </QueryClientProvider>
        </KyrosProvider>
      </body>
    </html>
  );
}
```

### Authenticated Layout
Layout for routes requiring authentication.

```tsx
// app/(app)/layout.tsx
import { Sidebar } from '@/components/ui/sidebar';
import { Header } from '@/components/ui/header';
import { Footer } from '@/components/ui/footer';

export default function AuthenticatedLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <>
      <Header fixed shadow background />
      <div className="flex h-[calc(100vh-64px)]">
        <Sidebar />
        <main className="flex-1 p-6 overflow-y-auto">
          {children}
        </main>
      </div>
      <Footer borderTop background />
    </>
  );
}
```

### Public Layout
Layout for public routes (login, register, etc.).

```tsx
// app/(auth)/layout.tsx
import { KyrosLogo } from '@/components/ui/logo';

export default function PublicLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="w-full max-w-md space-y-6">
        <KyrosLogo className="h-12 w-12 mx-auto mb-4" />
        {children}
      </div>
    </div>
  );
}
```

### Page Layouts
Individual pages may have specific layouts (e.g., dashboard with grid, detail pages with tabs).

#### Dashboard Layout
```tsx
// app/(app)/dashboard/page.tsx
export default function Dashboard() {
  return (
    <div className="space-y-6">
      {/* Dashboard content */}
    </div>
  );
}
```

#### Repository Detail Layout
```tsx
// app/(app)/repositories/[namespace]/[name]/page.tsx
import { RepositoryTabs } from '@/components/repository/RepositoryTabs';

export default function RepositoryPage({
  params: { namespace, name }
}: { 
  params: { 
    namespace: string; 
    name: string 
  } 
}) {
  const repository = `${namespace}/${name}`;
  
  return (
    <div className="space-y-6">
      <div className="flex flex-col lg:flex-row items-start justify-between gap-4">
        {/* Header content */}
      </div>
      
      <RepositoryTabs repository={repository} />
    </div>
  );
}
```

## State Management

### Server State (React Query)
Used for data that originates from the server and requires caching, background updates, and synchronization.

#### Query Client Configuration
```tsx
// In root layout (see above)
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60 * 5, // 5 minutes
      cacheTime: 1000 * 60 * 60, // 1 hour
      refetchOnWindowFocus: false,
      retry: 1,
    },
  },
});
```

#### Common Query Patterns
- **Fetching Lists**: `useQueries` for parallel queries
- **Fetching Details**: `useQuery` with query key including ID
- **Mutations**: `useMutation` for create/update/delete operations
- **Pagination**: `useInfiniteQuery` for cursor-based pagination
- **Background Refetching**: Configured via `refetchInterval` or `refetchOnReconnect`

#### Example: Fetching Repositories
```tsx
import { useQuery } from '@tanstack/react-query';
import { fetchRepositories } from '@/api/repositories';

export function useRepositories(namespaceId?: string) {
  return useQuery({
    queryKey: ['repositories', namespaceId],
    queryFn: () => fetchRepositories(namespaceId),
    staleTime: 1000 * 60 * 2, // 2 minutes
  });
}
```

#### Example: Mutating Repository
```tsx
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { updateRepository } from '@/api/repositories';
import { toast } from '@/components/ui/toast';

export function useUpdateRepository() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: updateRepository,
    onSuccess: (data, variables) => {
      // Invalidate and refetch
      queryClient.invalidateQueries({ queryKey: ['repositories'] });
      queryClient.invalidateQueries({ queryKey: ['repository', variables.id] });
      toast.success('Repository updated');
    },
    onError: (error) => {
      toast.error(`Failed to update repository: ${error.message}`);
    },
  });
}
```

### Client State (Zustand & Context)
Used for UI state that doesn't require server synchronization.

#### KyrosProvider (Context)
Provides global state like user, theme, and authentication status.

```tsx
// providers/KyrosProvider.tsx
import { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { getSession, signIn, signOut } from 'next-auth/react';

type KyrosContextType = {
  user: User | null;
  theme: 'light' | 'dark';
  setTheme: (theme: 'light' | 'dark') => void;
  isLoading: boolean;
  login: (credentials: Credentials) => Promise<void>;
  logout: () => Promise<void>;
};

const KyrosContext = createContext<KyrosContextType | undefined>(undefined);

export function KyrosProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [theme, setTheme] = useState<'light' | 'dark'>('light');
  const [isLoading, setIsLoading] = useState<boolean>(false);

  useEffect(() => {
    const getUser = async () => {
      const session = await getSession();
      setUser(session?.user ?? null);
    };
    getUser();
  }, []);

  const login = async (credentials: Credentials) => {
    setIsLoading(true);
    try {
      await signIn('credentials', credentials);
      const session = await getSession();
      setUser(session?.user ?? null);
    } finally {
      setIsLoading(false);
    }
  };

  const logout = async () => {
    setIsLoading(true);
    try {
      await signOut({ redirect: false });
      setUser(null);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <KyrosContext.Provider value={{ user, theme, setTheme, isLoading, login, logout }}>
      {children}
    </KyrosContext.Provider>
  );
}

export function useKyros() {
  const context = useContext(KyrosContext);
  if (context === undefined) {
    throw new Error('useKyros must be used within a KyrosProvider');
  }
  return context;
}
```

#### Zustand Store (for lightweight global state)
Used for non-persistent UI state like sidebar collapse, modal states, etc.

```tsx
// store/useUiStore.ts
import { create } from 'zustand';

interface UiState {
  sidebarCollapsed: boolean;
  setSidebarCollapsed: (collapsed: boolean) => void;
  modalOpen: boolean;
  setModalOpen: (open: boolean) => void;
  // ... other UI state
}

export const useUiStore = create<UiState>((set) => ({
  sidebarCollapsed: false,
  setSidebarCollapsed: (collapsed) => set({ sidebarCollapsed: collapsed }),
  modalOpen: false,
  setModalOpen: (open) => set({ modalOpen: open }),
  // ...
}));
```

## Data Fetching

### API Client
Centralized API client for consistent request handling, error management, and interceptors.

```tsx
// lib/api.ts
import axios from 'axios';

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('accessToken');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor
api.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    // Handle common errors (e.g., 401, 403)
    if (error.response?.status === 401) {
      // Trigger logout or redirect to login
    }
    return Promise.reject(error);
  }
);

export default api;
```

### React Query Integration
React Query is used for data fetching with automatic caching, background updates, and deduplication.

#### Query Keys
Consistent query key structure for easy invalidation and refetching.

```tsx
// Query key examples
['repositories']                    // List of repositories
['repository', repositoryId]        // Single repository
['repository', repositoryId, 'tags'] // Tags for a repository
['trust-scores']                    // List of trust scores
['trust-score', artifactId]         // Trust score for an artifact
['users']                           // List of users
['user', userId]                    // Single user
```

#### Mutations
Encapsulate create/update/delete operations with automatic cache updates.

```tsx
// Example: Creating a repository
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { createRepository } from '@/api/repositories';

export function useCreateRepository() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: createRepository,
    onSuccess: (data, variables) => {
      // Invalidate repository list to refetch
      queryClient.invalidateQueries({ queryKey: ['repositories'] });
      // Optionally, update the cache optimistically
    },
    onError: (error) => {
      // Handle error (show toast, etc.)
    },
  });
}
```

### Data Transformation
Transform API responses to match UI requirements, keeping components decoupled from API shape.

```tsx
// lib/transformers.ts
export function transformRepository(apiRepository: ApiRepository): Repository {
  return {
    id: apiRepository.id,
    name: apiRepository.name,
    namespace: apiRepository.namespace,
    description: apiRepository.description ?? '',
    visibility: apiRepository.visibility,
    createdAt: new Date(apiRepository.createdAt),
    updatedAt: new Date(apiRepository.updatedAt),
    // ... other transformations
  };
}
```

## Caching

### React Query Caching
Built-in caching with stale-while-revalidate strategy.

#### Cache Configuration
- **staleTime**: How long data is considered fresh (default: 5 minutes)
- **cacheTime**: How long unused data is kept in memory (default: 1 hour)
- **garbageCollection**: Inactive queries are garbage collected after cacheTime

#### Cache Sharing
Automatic deduplication: identical query keys share the same cache.

#### Manual Cache Updates
- `queryClient.setQueryData`: Update cache directly
- `queryClient.invalidateQueries`: Mark queries as stale for refetch
- `queryClient.resetQueries`: Reset queries to initial state

#### Example: Optimistic Update
```tsx
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { updateRepository } from '@/api/repositories';

export function useUpdateRepository() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: updateRepository,
    onMutate: async (variables) => {
      // Cancel any outgoing refetches
      await queryClient.cancelQueries({ queryKey: ['repository', variables.id] });
      
      // Snapshot the previous value
      const previousRepository = queryClient.getQueryData<Repository>([
        'repository',
        variables.id,
      ]);
      
      // Optimistically update to the new value
      queryClient.setQueryData<Repository>(
        ['repository', variables.id],
        (old) => ({ ...old, ...variables })
      );
      
      // Return a context object with the snapshot
      return { previousRepository };
    },
    onError: (err, variables, context) => {
      // Rollback to the previous value
      if (context?.previousRepository) {
        queryClient.setQueryData(
          ['repository', variables.id],
          context.previousRepository
        );
      }
    },
    onSettled: () => {
      // Refetch regardless of error or success
      queryClient.invalidateQueries({ queryKey: ['repository', variables.id] });
    },
  });
}
```

### Additional Caching Layers
#### Service Workers
- **Workbox**: For caching static assets and API responses (optional)
- **Network First**: For API requests to ensure fresh data
- **Cache First**: For static assets (images, fonts, etc.)

#### Local Storage
- **Authentication Tokens**: Stored securely (httpOnly cookies preferred, but localStorage fallback)
- **User Preferences**: Theme, sidebar state, etc.
- **API Responses**: Not recommended for sensitive data; use with caution

#### IndexedDB
- **Large Data Sets**: For offline capabilities (future enhancement)
- **Cached SBOMs**: Large SBOM documents
- **Scan Results**: Vulnerability scan results for offline viewing

## Real-time Updates

### WebSocket Connection
Persistent WebSocket connection for real-time event updates from the backend.

#### WebsocketProvider
Provides WebSocket connection and event handling across the app.

```tsx
// providers/WebsocketProvider.tsx
import { createContext, useContext, useEffect, useState, ReactNode } from 'react';
import { io, Socket } from 'socket.io-client';

type WebsocketContextType = {
  socket: Socket | null;
  connected: boolean;
  // Event listeners can be added as needed
};

const WebsocketContext = createContext<WebsocketContextType | undefined>(undefined);

export function WebsocketProvider({ children }: { children: ReactNode }) {
  const [socket, setSocket] = useState<Socket | null>(null);
  const [connected, setConnected] = useState<boolean>(false);

  useEffect(() => {
    const newSocket = io(process.env.NEXT_PUBLIC_WS_URL, {
      transports: ['websocket'],
      auth: {
        token: localStorage.getItem('accessToken'),
      },
    });

    newSocket.on('connect', () => {
      setConnected(true);
    });

    newSocket.on('disconnect', () => {
      setConnected(false);
    });

    // Handle specific events
    newSocket.on('trust-score.updated', (data) => {
      // Invalidate relevant queries
      // Example: queryClient.invalidateQueries({ queryKey: ['trust-score', data.artifactId] });
    });

    newSocket.on('artifact.pushed', (data) => {
      // Invalidate repository lists, etc.
    });

    setSocket(newSocket);

    return () => {
      newSocket.disconnect();
    };
  }, []);

  return (
    <WebsocketContext.Provider value={{ socket, connected }}>
      {children}
    </WebsocketContext.Provider>
  );
}

export function useWebsocket() {
  const context = useContext(WebsocketContext);
  if (context === undefined) {
    throw new Error('useWebsocket must be used within a WebsocketProvider');
  }
  return context;
}
```

#### Custom Hook for Real-time Data
Hook to subscribe to real-time updates and trigger refetches.

```tsx
// hooks/useRealtimeRefetch.ts
import { useEffect } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { useWebsocket } from '@/providers/WebsocketProvider';

export function useRealtimeRefetch(queryKey: unknown[], eventName: string) {
  const queryClient = useQueryClient();
  const { socket } = useWebsocket();

  useEffect(() => {
    if (!socket) return;

    const handler = (data: unknown) => {
      // Invalidate the query to trigger refetch
      queryClient.invalidateQueries({ queryKey });
    };

    socket.on(eventName, handler);

    return () => {
      socket.off(eventName, handler);
    };
  }, [queryClient, socket, queryKey, eventName]);
}
```

#### Usage in Components
```tsx
import { useQuery } from '@tanstack/react-query';
import { useRealtimeRefetch } from '@/hooks/useRealtimeRefetch';
import { fetchTrustScore } from '@/api/trust-score';

export function useTrustScore(artifactId: string) {
  const { data, isLoading, error } = useQuery({
    queryKey: ['trust-score', artifactId],
    queryFn: () => fetchTrustScore(artifactId),
  });

  // Refetch when trust score is updated via WebSocket
  useRealtimeRefetch(['trust-score', artifactId], 'trust-score.updated');

  return { data, isLoading, error };
}
```

### Server-Sent Events (SSE) Alternative
For simpler use cases, SSE can be used instead of WebSockets.

```tsx
// Example: Using EventSource for trust score updates
import { useEffect } from 'react';
import { useQueryClient } from '@tanstack/react-query';

export function useSseRefetch(queryKey: unknown[], url: string) {
  const queryClient = useQueryClient();

  useEffect(() => {
    const eventSource = new EventSource(url);

    eventSource.onmessage = (event) => {
      const data = JSON.parse(event.data);
      // Invalidate relevant queries
      queryClient.invalidateQueries({ queryKey: ['trust-score', data.artifactId] });
    };

    return () => {
      eventSource.close();
    };
  }, [queryClient, url]);
}
```

## Permissions

### Permission Model
Kyros uses a combination of RBAC (Role-Based Access Control) and PBAC (Policy-Based Access Control) via OPA.

#### Permission Checking
Centralized permission checking utility.

```tsx
// lib/permissions.ts
import { useKyros } from '@/providers/KyrosProvider';

export function usePermissions() {
  const { user } = useKyros();
  
  // In a real implementation, this would fetch user permissions from context or API
  // For now, we'll mock based on user role
  const permissions = user ? getPermissionsForUser(user) : [];
  
  return {
    hasPermission: (permission: string) => permissions.includes(permission),
    permissions,
  };
}

function getPermissionsForUser(user: User): string[] {
  // Mock implementation - replace with actual permission lookup
  const rolePermissions: Record<string, string[]> = {
    admin: ['*'], // All permissions
    'platform-admin': [
      'repository:create',
      'repository:read',
      'repository:update',
      'repository:delete',
      'namespace:create',
      'namespace:read',
      'namespace:update',
      'namespace:delete',
      'user:read',
      'group:read',
      'role:read',
      'policy:read',
    ],
    'tenant-admin': [
      'repository:create',
      'repository:read',
      'repository:update',
      'namespace:create',
      'namespace:read',
      'namespace:update',
      'user:create',
      'user:read',
      'user:update',
      'group:create',
      'group:read',
      'group:update',
    ],
    developer: [
      'repository:read',
      'repository:push',
      'repository:pull',
      'repository:tag:create',
      'repository:tag:delete',
      'trust:score:read',
    ],
    viewer: [
      'repository:read',
      'trust:score:read',
    ],
    scanner: [
      'trust:score:read',
      'trust:vulnerability:read',
      'trust:sbom:read',
      'trust:signature:read',
    ],
  };
  
  return rolePermissions[user.role] || [];
}
```

#### Permission Guard Component
Component to conditionally render content based on permissions.

```tsx
// components/ui/PermissionGuard.tsx
import { usePermissions } from '@/lib/permissions';

interface PermissionGuardProps {
  permission: string;
  fallback?: React.ReactNode;
  children: React.ReactNode;
}

export function PermissionGuard({ 
  permission, 
  fallback, 
  children 
}: PermissionGuardProps) {
  const { hasPermission } = usePermissions();
  
  if (!hasPermission(permission)) {
    return fallback ?? null;
  }
  
  return children;
}
```

#### Usage in Pages and Components
```tsx
import { PermissionGuard } from '@/components/ui/PermissionGuard';
import { Button } from '@/components/ui/button';

export function RepositoryActions({ repositoryId }: { repositoryId: string }) {
  return (
    <div className="flex space-x-3">
      <PermissionGuard permission="repository:push">
        <Button variant="default">
          Push Image
        </Button>
      </PermissionGuard>
      
      <PermissionGuard permission="repository:scan">
        <Button variant="outline">
          Scan for Vulnerabilities
        </Button>
      </PermissionGuard>
      
      <PermissionGuard permission="repository:delete">
        <Button variant="destructive">
          Delete Repository
        </Button>
      </PermissionGuard>
    </div>
  );
}
```

#### API-level Permission Checking
Permissions are also enforced on the backend, but frontend checks prevent unnecessary requests.

```tsx
// Example: API client with permission check (optional)
import { usePermissions } from '@/lib/permissions';
import api from '@/lib/api';

export function protectedApiRequest<T>(url: string, options: RequestInit = {}): Promise<T> {
  const { hasPermission } = usePermissions();
  
  // This is a simplified example - in practice, you'd check specific permissions per endpoint
  // For demonstration, we assume all endpoints require 'api:access' permission
  if (!hasPermission('api:access')) {
    return Promise.reject(new Error('Insufficient permissions'));
  }
  
  return api(url, options).then((res) => res.data);
}
```

## Feature Flags

### Purpose
Feature flags allow for gradual rollout, A/B testing, and killing features without deployment.

### Implementation
Using a simple flag service with optional integration with a third-party service (LaunchDarkly, Unleash, etc.).

#### Feature Flag Service
```tsx
// lib/featureFlags.ts
type FeatureFlag = 
  | 'new-dashboard'
  | 'advanced-trust-score'
  | 'webhook-builder-v2'
  | 'policy-builder-v2'
  | 'audit-viewer-enhancements'
  | 'real-time-collaboration'
  | 'ai-powered-insights';

class FeatureFlagService {
  private flags: Record<FeatureFlag, boolean> = {};
  
  // Initialize with default values (could be fetched from API)
  constructor() {
    // In production, this might fetch from a service
    this.flags = {
      'new-dashboard': false,
      'advanced-trust-score': false,
      'webhook-builder-v2': false,
      'policy-builder-v2': false,
      'audit-viewer-enhancements': false,
      'real-time-collaboration': false,
      'ai-powered-insights': false,
    };
  }
  
  isEnabled(flag: FeatureFlag): boolean {
    return this.flags[flag] ?? false;
  }
  
  // Method to update flags (e.g., from API or webhook)
  updateFlags(newFlags: Partial<Record<FeatureFlag, boolean>>) {
    this.flags = { ...this.flags, ...newFlags };
  }
}

export const featureFlagService = new FeatureFlagService();

// Hook for easy access
export function useFeatureFlag(flag: FeatureFlag) {
  // In a real app, you might subscribe to flag changes
  return featureFlagService.isEnabled(flag);
}
```

#### Usage in Components
```tsx
import { useFeatureFlag } from '@/lib/featureFlags';
import { NewDashboard } from '@/components/dashboard/NewDashboard';
import { OldDashboard } from '@/components/dashboard/OldDashboard';

export function Dashboard() {
  const isNewDashboardEnabled = useFeatureFlag('new-dashboard');
  
  return isNewDashboardEnabled ? <NewDashboard /> : <OldDashboard />;
}
```

#### Usage in Routes
```tsx
// app/(app)/dashboard/page.tsx
import { useFeatureFlag } from '@/lib/featureFlags';
import NewDashboard from '@/components/dashboard/NewDashboard';
import OldDashboard from '@/components/dashboard/OldDashboard';

export default function Dashboard() {
  const isNewDashboardEnabled = useFeatureFlag('new-dashboard');
  
  return isNewDashboardEnabled ? <NewDashboard /> : <OldDashboard />;
}
```

### Feature Flag Lifecycle
1. **Development**: Flag defaults to `false`
2. **Testing**: Enable flag in staging environment for QA
3. **Rollout**: Gradually enable for percentage of users
4. **Validation**: Monitor metrics and user feedback
5. **Completion**: Flag set to `true` for all users, old code removed
6. **Removal**: Flag removed after cleanup period

### Configuration Sources
- **Default Values**: Hardcoded defaults in service
- **API**: Fetch flags from `/api/v1/feature-flags` on app load
- **WebSocket**: Real-time flag updates via events
- **Local Storage**: Persist flag overrides for debugging
- **Environment Variables**: Override flags in different environments (dev, staging, prod)

#### Example: Fetching Flags from API
```tsx
// lib/featureFlags.ts (enhanced)
import { useEffect } from 'react';
import { fetchFeatureFlags } from '@/api/featureFlags';

export function useFeatureFlags() {
  const [flags, setFlags] = useState<Record<FeatureFlag, boolean>>({});
  
  useEffect(() => {
    async function loadFlags() {
      const fetchedFlags = await fetchFeatureFlags();
      setFlags(fetchedFlags);
    }
    
    loadFlags();
  }, []);
  
  // Return a memoized service with current flags
  const service = useMemo(() => {
    return new FeatureFlagService(flags);
  }, [flags]);
  
  return service;
}
```

## Error Boundaries

### Global Error Boundary
Catches unexpected errors in the React tree and displays a fallback UI.

```tsx
// app/error.tsx
import { useEffect, useState } from 'react';

export default function Error({ 
  error, 
  resetErrorBoundary 
}: { 
  error: Error & { digest?: string }; 
  resetErrorBoundary: () => void; 
}) {
  useEffect(() => {
    // Log error to monitoring service
    console.error('ErrorBoundary caught:', error);
    // Optionally send to error tracking service (Sentry, etc.)
  }, [error]);
  
  return (
    <div className="flex min-h-[calc(100vh-64px)] items-center justify-center px-6">
      <div className="w-full max-w-2xl space-y-6 text-center">
        <h1 className="text-3xl font-bold text-gray-900">
          Something went wrong
        </h1>
        <p className="text-lg text-gray-600">
          We've encountered an unexpected error. Please try again later.
        </p>
        <div className="flex flex-col sm:flex-row gap-4 justify-center">
          <button 
            onClick={resetErrorBoundary} 
            className="btn btn-outline"
          >
            Try again
          </button>
          <button 
            onClick={() => window.location.href = '/'} 
            className="btn btn-secondary"
          >
            Go home
          </button>
        </div>
      </div>
    </div>
  );
}
```

### Component-level Error Boundaries
For isolating errors in specific components.

```tsx
// components/ui/ErrorBoundary.tsx
import { useState } from 'react';

interface ErrorBoundaryProps {
  fallback: React.ReactNode;
  children: React.ReactNode;
}

export function ErrorBoundary({ 
  fallback, 
  children 
}: ErrorBoundaryProps) {
  const [hasError, setHasError] = useState(false);
  
  if (hasError) {
    return fallback;
  }
  
  return (
    <ErrorBoundary 
      onError={() => setHasError(true)}
    >
      {children}
    </ErrorBoundary>
  );
}
```

#### Usage
```tsx
import { ErrorBoundary } from '@/components/ui/ErrorBoundary';
import { TrustScoreChart } from '@/components/trust-score/TrustScoreChart';

export function TrustScoreSection() {
  return (
    <ErrorBoundary fallback={<div>Failed to load trust score chart</div>}>
      <TrustScoreChart />
    </ErrorBoundary>
  );
}
```

## Performance Optimization

### Code Splitting
- **Route-based Splitting**: Automatic with Next.js App Router
- **Component-based Splitting**: Dynamic import for heavy components
```tsx
const HeavyComponent = dynamic(() => import('@/components/HeavyComponent'), {
  loading: () => <Spinner />,
  ssr: false,
});
```

### Image Optimization
- **Next.js Image**: Automatic optimization, resizing, and format selection
- **Placeholder**: Blurred placeholder (LQIP) or dominant color
- **Priority**: Mark above-the-fold images as `priority`

### Bundle Analysis
- **Next.js Build Analyzer**: `next build && next-analyzer`
- **Webpack Bundle Analyzer**: For deeper analysis
- **Budget**: Set performance budgets for JavaScript and CSS

### Rendering Optimization
- **React.memo**: For components that render frequently with same props
- **useMemo**: For expensive computations
- **useCallback**: For functions passed as props
- **Virtual Scrolling**: For large lists (using `react-window` or similar)
- **Intersection Observer**: For lazy loading of off-screen content

### Network Optimization
- **Prefetching**: Next.js automatically prefetches linked routes
- **Data Prefetching**: Prefetch queries for predicted navigation
- **HTTP/2**: Leverage multiplexing for multiple requests
- **Compression**: Enable gzip/brotli on server
- **Caching**: Proper cache headers for API responses and static assets

## Accessibility

### WCAG 2.1 AA Compliance
Kyros aims for WCAG 2.1 AA compliance through:

#### Semantic HTML
- Proper use of header elements (h1-h6)
- Lists for navigation (ul/ol/li)
- Buttons for actions (not divs)
- Labels for form inputs
- Landmark elements (header, nav, main, footer)

#### Keyboard Navigation
- All interactive elements accessible via Tab
- Visible focus indicators (minimum 3:1 contrast)
- Logical tab order following visual flow
- Skip navigation links ("Skip to main content")

#### Color Contrast
- Minimum 4.5:1 for normal text
- Minimum 3:1 for large text and graphical objects
- Text over images uses backdrop or overlay for contrast
- Color not sole means of conveying information

#### Screen Reader Support
- ARIA labels and roles where needed
- Live regions for dynamic content updates
- Proper labeling of form fields and error messages
- Descriptive alt text for meaningful images
- Language attributes for multilingual content

#### Responsive Design
- Content reflows at different breakpoints
- Touch targets minimum 44x44px
- Navigation adapts to screen size (collapsible sidebar, etc.)
- Text remains readable without zooming

#### Testing
- **Automated**: axe-core integration in tests
- **Manual**: Screen reader testing (VoiceOver, NVDA, JAWS)
- **User Testing**: With participants of diverse abilities

## Internationalization (i18n)

### Strategy
Kyros uses next-i18next for internationalization.

#### Configuration
```tsx
// i18n.ts
import nextI18Next from 'next-i18next';

const nextI18NextInstance = nextI18Next({
  defaultLanguage: 'en',
  languages: ['en', 'es', 'fr', 'de', 'ja', 'pt', 'zh'],
  localePath: typeof window === 'undefined' ? './public/locales' : '/locales',
});

export default nextI18NextInstance;
export const { 
  useTranslation, 
  getTranslation, 
  router, 
  i18n, 
  Link 
} = nextI18NextInstance;
```

#### Usage
```tsx
import { useTranslation } from 'next-i18next';
import { Button } from '@/components/ui/button';

export function SubmitButton() {
  const { t } = useTranslation('common');
  
  return (
    <Button variant="default" type="submit">
      {t('submit')}
    </Button>
  );
}
```

#### Language Switcher
```tsx
import { useTranslation } from 'next-i18next';
import { useRouter } from 'next/navigation';
import { Button } from '@/components/ui/button';

export function LanguageSwitcher() {
  const { i18n } = useTranslation();
  const router = useRouter();
  
  const changeLanguage = (lng: string) => {
    i18n.changeLanguage(lng);
    router.refresh();
  };
  
  return (
    <div className="flex space-x-2">
      {['en', 'es', 'fr', 'de', 'ja', 'pt', 'zh'].map((lng) => (
        <Button 
          key={lng} 
          variant="outline" 
          size="sm" 
          onClick={() => changeLanguage(lng)}
        >
          {lng.toUpperCase()}
        </Button>
      ))}
    </div>
  );
}
```

## Testing Strategy

### Unit Tests
- **Framework**: Jest with React Testing Library
- **Coverage**: 80%+ target for components and hooks
- **Examples**:
  - Rendering
  - User interactions
  - State changes
  - Edge cases

#### Example: Button Test
```tsx
// __tests__/Button.test.tsx
import { render, screen, fireEvent } from '@testing-library/react';
import { Button } from '@/components/ui/button';

describe('Button', () => {
  test('renders with correct text', () => {
    render(<Button>Click me</Button>);
    expect(screen.getByRole('button', { name: /click me/i })).toBeInTheDocument();
  });

  test('calls onClick when clicked', () => {
    const handleClick = jest.fn();
    render(<Button onClick={handleClick}>Click me</Button>);
    fireEvent.click(screen.getByRole('button', { name: /click me/i }));
    expect(handleClick).toHaveBeenCalledOnce();
  });

  test('applies variant class', () => {
    render(<Button variant="destructive">Delete</Button>);
    const button = screen.getByRole('button', { name: /delete/i });
    expect(button).toHaveClass('bg-destructive');
  });
});
```

### Integration Tests
- **Framework**: React Testing Library or Cypress
- **Scope**: Component interactions, user flows, API mocking
- **Examples**:
  - Login flow
  - Repository creation and deletion
  - Trust score calculation and display
  - Webhook creation and testing

### End-to-End Tests
- **Framework**: Playwright
- **Scope**: Critical user journeys
- **Examples**:
  - New user sign up and onboarding
  - Developer pushing and pulling images
  - Security engineer investigating vulnerabilities
  - System administrator performing upgrade
  - Organization admin managing teams and permissions

### Visual Regression Tests
- **Framework**: Chromatic or Percy
- **Scope**: UI components and pages
- **Purpose**: Detect unintended visual changes

### Performance Tests
- **Framework**: Lighthouse CI
- **Metrics**: 
  - First Contentful Paint (FCP) < 1.5s
  - Largest Contentful Paint (LCP) < 2.5s
  - Cumulative Layout Shift (CLS) < 0.1
  - First Input Delay (FID) < 100ms
  - Time to Interactive (TTI) < 3.5s
- **Thresholds**: Set budgets for JavaScript, CSS, and image sizes

## Deployment and Build Process

### Build Commands
```bash
# Development
next dev

# Production build
next build

# Production start
next start

# Export static (if needed)
next export
```

### Environment Variables
- `NEXT_PUBLIC_API_URL`: Base URL for API requests
- `NEXT_PUBLIC_WS_URL`: WebSocket URL for real-time updates
- `NEXT_PUBLIC_ENABLE_FEATURE_FLAGS`: Comma-separated list of flags to enable
- `NEXT_PUBLIC_GA_ID`: Google Analytics ID (if used)
- `NEXT_PUBLIC_SENTRY_DSN`: Sentry DSN for error tracking
- `NEXT_PUBLIC_APP_VERSION`: Application version (from package.json)

### Dockerfile
```dockerfile
# Build stage
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# Production stage
FROM node:20-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/public ./public
EXPOSE 3000
CMD ["npm", "start"]
```

### Kubernetes Deployment
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kyros-frontend
  labels:
    app: kyros-frontend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: kyros-frontend
  template:
    metadata:
      labels:
        app: kyros-frontend
    spec:
      containers:
      - name: frontend
        image: kyros-frontend:latest
        ports:
        - containerPort: 3000
        env:
        - name: NEXT_PUBLIC_API_URL
          valueFrom:
            configMapKeyRef:
              name: kyros-config
              key: api-url
        - name: NEXT_PUBLIC_WS_URL
          valueFrom:
            configMapKeyRef:
              name: kyros-config
              key: ws-url
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /_next/data/v1/en/dashboard.json
            port: 3000
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /_next/data/v1/en/dashboard.json
            port: 3000
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: kyros-frontend-service
  labels:
    app: kyros-frontend
spec:
  selector:
    app: kyros-frontend
  ports:
    - protocol: TCP
      port: 80
      targetPort: 3000
  type: ClusterIP
```

## Conclusion

The Kyros frontend architecture is designed to be scalable, maintainable, and performant while providing an excellent user experience. By leveraging modern technologies like Next.js 15, React 18, TypeScript, and Tailwind CSS, and following best practices for state management, data fetching, and real-time updates, the frontend delivers a responsive and reliable interface for managing the software supply chain.

The architecture emphasizes:
- **Type Safety**: End-to-end TypeScript for fewer runtime errors
- **Modularity**: Reusable components and clear separation of concerns
- **Performance**: Efficient data fetching, caching, and code splitting
- **Real-time Capabilities**: WebSocket connections for live updates
- **Accessibility**: WCAG 2.1 AA compliance as a core requirement
- **Internationalization**: Ready for global audiences
- **Testing**: Comprehensive testing strategy to ensure quality
- **Deployability**: Docker and Kubernetes ready for cloud-native deployment

This foundation enables rapid feature development while maintaining a high-quality user experience that meets the needs of platform engineers, developers, security engineers, organization administrators, and system administrators.