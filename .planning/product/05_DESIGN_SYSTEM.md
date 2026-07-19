# Kyros Design System

## Overview
The Kyros Design System provides a cohesive visual language, component library, and interaction patterns for building consistent, accessible, and beautiful user interfaces across the platform. It ensures that all parts of Kyros feel like a unified product while allowing for flexibility and customization where needed.

## Design Principles
1. **Clarity**: Interface elements should be clear and unambiguous in their purpose and function.
2. **Consistency**: Similar elements should look and behave similarly across the application.
3. **Feedback**: Users should receive immediate and clear feedback for their actions.
4. **Efficiency**: Common tasks should be accessible with minimal steps.
5. **Accessibility**: The interface should be usable by people with diverse abilities.
6. **Scalability**: The system should accommodate growth in features and users.
7. **Trust**: Visual design should convey security, reliability, and professionalism.

## Foundation

### Design Tokens
Design tokens are the visual atoms of the design system — named entities that store visual design attributes. We use them instead of hard-coded values to maintain consistency and enable easy theming.

#### Color System
Kyros uses a semantic color system where colors are defined by their usage rather than their specific hue.

##### Primary Palette
- **Primary**: `#2563EB` (Blue-600) - Used for primary actions, links, and key brand elements
- **Primary Foreground**: `#FFFFFF` (White) - Text and icons on primary background
- **Primary Dark**: `#1D4ED8` (Blue-700) - Hover/active states for primary elements
- **Primary Light**: `#DBEAFE` (Blue-100) - Subtle backgrounds and accents

##### Neutral Palette
Used for text, backgrounds, borders, and subtle accents.
- **Gray-50**: `#F9FAFB`
- **Gray-100**: `#F3F4F6`
- **Gray-200**: `#E5E7EB`
- **Gray-300**: `#D1D5DB`
- **Gray-400**: `#9CA3AF`
- **Gray-500**: `#6B7280`
- **Gray-600**: `#4B5563`
- **Gray-700**: `#374151`
- **Gray-800**: `#1F2937`
- **Gray-900**: `#111827`

##### Semantic Colors
Colors that convey specific meaning or state.
- **Success**: `#10B981` (Emerald-500)
- **Success Dark**: `#059669` (Emerald-600)
- **Success Light**: `#D1FAE5` (Emerald-100)
- **Warning**: `#F59E0B` (Amber-500)
- **Warning Dark**: `#D97706` (Amber-600)
- **Warning Light**: `#FEF3C7` (Amber-100)
- **Error**: `#EF4444` (Red-500)
- **Error Dark**: `#DC2626` (Red-600)
- **Error Light**: `#FEE2E2` (Red-100)
- **Info**: `#3B82F6` (Blue-500)
- **Info Dark**: `#2563EB` (Blue-600)
- **Info Light**: `#DBEAFE` (Blue-100)

##### Trust Score Colors
Specialized color scale for trust score visualization.
- **Trusted** (0.9-1.0): `#10B981` (Emerald-500)
- **High** (0.7-0.89): `#3B82F6` (Blue-500)
- **Medium** (0.5-0.69): `#F59E0B` (Amber-500)
- **Low** (0.3-0.49): `#EF4444` (Red-500)
- **Untrusted** (0.0-0.29): `#7F1D1D` (Red-800)

#### Typography
Typography settings ensure readable, accessible, and harmonious text throughout the interface.

##### Font Family
- **Primary**: `Inter, system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, 'Noto Sans', sans-serif`
- **Monospace**: `'JetBrains Mono', 'Fira Code', 'Droid Sans Mono', 'Source Code Pro', Menlo, Monaco, Consolas, 'Courier New', monospace`

##### Font Weights
- **Light**: 300
- **Regular**: 400
- **Medium**: 500
- **Semi-bold**: 600
- **Bold**: 700
- **Extra-bold**: 800

##### Font Sizes (Base 16px)
- **2xs**: 0.625rem (10px)
- **xs**: 0.75rem (12px)
- **sm**: 0.875rem (14px)
- **base**: 1rem (16px)
- **lg**: 1.125rem (18px)
- **xl**: 1.25rem (20px)
- **2xl**: 1.5rem (24px)
- **3xl**: 1.875rem (30px)
- **4xl**: 2.25rem (36px)
- **5xl**: 3rem (48px)
- **6xl**: 3.75rem (60px)
- **7xl**: 4.5rem (72px)
- **8xl**: 6rem (96px)
- **9xl**: 8rem (128px)

##### Line Heights
- **Tight**: 1.25
- **Snug**: 1.375
- **Normal**: 1.5
- **Relaxed**: 1.625
- **Loose**: 2

##### Letter Spacing
- **Tighter**: -0.05em
- **Tight**: -0.025em
- **Normal**: 0
- **Wide**: 0.025em
- **Wider**: 0.05em
- **Widest**: 0.1em

#### Spacing System
Based on a 4px grid for consistent spacing and alignment.

##### Spatial Scale
- **0**: 0px
- **1**: 0.25rem (4px)
- **2**: 0.5rem (8px)
- **3**: 0.75rem (12px)
- **4**: 1rem (16px)
- **5**: 1.25rem (20px)
- **6**: 1.5rem (24px)
- **7**: 1.75rem (28px)
- **8**: 2rem (32px)
- **9**: 2.25rem (36px)
- **10**: 2.5rem (40px)
- **11**: 2.75rem (44px)
- **12**: 3rem (48px)
- **14**: 3.5rem (56px)
- **16**: 4rem (64px)
- **20**: 5rem (80px)
- **24**: 6rem (96px)
- **28**: 7rem (112px)
- **32**: 8rem (128px)
- **36**: 9rem (144px)
- **40**: 10rem (160px)
- **44**: 11rem (176px)
- **48**: 12rem (192px)
- **52**: 13rem (208px)
- **56**: 14rem (224px)
- **60**: 15rem (240px)
- **64**: 16rem (256px)
- **72**: 18rem (288px)
- **80**: 20rem (320px)
- **96**: 24rem (384px)

#### Border Radius
- **None**: 0px
- **Sm**: 0.25rem (4px)
- **Md**: 0.375rem (6px)
- **Lg**: 0.5rem (8px)
- **Xl**: 0.75rem (12px)
- **2xl**: 1rem (16px)
- **3xl**: 1.5rem (24px)
- **Full**: 9999px (pill)

#### Shadows
- **Sm**: `0 1px 2px 0 rgb(0 0 0 / 0.05)`
- **Md**: `0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -1px rgb(0 0 0 / 0.06)`
- **Lg**: `0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -2px rgb(0 0 0 / 0.05)`
- **Xl**: `0 20px 25px -5px rgb(0 0 0 / 0.1), 0 10px 10px -5px rgb(0 0 0 / 0.04)`
- **2xl**: `0 25px 50px -12px rgb(0 0 0 / 0.25)`
- **Inner**: `inset 0 2px 4px 0 rgb(0 0 0 / 0.05)`
- **None**: `none`

#### Opacity
- **0**: 0%
- **5**: 5%
- **10**: 10%
- **20**: 20%
- **25**: 25%
- **30**: 30%
- **40**: 40%
- **50**: 50%
- **60**: 60%
- **70**: 70%
- **75**: 75%
- **80**: 80%
- **90**: 90%
- **95**: 95%
- **100**: 100%

#### Transitions
- **None**: `none`
- **All**: `all 0.2s cubic-bezier(0.4, 0, 0.2, 1)`
- **Common**: `background-color 0.2s cubic-bezier(0.4, 0, 0.2, 1), border-color 0.2s cubic-bezier(0.4, 0, 0.2, 1), color 0.2s cubic-bezier(0.4, 0, 0.2, 1), opacity 0.2s cubic-bezier(0.4, 0, 0.2, 1), box-shadow 0.2s cubic-bezier(0.4, 0, 0.2, 1), transform 0.2s cubic-bezier(0.4, 0, 0.2, 1)`
- **Colors**: `background-color 0.2s cubic-bezier(0.4, 0, 0.2, 1), border-color 0.2s cubic-bezier(0.4, 0, 0.2, 1), color 0.2s cubic-bezier(0.4, 0, 0.2, 1), opacity 0.2s cubic-bezier(0.4, 0, 0.2, 1)`
- **Opacity**: `opacity 0.2s cubic-bezier(0.4, 0, 0.2, 1)`
- **Shadow**: `box-shadow 0.2s cubic-bezier(0.4, 0, 0.2, 1)`
- **Transform**: `transform 0.2s cubic-bezier(0.4, 0, 0.2, 1)`

### Grid System
Kyros uses a flexible 12-column grid system for layout consistency.

#### Container Widths
- **Sm**: 640px
- **Md**: 768px
- **Lg**: 1024px
- **Xl**: 1280px
- **2xl**: 1536px

#### Column Gutter
- **Default**: 1.5rem (24px)
- **None**: 0px

#### Breakpoints
- **Sm**: ≥640px
- **Md**: ≥768px
- **Lg**: ≥1024px
- **Xl**: ≥1280px
- **2xl**: ≥1536px

### Iconography
Kyros uses Heroicons as its primary icon set, with fallback to custom icons when needed.

#### Icon Styles
- **Outline**: Default style for inactive/neutral states
- **Solid**: Used for active/emphasized states
- **Mini**: Smaller variant for compact spaces (16x16px)
- **Standard**: Default size (20x20px)
- **Large**: Larger variant for emphasis (24x24px)

#### Icon Usage Guidelines
- Use icons to supplement text labels, not replace them (except in space-constrained contexts)
- Maintain consistent icon style within a context
- Ensure icons are recognizable and convey clear meaning
- Provide accessible labels for icon-only buttons
- Scale icons appropriately for different contexts

### Motion and Animation
Thoughtful use of motion enhances usability and provides feedback.

#### Duration Standards
- **Fast**: 80ms
- **Normal**: 150ms
- **Slow**: 250ms
- **Slower**: 350ms

#### Easing Functions
- **Default**: `cubic-bezier(0.4, 0, 0.2, 1)` (standard material ease-out)
- **In**: `cubic-bezier(0.4, 0, 1, 1)`
- **Out**: `cubic-bezier(0, 0, 0.2, 1)`
- **In-out**: `cubic-bezier(0.4, 0, 0.2, 1)`

#### Animation Principles
- **Purposeful**: Every animation should serve a clear purpose (feedback, orientation, etc.)
- **Natural**: Movements should follow physical principles (mass, velocity, friction)
- **Non-distracting**: Animations should not interfere with primary tasks
- **Performance-conscious**: Prefer transform and opacity animations for better performance
- **Respects preferences**: Honor `prefers-reduced-motion` media query

### Accessibility Standards
Kyros is committed to WCAG 2.1 AA compliance.

#### Color Contrast
- **Normal text**: Minimum 4.5:1 contrast ratio
- **Large text**: Minimum 3:1 contrast ratio (18pt+ or 14pt+ bold)
- **Graphical objects**: Minimum 3:1 contrast ratio
- **UI components**: Minimum 3:1 contrast ratio for interactive elements

#### Touch Targets
- Minimum size: 44x44px (with 8px spacing between touch targets)
- Preferred size: 48x48px

#### Keyboard Navigation
- All interactive elements must be accessible via keyboard
- Logical tab order that follows visual flow
- Visible focus indicators with minimum 3:1 contrast
- Skip navigation links for screen reader users

#### Screen Reader Support
- Proper semantic HTML (headings, lists, landmarks)
- Descriptive alt text for meaningful images
- ARIA labels and roles where native semantics are insufficient
- Live regions for dynamic content updates
- Proper form labeling and error messaging

#### Responsive Design
- Content reflows appropriately at different breakpoints
- Touch targets remain accessible on small screens
- Navigation adapts to screen size (collapsible sidebar, etc.)
- Text remains readable without zooming

## Component Guidelines

### Component Categories
Components are organized by their primary function in the interface.

#### 1. Layout Components
- **Container**: Centers content and constrains width
- **Grid**: 12-column responsive grid system
- **Flex**: Flexbox container for one-dimensional layouts
- **Stack**: Vertical or horizontal spacing utility
- **Sidebar**: Collapsible navigation panel
- **Header**: Page/application header
- **Footer**: Page/application footer

#### 2. Visualization Components
- **Avatar**: User/image representation
- **Badge**: Small status indicator
- **Banner**: Important notice or call-to-action
- **Divider**: Visual separator
- **Image**: Optimized image display with fallback
- **KV**: Key-value pair display
- **List**: Vertical list of items
- **Table**: Structured data display with sorting/filtering
- **Tag**: Categorical label
- **Tooltip**: Contextual information on hover/focus

#### 3. Feedback Components
- **Alert**: Important message requiring attention
- **Progress**: Visual indicator of completion
- **Spinner**: Loading indicator
- **Skeleton**: Placeholder UI during loading
- **Toast**: Temporary notification
- **Validation**: Form field validation state

#### 4. Input Components
- **Button**: Trigger for actions
- **Checkbox**: Binary selection
- **Combobox**: Text input with dropdown suggestions
- **Form**: Container for related inputs
- **Input**: Text entry field
- **Label**: Text label for form elements
- **Radio Group**: Mutually exclusive selection
- **Select**: Dropdown selection
- **Switch**: Toggle for on/off state
- **Textarea**: Multi-line text input

#### 5. Navigation Components
- **Breadcrumb**: Hierarchical navigation trail
- **Link**: Navigation to another resource
- **Menu**: List of actions or options
- **Pagination**: Navigation through paginated content
- **Tabs**: Organized content sections
- **Stepper**: Multi-step process indicator

#### 6. Overlay Components
- **Dialog**: Modal window requiring user action
- **Drawer**: Slide-in panel from screen edge
- **Popover**: Temporary content anchored to an element
- **Tooltip**: Brief contextual information

#### 7. Data Visualization Components
- **Chart**: Generic chart container
- **Graph**: Node-link visualization
- **Map**: Geographic data display
- **Metric**: Single important value with context
- **Sparkline**: Simple line chart for trends
- **Gauge**: Circular indicator for values within range

### Component Composition Principles
1. **Atomic Design**: Build from atoms → molecules → organisms → templates → pages
2. **Single Responsibility**: Each component should do one thing well
3. **Composability**: Simple components combine to create complex ones
4. **Configurability**: Use props for variation rather than creating new components
5. **Extensibility**: Design for future enhancement without breaking changes
6. **Performance**: Minimize re-renders and expensive operations
7. **Accessibility**: Build with accessibility in mind from the start
8. **Documentation**: Clear usage guidelines and examples

## Theming Capabilities

### Light/Dark Mode
Kyros supports both light and dark themes with automatic switching based on system preference.

#### Theme Variables
All colors are defined as CSS variables that switch based on the `[data-theme]` attribute.

```css
:root {
  --color-background: var(--gray-50);
  --color-foreground: var(--gray-900);
  --color-primary: var(--blue-600);
  --color-primary-foreground: var(--white);
  /* ... */
}

[data-theme="dark"] {
  --color-background: var(--gray-900);
  --color-foreground: var(--gray-50);
  --color-primary: var(--blue-400);
  --color-primary-foreground: var(--gray-900);
  /* ... */
}
```

#### Implementation Approach
- CSS variables for all themeable values
- Class-based toggling (`data-theme="light"` vs `data-theme="dark"`)
- Automatic detection via `prefers-color-scheme` media query
- Manual override capability for users
- Transition support for smooth theme changes

### Branding and White Labeling
Organizations can customize Kyros to match their brand identity.

#### Customizable Elements
- **Logo**: Replace default logo with organization's logo
- **Primary Color**: Adjust primary color to match brand
- **Favicon**: Custom browser tab icon
- **Footer Text**: Customizable footer content
- **Login Page**: Custom background, logo, and text
- **Email Templates**: Branded notification emails
- **Font Family**: Option to use custom fonts (with performance considerations)

#### Limitations
- Secondary colors are derived from primary for consistency
- Structural layout remains consistent for usability
- Interactive patterns preserved for familiarity
- Accessibility standards must be maintained

## Implementation Guidelines

### Development Practices
1. **Component First**: Develop components in isolation using Storybook
2. **Type Safety**: Use TypeScript for all frontend code
3. **Accessibility First**: Test with screen readers and keyboard navigation
4. **Performance Conscious**: Monitor bundle size and render performance
5. **Testing**: Write unit tests for components and integration tests for interactions
6. **Documentation**: Maintain Storybook stories and usage documentation
7. **Versioning**: Follow semantic versioning for breaking changes

### CSS Architecture
- **Utility-First**: Leverage Tailwind CSS for rapid UI development
- **Component Styles**: Use `@apply` for extracting utility classes into CSS classes
- **BEM Naming**: For custom CSS, use Block-Element-Modifier naming
- **CSS Variables**: For themeable values and design tokens
- **PostCSS**: For autoprefixing and modern CSS features
- **PurgeCSS**: Remove unused CSS in production builds

### JavaScript/TypeScript Practices
- **Functional Components**: Prefer React functional components with hooks
- **Custom Hooks**: Extract reusable logic into custom hooks
- **Context API**: For global state that doesn't require external libraries
- **State Management**: Use React Query for server state, Zustand/Jotai for client state
- **Error Boundaries**: Catch and gracefully handle JavaScript errors
- **Suspense**: For code splitting and lazy loading
- **Testing Library**: For user-centric component testing

### Asset Management
- **Icons**: Heroicons as primary set, SVGs for custom icons
- **Images**: Optimized formats (WebP, AVIF) with responsive sizes
- **Fonts**: Self-hosted with font-display: swap
- **Animations**: Framer Motion for complex interactions, CSS for simple transitions
- **Internationalization**: next-i18next for translation support

## Usage Examples

### Button Variants
```tsx
import { Button } from '@/components/ui/button'

// Default button
<Button>Default</Button>

// Destructive action
<Button variant="destructive">Delete</Button>

// Outline style
<Button variant="outline">Secondary</Button>

// Ghost style (transparent background)
<Button variant="ghost">Text</Button>

// Link style
<Button variant="link">Link</Button>

// Sizes
<Button size="icon">
  <Search />
</Button>
<Button size="sm">Small</Button>
<Button size="lg">Large</Button>
```

### Form Elements
```tsx
import { 
  Form, 
  Field, 
  Label, 
  Input, 
  Textarea, 
  Select, 
  Checkbox, 
  RadioGroup, 
  Switch,
  Button 
} from '@/components/ui/form'

<Form>
  <Field>
    <Label>Email address</Label>
    <Input type="email" placeholder="you@example.com" />
    <p className="mt-2 text-sm text-gray-500">
      We'll never share your email with anyone else.
    </p>
  </Field>
  
  <Field>
    <Label>Password</Label>
    <Input type="password" placeholder="••••••••" />
  </Field>
  
  <Field>
    <Label>Newsletter preferences</Label>
    <div className="space-y-2">
      <Checkbox>
        <Label>Product updates</Label>
      </Checkbox>
      <Checkbox>
        <Label>Promotional offers</Label>
      </Checkbox>
    </div>
  </Field>
  
  <Field>
    <Label>Account type</Label>
    <RadioGroup>
      <RadioGroup.Label>Personal</RadioGroup.Label>
      <RadioGroup.Value value="personal" />
      <RadioGroup.Label>Business</RadioGroup.Label>
      <RadioGroup.Value value="business" />
    </RadioGroup>
  </Field>
  
  <Field>
    <Label>Notifications</Label>
    <Switch checked={true} />
    <span className="ml-2">Enable notifications</span>
  </Field>
  
  <Field>
    <Button type="submit">Save changes</Button>
  </Field>
</Form>
```

### Card Layout
```tsx
import { 
  Card, 
  CardHeader, 
  CardTitle, 
  CardDescription, 
  CardContent, 
  CardFooter,
  Button
} from '@/components/ui/card'

<Card className="w-full">
  <CardHeader className="pb-4">
    <CardTitle className="text-lg font-semibold">
      Repository Statistics
    </CardTitle>
    <CardDescription className="text-sm text-gray-500">
      Overview of your repository activity
    </CardDescription>
  </CardHeader>
  <CardContent className="space-y-4">
    <div className="grid grid-cols-2 gap-4">
      <div className="bg-gray-50 p-4 rounded-lg">
        <div className="text-sm font-medium text-gray-500">
          Total Images
        </div>
        <div className="text-2xl font-bold text-gray-900">
          1,248
        </div>
        <div className="text-sm text-green-600">
          +12% this month
        </div>
      </div>
      
      <div className="bg-gray-50 p-4 rounded-lg">
        <div className="text-sm font-medium text-gray-500">
          Average Trust Score
        </div>
        <div className="text-2xl font-bold text-gray-900">
          0.78
        </div>
        <div className="text-sm text-green-600">
          +0.05 this month
        </div>
      </div>
    </div>
  </CardContent>
  <CardFooter className="pt-4">
    <Button variant="outline">
      View Detailed Report
    </Button>
  </CardFooter>
</Card>
```

### Navigation Patterns
```tsx
import { 
  Sidebar, 
  SidebarHeader, 
  SidebarFooter, 
  SidebarMenu, 
  MenuItem,
  Breadcrumb,
  BreadcrumbItem,
  Tabs,
  TabsList,
  TabsTrigger,
  TabsContent
} from '@/components/ui/navigation'

<div className="flex h-screen bg-gray-50">
  {/* Sidebar Navigation */}
  <Sidebar className="w-64 border-r">
    <SidebarHeader className="px-4 py-3">
      <h2 className="text-lg font-semibold">Kyros</h2>
    </SidebarHeader>
    
    <SidebarMenu>
      <MenuItem href="/" icon={<Home />}>
        Dashboard
      </MenuItem>
      <MenuItem href="/repositories" icon={<Folder />}>
        Repositories
      </MenuItem>
      <MenuItem href="/images" icon={<Image />}>
        Images
      </MenuItem>
      <MenuItem href="/trust-score" icon={<ShieldCheck />}>
        Trust Score
      </MenuItem>
      <MenuItem href="/security" icon={<ShieldAlert />}>
        Security
      </MenuItem>
      <MenuItem href="/users" icon={<Users />}>
        Users & Teams
      </MenuItem>
      <MenuItem href="/settings" icon={<Settings />}>
        Settings
      </MenuItem>
    </SidebarMenu>
    
    <SidebarFooter className="px-4 py-3">
      <Button variant="outline" size="icon" href="/help">
        <QuestionCircle size={20} />
      </Button>
    </SidebarFooter>
  </Sidebar>
  
  {/* Main Content */}
  <main className="flex-1 p-6">
    <Breadcrumb className="mb-4">
      <BreadcrumbItem href="/">Home</BreadcrumbItem>
      <BreadcrumbItem href="/repositories">Repositories</BreadcrumbItem>
      <BreadcrumbItem>web-app</BreadcrumbItem>
    </Breadcrumb>
    
    <Tabs defaultValue="overview" className="w-full">
      <TabsList className="grid w-full grid-cols-4">
        <TabsTrigger value="overview">Overview</TabsTrigger>
        <TabsTrigger value="tags">Tags</TabsTrigger>
        <TabsTrigger value="security">Security</TabsTrigger>
        <TabsTrigger value="analytics">Analytics</TabsTrigger>
      </TabsList>
      
      <TabsContent value="overview">
        {/* Overview content */}
      </TabsContent>
      
      <TabsContent value="tags">
        {/* Tags content */}
      </TabsContent>
      
      <TabsContent value="security">
        {/* Security content */}
      </TabsContent>
      
      <TabsContent value="analytics">
        {/* Analytics content */}
      </TabsContent>
    </Tabs>
  </main>
</div>
```

## Extensibility and Customization

### Theme Extension
Organizations can extend the design system with their own tokens:

```css
/* In your custom CSS file */
:root {
  /* Extend or override existing tokens */
  --color-brand: #ff6b6b; /* Custom brand color */
  --color-brand-foreground: #ffffff;
  
  /* Add new semantic tokens */
  --color-status-new: #ff6b6b;
  --color-status-processing: #4cc9f0;
  --color-status-complete: #4caf50;
}

/* Use in components */
.button-custom {
  background-color: var(--color-brand);
  color: var(--color-brand-foreground);
}
```

### Component Extension
Create new components that compose existing ones:

```tsx
// CustomAvatarWithStatus.tsx
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { useUserStatus } from '@/hooks/useUserStatus'

export function CustomAvatarWithStatus({ userId, size = 'md' }: { 
  userId: string; 
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl' 
}) {
  const { status, isOnline } = useUserStatus(userId)
  
  return (
    <div className="relative inline-block">
      <Avatar size={size}>
        <AvatarImage src={`/avatars/${userId}.jpg`} alt="User avatar" />
        <AvatarFallback>{getInitials(userId)}</AvatarFallback>
      </Avatar>
      {isOnline && (
        <Badge 
          variant={status === 'active' ? 'success' : 'secondary'} 
          size="xs" 
          className="absolute -bottom-1 -right-1"
        >
          {status}
        </Badge>
      )}
    </div>
  )
}
```

### Pattern Library
Document common composition patterns:

#### Action Bar Pattern
```tsx
// ActionBar.tsx
import { Button } from '@/components/ui/button'

interface ActionBarProps {
  primaryAction: React.ReactNode
  secondaryActions?: React.ReactNode[]
  destructiveActions?: React.ReactNode[]
}

export function ActionBar({ 
  primaryAction, 
  secondaryActions = [], 
  destructiveActions = [] 
}: ActionBarProps) {
  return (
    <div className="flex flex-wrap items-center gap-3">
      <div className="flex-1">{primaryAction}</div>
      <div className="flex flex-wrap gap-2">
        {secondaryActions.map((action, index) => (
          <Button key={index} variant="outline">{action}</Button>
        ))}
        {destructiveActions.map((action, index) => (
          <Button 
            key={index} 
            variant="destructive" 
            size="sm"
          >
            {action}
          </Button>
        ))}
      </div>
    </div>
  )
}
```

#### Data Table Pattern
```tsx
// DataTable.tsx
import { 
  Table, 
  TableHeader, 
  TableBody, 
  TableRow, 
  TableCell, 
  TableHead, 
  TableCaption 
} from '@/components/ui/table'
import { Checkbox } from '@/components/ui/checkbox'
import { Button } from '@/components/ui/button'
import { Menu, MenuTrigger, MenuContent, MenuItem } from '@/components/ui/menu'

interface DataTableProps<T> {
  columns: Array<{
    key: keyof T
    label: string
    sortable?: boolean
    format?: (value: T[keyof T]) => string
  }>
  data: T[]
  onRowAction?: (item: T, action: string) => void
  enableSelection?: boolean
}

export function DataTable<T>({ 
  columns, 
  data, 
  onRowAction, 
  enableSelection = false 
}: DataTableProps<T>) {
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set())
  
  const handleSelectAll = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.checked) {
      setSelectedIds(new Set(data.map(item => String((item as any).id))))
    } else {
      setSelectedIds.clear()
    }
  }
  
  const handleToggleSelect = (id: string, e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.checked) {
      selectedIds.add(id)
    } else {
      selectedIds.delete(id)
    }
  }
  
  return (
    <div className="space-y-4">
      {enableSelection && (
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Checkbox 
              checked={selectedIds.size === data.length}
              onChange={handleSelectAll}
              className="h-4 w-4"
            />
            <span className="text-sm font-medium">
              {selectedIds.size} of {data.length} selected
            </span>
          </div>
          <Button variant="outline" size="sm" disabled={selectedIds.size === 0}>
            Actions
          </Button>
        </div>
      )}
      
      <Table>
        <TableCaption>
          {data.length} items
        </TableCaption>
        <Thead>
          <Tr>
            {enableSelection && (
              <Th className="w-4">
                <Checkbox 
                  checked={selectedIds.size === data.length}
                  onChange={handleSelectAll}
                  className="h-4 w-4"
                />
              </Th>
            )}
            {columns.map(column => (
              <Th key={column.key} className="text-left">
                {column.label}
                {column.sortable && (
                  <SortIcon className="ml-1 h-4 w-4 text-gray-400" />
                )}
              </Th>
            ))}
            {enableSelection && (
              <Th className="text-right">Actions</Th>
            )}
          </Tr>
        </Thead>
        <Tbody>
          {data.map((item, index) => {
            const itemId = String((item as any).id)
            const isSelected = selectedIds.has(itemId)
            
            return (
              <Tr 
                key={index} 
                className={isSelected ? 'bg-muted/50' : ''}
              >
                {enableSelection && (
                  <Td className="w-4">
                    <Checkbox 
                      checked={isSelected}
                      onChange={(e) => handleToggleSelect(itemId, e)}
                      className="h-4 w-4"
                    />
                  </Td>
                )}
                {columns.map((column, colIndex) => {
                  const value = (item as any)[column.key]
                  const formatted = column.format 
                    ? column.format(value) 
                    : String(value)
                  
                  return (
                    <Td 
                      key={colIndex} 
                      className={colIndex === columns.length - 1 
                        ? 'text-right' 
                        : 'text-left'}
                    >
                      {formatted}
                    </Td>
                  )
                })}
                </Tr>
              ))}
            </Tbody>
          </Table>
          
          {enableSelection && selectedIds.size > 0 && (
            <div className="flex items-center justify-end space-x-3">
              <Button 
                variant="destructive" 
                onClick={() => {
                  // Handle bulk delete
                  selectedIds.clear()
                }}
              >
                Delete {selectedIds.size} Selected
              </Button>
              <Button 
                variant="default" 
                onClick={() => {
                  // Handle bulk action
                  selectedIds.clear()
                }}
              >
                Perform Action on {selectedIds.size} Selected
              </Button>
            </div>
          )}
        </div>
      )
    }
  )
}
```

## Maintenance and Evolution

### Contribution Guidelines
1. **Follow Existing Patterns**: Match the coding style and conventions of existing components
2. **Include Tests**: Write unit tests for new components
3. **Update Documentation**: Add component to Storybook and document usage
4. **Consider Accessibility**: Ensure new components meet WCAG 2.1 AA standards
5. **Performance Review**: Consider performance implications of new components
6. **Design Review**: Get feedback from design team on visual and interaction design
7. **Changelog**: Document changes in the release notes

### Versioning Strategy
- **Semantic Versioning**: MAJOR.MINOR.PATCH
- **Major**: Breaking changes requiring migration
- **Minor**: Backward-compatible new features
- **Patch**: Backward-compatible bug fixes
- **Deprecation**: Clear warnings with migration paths before removal

### Design System Updates
- **Regular Reviews**: Quarterly reviews of design system usage and effectiveness
- **Feedback Channels**: Dedicated channels for designer and developer feedback
- **Adoption Metrics**: Track usage of design system components vs. custom implementations
- **Evolution Process**: Clear proposal, review, and approval process for changes

## Conclusion

The Kyros Design System provides a solid foundation for building consistent, accessible, and beautiful user interfaces. By establishing clear principles for color, typography, spacing, components, and patterns, it enables teams to work efficiently while maintaining a cohesive user experience.

The system is designed to be extensible, allowing for customization and evolution while preserving the core principles that make Kyros feel like a unified product. Through thoughtful implementation and ongoing maintenance, the design system will continue to serve as a valuable asset for the Kyros platform.