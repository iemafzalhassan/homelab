# Kyros Component Library

## Overview
This document documents all reusable components in the Kyros platform, including their purpose, props, usage examples, and implementation guidelines. These components form the building blocks of the user interface, ensuring consistency and maintainability.

## Component Categories

### 1. Layout Components
Components that structure the page layout and positioning.

#### Container
Centers content and constrains width to a maximum size.

**Props**:
- `children`: React.ReactNode - Content to be centered
- `className`: string - Additional CSS classes
- `as`: keyof JSX.IntrinsicElements - HTML element to render (default: "div")

**Usage**:
```tsx
import { Container } from '@/components/ui/container'

<Container>
  <h1>Welcome to Kyros</h1>
  <p>This content is centered and constrained.</p>
</Container>
```

**Variants**:
- `sm`: max-width: 640px
- `md`: max-width: 768px
- `lg`: max-width: 1024px
- `xl`: max-width: 1280px
- `2xl`: max-width: 1536px
- `fluid`: no max-width (100% width)

#### Grid
12-column responsive grid system.

**Props**:
- `children`: React.ReactNode - Grid items
- `className`: string - Additional CSS classes
- `cols`: number - Number of columns (default: 12)
- `gap`: number or string - Gap between columns (default: 24px)
- `className`: string - Additional CSS classes

**Usage**:
```tsx
import { Grid } from '@/components/ui/grid'

<Grid cols={3} gap={4}>
  <div>Item 1</div>
  <div>Item 2</div>
  <div>Item 3</div>
</Grid>
```

**Responsive Columns**:
```tsx
<Grid 
  cols={{ 
    base: 1,      // 1 column on mobile
    sm: 2,        // 2 columns on small screens
    md: 3,        // 3 columns on medium screens
    lg: 4         // 4 columns on large screens
  }} 
  gap={4}
>
  {/* Grid items */}
</Grid>
```

#### Flex
Flexbox container for one-dimensional layouts.

**Props**:
- `children`: React.ReactNode - Flex items
- `className`: string - Additional CSS classes
- `direction`: 'row' | 'row-reverse' | 'column' | 'column-reverse' (default: 'row')
- `justify`: 'start' | 'end' | 'center' | 'between' | 'around' | 'evenly'evenly')
- `align`: 'start' | 'end' | 'center' | 'baseline' | 'stretch' (default: 'stretch')
- `wrap`: boolean | 'wrap' | 'wrap-reverse' (default: false)
- `gap`: number or string - Gap between items (default: 0)

**Usage**:
```tsx
import { Flex } from '@/components/ui/flex'

<Flex direction="row" justify="between" align="center" gap={4}>
  <div>Left</div>
  <div>Center</div>
  <div>Right</div>
</Flex>
```

#### Stack
Vertical or horizontal spacing utility.

**Props**:
- `children`: React.ReactNode - Stacked items
- `className`: string - Additional CSS classes
- `direction`: 'row' | 'column' (default: 'column')
- `spacing`: number or string - Space between items (default: 4px)
- `divider`: boolean - Show dividers between items (default: false)
- `dividerClass`: string - CSS class for dividers

**Usage**:
```tsx
import { Stack } from '@/components/ui/stack'

<Stack spacing={4}>
  <div>Item 1</div>
  <div>Item 2</div>
  <div>Item 3</div>
</Stack>
```

#### Sidebar
Collapsible navigation panel.

**Props**:
- `children`: React.ReactNode - Sidebar content
- `className`: string - Additional CSS classes
- `collapsible`: boolean - Whether sidebar can be collapsed (default: true)
- `collapsed`: boolean - Initial collapsed state (default: false)
- `defaultCollapsed`: boolean - Default collapsed state (default: false)
- `onToggle`: (collapsed: boolean) => void - Callback when collapsed state changes
- `breakpoint`: 'sm' | 'md' | 'lg' | 'xl' | '2xl' - Breakpoint for auto-collapse (default: 'lg')
- `width`: number or string - Sidebar width (default: 240px)
- `collapsedWidth`: number or string - Collapsed sidebar width (default: 64px)

**Usage**:
```tsx
import { Sidebar } from '@/components/ui/sidebar'
import { SidebarHeader } from '@/components/ui/sidebar'
import { SidebarMenu } from '@/components/ui/sidebar'
import { MenuItem } from '@/components/ui/menu'

<Sidebar defaultCollapsed={false} breakpoint="md">
  <SidebarHeader>
    <h2 className="text-lg font-semibold">Kyros</h2>
  </SidebarHeader>
  
  <SidebarMenu>
    <MenuItem href="/" icon={<Home />}>
      Dashboard
    </MenuItem>
    <MenuItem href="/repositories" icon={<Folder />}>
      Repositories
    </MenuItem>
    {/* ... */}
  </SidebarMenu>
</Sidebar>
```

#### Header
Page/application header.

**Props**:
- `children`: React.ReactNode - Header content
- `className`: string - Additional CSS classes
- `fixed`: boolean - Whether header is fixed to top (default: false)
- `height`: number or string - Header height (default: 64px)
- `shadow`: boolean - Whether to show bottom shadow (default: true)
- `background`: boolean - Whether to show background (default: true)

**Usage**:
```tsx
import { Header } from '@/components/ui/header'
import { HeaderLeft } from '@/components/ui/header'
import { HeaderCenter } from '@/components/ui/header'
import { HeaderRight } from '@/components/ui/header'

<Header fixed shadow background>
  <HeaderLeft>
    <a href="/" className="flex items-center space-x-2">
      <Logo className="h-8 w-8" />
      <span className="text-xl font-semibold">Kyros</span>
    </a>
  </HeaderLeft>
  
  <HeaderCenter>
    <SearchBar placeholder="Search repositories..." />
  </HeaderCenter>
  
  <HeaderRight>
    <UserAvatar />
    <Button variant="outline">New Repository</Button>
  </HeaderRight>
</Header>
```

#### Footer
Page/application footer.

**Props**:
- `children`: React.ReactNode - Footer content
- `className`: string - Additional CSS classes
- `height`: number or string - Footer height (default: 48px)
- `borderTop`: boolean - Whether to show top border (default: true)
- `background`: boolean - Whether to show background (default: true)

**Usage**:
```tsx
import { Footer } from '@/components/ui/footer'

<Footer borderTop background>
  <div className="flex justify-between items-center px-6">
    <span className="text-sm text-gray-500">
      Kyros v1.0.0 • © 2026 Kyros Project
    </span>
    <div className="flex space-x-4">
      <a href="/docs" className="text-sm text-gray-500 hover:text-gray-600">
        Documentation
      </a>
      <a href="/support" className="text-sm text-gray-500 hover:text-gray-600">
        Support
      </a>
    </div>
  </div>
</Footer>
```

### 2. Visualization Components
Components that display information visually.

#### Avatar
User or image representation.

**Props**:
- `src`: string - Image URL (optional)
- `alt`: string - Alternative text (optional)
- `size`: 'xs' | 'sm' | 'md' | 'lg' | 'xl' | '2xl' (default: 'md')
- `shape`: 'circle' | 'square' (default: 'circle')
- `fallback`: React.ReactNode - Content to show when image fails to load (optional)
- `onError`: () => void - Callback when image fails to load
- `className`: string - Additional CSS classes

**Usage**:
```tsx
import { Avatar } from '@/components/ui/avatar'

// With image
<Avatar 
  src="https://example.com/avatar.jpg" 
  alt="User Avatar" 
  size="lg" 
/>

// With initials fallback
<Avatar 
  size="md" 
  fallback="JD"
/>

// With custom fallback
<Avatar 
  size="sm" 
  fallback={<span className="text-xs font-medium">U</span>}
/>
```

#### Badge
Small status indicator.

**Props**:
- `children`: React.ReactNode - Badge content
- `className`: string - Additional CSS classes
- `variant`: 'default' | 'primary' | 'secondary' | 'destructive' | 'warning' | 'success' (default: 'secondary')
- `size`: 'xs' | 'sm' | 'md' | 'lg' | 'xl' (default: 'md')
- `pill`: boolean - Whether to use pill shape (default: false)

**Usage**:
```tsx
import { Badge } from '@/components/ui/badge'

<Badge variant="success">Success</Badge>
<Badge variant="destructive">Error</Badge>
<Badge variant="warning" size="sm">Warning</Badge>
<Badge variant="secondary" pill>New</Badge>
```

#### Banner
Important notice or call-to-action.

**Props**:
- `children`: React.ReactNode - Banner content
- `className`: string - Additional CSS classes
- `variant`: 'default' | 'primary' | 'secondary' | 'destructive' | 'warning' | 'info' (default: 'default')
- `title`: string - Banner title (optional)
- `description`: string - Banner description (optional)
- `action`: React.ReactNode - Action button/link (optional)
- `dismissible`: boolean - Whether banner can be dismissed (default: false)
- `onDismiss`: () => void - Callback when banner is dismissed
- `icon`: React.ReactNode - Icon to show (optional)

**Usage**:
```tsx
import { Banner } from '@/components/ui/banner'

<Banner 
  variant="warning" 
  title="Maintenance Scheduled"
  description="System maintenance will occur tonight from 2:00 AM to 4:00 AM UTC."
  action={<Button variant="outline" size="sm">Learn More</Button>}
  dismissible
/>
```

#### Divider
Visual separator.

**Props**:
- `className`: string - Additional CSS classes
- `orientation`: 'horizontal' | 'vertical' (default: 'horizontal')
- `variant`: 'solid' | 'dashed' | 'dotted' (default: 'solid')
- `color`: string - CSS color value (optional, uses currentColor by default)

**Usage**:
```tsx
import { Divider } from '@/components/ui/divider'

<Divider orientation="horizontal" className="my-4" />
<Divider orientation="vertical" className="mx-4" />
<Divider variant="dashed" className="my-2" />
```

#### Image
Optimized image display with fallback.

**Props**:
- `src`: string - Image URL
- `alt`: string - Alternative text
- `width`: number | string - Image width
- `height`: number | string - Image height
- `objectFit`: 'contain' | 'cover' | 'fill' | 'none' | 'scale-down' (default: 'cover')
- `objectPosition`: string - CSS object-position value (default: 'center')
- `loading`: 'eager' | 'lazy' (default: 'lazy')
- `placeholder`: 'blur' | 'empty' (default: 'empty')
- `blurDataURL`: string - Base64 encoded blurred placeholder (optional)
- `className`: string - Additional CSS classes
- `onError`: () => void - Callback when image fails to load
- `onLoad`: () => void - Callback when image loads

**Usage**:
```tsx
import { Image } from '@/components/ui/image'

<Image 
  src="https://example.com/image.jpg" 
  alt="Description" 
  width={400} 
  height={300} 
  objectFit="cover"
  className="rounded-lg"
/>
```

#### KV
Key-value pair display.

**Props**:
- `key`: string - Label/key
- `value`: React.ReactNode - Value/content
- `className`: string - Additional CSS classes
- `keyClassName`: string - CSS classes for key
- `valueClassName`: string - CSS classes for value
- `truncate`: boolean - Whether to truncate long values (default: false)
- `maxLength`: number - Maximum length before truncation (default: 100)

**Usage**:
```tsx
import { KV } from '@/components/ui/kv'

<KV key="Name" value="John Doe" />
<KV key="Email" value="john.doe@example.com" />
<KV key="Bio" value="Software engineer passionate about open source..." truncate maxLength={50} />
```

#### List
Vertical list of items.

**Props**:
- `children`: React.ReactNode - List items
- `className`: string - Additional CSS classes
- `spaced`: boolean - Whether to add space between items (default: false)
- `divider`: boolean - Whether to show dividers between items (default: false)
- `dividerClass`: string - CSS class for dividers
- `as`: keyof JSX.IntrinsicElements - HTML element to render (default: "ul")
- `role`: string - ARIA role (optional)

**Usage**:
```tsx
import { List } from '@/components/ui/list'
import { ListItem } from '@/components/ui/list'

<List spaced divider>
  <ListItem>Item 1</ListItem>
  <ListItem>Item 2</ListItem>
  <ListItem>Item 3</ListItem>
</List>
```

#### Table
Structured data display with sorting/filtering capabilities.

**Props**:
- `columns`: Array<{
  accessor: keyof T | string
  header: string
  width?: number | string
  align?: 'left' | 'center' | 'right'
  format?: (value: any) => string
  sortable?: boolean
  filterable?: boolean
  className?: string
  headerClassName?: string
  cellClassName?: string
}> - Column definitions
- `data`: T[] - Row data
- `className`: string - Additional CSS classes
- `striped`: boolean - Whether to show striped rows (default: false)
- `highlight`: boolean - Whether to highlight row on hover (default: true)
- `bordered`: boolean - Whether to show borders (default: true)
- `compact`: boolean - Whether to use compact padding (default: false)
- `sortBy`: keyof T | string - Initial sort column (optional)
- `sortDesc`: boolean - Whether initial sort is descending (default: false)
- `onSort`: (column: keyof T | string, desc: boolean) => void - Sort callback
- `filterValue`: string - Current filter value (optional)
- `onFilter`: (value: string) => void - Filter callback
- `showFilter`: boolean - Whether to show filter input (default: true)
- `placeholder`: string - Filter input placeholder (default: "Filter...")
- `pageSize`: number - Rows per page (default: 10)
- `currentPage`: number - Current page (default: 1)
- `onPageChange`: (page: number) => void - Page change callback
- `showPagination`: boolean - Whether to show pagination controls (default: true)
- `selectable`: boolean - Whether rows are selectable (default: false)
- `selected`: Set<string | number> - Selected row IDs (default: empty set)
- `onSelectionChange`: (selected: Set<string | number>) => void - Selection change callback
- `showSelection`: boolean - Whether to show selection checkboxes (default: false)
- `emptyState`: React.ReactNode - Content to show when no data (optional)
- `loading`: boolean - Whether to show loading state (default: false)
- `loadingSpinner`: React.ReactNode - Custom loading spinner (optional)
- `error`: string | null - Error message (optional)
- `className`: string - Additional CSS classes

**Usage**:
```tsx
import { Table } from '@/components/ui/table'

const columns = [
  { accessor: 'name', header: 'Name', sortable: true, filterable: true },
  { accessor: 'email', header: 'Email', sortable: true },
  { accessor: 'role', header: 'Role', sortable: false },
  { accessor: 'lastLogin', header: 'Last Access', sortable: true, format: (date) => formatDate(date) }
]

const users = [
  { id: 1, name: 'John Doe', email: 'john@example.com', role: 'admin', lastLogin: new Date() },
  { id: 2, name: 'Jane Smith', email: 'jane@example.com', role: 'user', lastLogin: new Date(Date.now() - 86400000) }
]

<Table 
  columns={columns} 
  data={users} 
  striped 
  highlight 
  bordered
  selectable
  showSelection
/>
```

#### Tag
Categorical label.

**Props**:
- `children`: React.ReactNode - Tag content
- `className`: string - Additional CSS classes
- `variant`: 'default' | 'primary' | 'secondary' | 'destructive' | 'warning' | 'info' (default: 'secondary')
- `size`: 'xs' | 'sm' | 'md' | 'lg' | 'xl' (default: 'md')
- `pill`: boolean - Whether to use pill shape (default: true)

**Usage**:
```tsx
import { Tag } from '@/components/ui/tag'

<Tag variant="primary">Production</Tag>
<Tag variant="success" size="sm">Verified</Tag>
<Tag variant="warning">Pending</Tag>
```

#### Tooltip
Contextual information on hover/focus.

**Props**:
- `children`: React.ReactNode - Content that triggers the tooltip
- `content`: React.ReactNode - Tooltip content
- `className`: string - Additional CSS classes
- `contentClassName`: string - CSS classes for tooltip content
- `placement`: 'top' | 'top-start' | 'top-end' | 'bottom' | 'bottom-start' | 'bottom-end' | 'left' | 'left-start' | 'left-end' | 'right' | 'right-start' | 'right-end' (default: 'top')
- `offset`: number - Distance from trigger (default: 4)
- `delay`: number - Delay in ms before showing (default: 0)
- `hideDelay`: number - Delay in ms before hiding (default: 0)
- `arrow`: boolean - Whether to show arrow (default: true)
- `className`: string - Additional CSS classes

**Usage**:
```tsx
import { Tooltip } from '@/components/ui/tooltip'

<Tooltip content="View detailed repository information">
  <Button variant="outline" size="icon">
    <Info size={20} />
  </Button>
</Tooltip>
```

### 3. Feedback Components
Components that provide feedback to users.

#### Alert
Important message requiring attention.

**Props**:
- `children`: React.ReactNode - Alert content
- `className`: string - Additional CSS classes
- `variant`: 'default' | 'primary' | 'secondary' | 'destructive' | 'warning' | 'success' (default: 'default')
- `title`: string - Alert title (optional)
- `description`: string - Alert description (optional)
- `icon`: React.ReactNode - Icon to show (optional)
- `closable`: boolean - Whether alert can be closed (default: false)
- `onClose`: () => void - Callback when alert is closed
- `className`: string - Additional CSS classes

**Usage**:
```tsx
import { Alert } from '@/components/ui/alert'

<Alert 
  variant="warning" 
  title="Maintenance Scheduled"
  description="System maintenance will occur tonight from 2:00 AM to 4:00 AM UTC."
  closable
/>
```

#### Progress
Visual indicator of completion.

**Props**:
- `value`: number - Progress value (0-100)
- `className`: string - Additional CSS classes
- `variant`: 'default' | 'primary' | 'secondary' | 'destructive' | 'warning' | 'success' (default: 'default')
- `showValue`: boolean - Whether to show percentage value (default: true)
- `size`: 'xs' | 'sm' | 'md' | 'lg' | 'xl' (default: 'md')
- `striped`: boolean - Whether to use striped pattern (default: false)
- `animated`: boolean - Whether to animate (default: true)

**Usage**:
```tsx
import { Progress } from '@/components/ui/progress'

<Progress value={75} showValue />
<Progress value={50} variant="success" showValue={false} />
<Progress value={30} variant="warning" striped />
```

#### Spinner
Loading indicator.

**Props**:
- `className`: string - Additional CSS classes
- `variant`: 'default' | 'primary' | 'secondary' | 'destructive' | 'warning' | 'success' (default: 'default')
- `size`: 'xs' | 'sm' | 'md' | 'lg' | 'xl' (default: 'md')
- `thickness`: number - Stroke thickness (default: 2)
- `speed`: number - Animation speed multiplier (default: 1)
- `ariaLabel`: string - ARIA label (default: "Loading")

**Usage**:
```tsx
import { Spinner } from '@/components/ui/spinner'

<Spinner size="lg" />
<Spinner variant="success" size="md" />
<Spinner thickness={3} speed={1.5} />
```

#### Skeleton
Placeholder UI during loading.

**Props**:
- `className`: string - Additional CSS classes
- `variant`: 'text' | 'rectangle' | 'circle' | 'avatar' (default: 'rectangle')
- `width`: number | string - Width (default: 100%)
- `height`: number | string - Height (default: 16px for text, 100% for others)
- `count`: number - Number of skeletons to show (default: 1)
- `animation`: boolean - Whether to show loading animation (default: true)
- `radius`: number | string - Border radius (default: 4px)

**Usage**:
```tsx
import { Skeleton } from '@/components/ui/skeleton'

<Skeleton variant="text" count={3} className="space-y-2" />
<Skeleton variant="avatar" size="lg" />
<Skeleton variant="rectangle" width={200} height={100} count={2} />
```

#### Toast
Temporary notification.

**Props**:
- `children`: React.ReactNode - Toast content
- `className`: string - Additional CSS classes
- `variant`: 'default' | 'primary' | 'secondary' | 'destructive' | 'warning' | 'success' (default: 'default')
- `position`: 'top-start' | 'top-center' | 'top-end' | 'bottom-start' | 'bottom-center' | 'bottom-end' (default: 'bottom-end')
- `duration`: number - Duration in ms before auto-dismiss (default: 5000)
- `closable`: boolean - Whether toast can be manually closed (default: true)
- `onClose`: () => void - Callback when toast is closed
- `onTimeout`: () => void - Callback when toast auto-dismisses
- `icon`: React.ReactNode - Icon to show (optional)
- `progress`: boolean - Whether to show progress bar (default: true)

**Usage**:
```tsx
import { Toast } from '@/components/ui/toast'

// Typically used via a toast manager
<ToastManager position="top-end">
  <Toast 
    variant="success" 
    title="Image Pushed Successfully"
    description="Your image has been pushed to the registry."
    duration={3000}
  />
</ToastManager>
```

#### Validation
Form field validation state.

**Props**:
- `children`: React.ReactNode - Validation message/content
- `className`: string - Additional CSS classes
- `variant`: 'default' | 'primary' | 'secondary' | 'destructive' | 'warning' | 'success' (default: 'default')
- `icon`: React.ReactNode - Icon to show (optional)
- `showIcon`: boolean - Whether to show icon (default: true)
- `visible`: boolean - Whether validation is visible (default: true)

**Usage**:
```tsx
import { Validation } from '@/components/ui/validation'
import { Input } from '@/components/ui/input'

<Input 
  type="password" 
  placeholder="Password" 
  className="mb-2"
/>
<Validation variant="destructive">
  Password must be at least 8 characters long
</Validation>
```

### 4. Input Components
Components for user input.

#### Button
Trigger for actions.

**Props**:
- `children`: React.ReactNode - Button content
- `className`: string - Additional CSS classes
- `variant`: 'default' | 'primary' | 'secondary' | 'destructive' | 'warning' | 'success' | 'outline' | 'ghost' | 'link' (default: 'default')
- `size`: 'xs' | 'sm' | 'md' | 'lg' | 'xl' (default: 'md')
- `disabled`: boolean - Whether button is disabled (default: false)
- `loading`: boolean - Whether to show loading state (default: false)
- `loadingSpinner`: React.ReactNode - Custom loading spinner (optional)
- `type`: 'button' | 'submit' | 'reset' (default: 'button')
- `as`: keyof JSX.IntrinsicElements - HTML element to render (default: "button")
- `onClick`: (event: React.MouseEvent) => void - Click handler
- `href`: string - URL for link variant (required when as="a")
- `target`: string - Target for link variant (optional when as="a")
- `rel`: string - Rel for link variant (optional when as="a")

**Usage**:
```tsx
import { Button } from '@/components/ui/button'

<Button>Default</Button>
<Button variant="destructive">Delete</Button>
<Button variant="outline">Secondary</Button>
<Button variant="ghost">Text</Button>
<Button variant="link" href="/docs">Link</Button>
<Button size="icon">
  <Search size={20} />
</Button>
<Button loading>
  Loading...
</Button>
```

#### Checkbox
Binary selection.

**Props**:
- `checked`: boolean - Whether checkbox is checked
- `className`: string - Additional CSS classes
- `disabled`: boolean - Whether checkbox is disabled (default: false)
- `onChange`: (event: React.ChangeEvent<HTMLInputElement>) => void - Change handler
- `indeterminate`: boolean - Whether to show indeterminate state (default: false)
- `label`: React.ReactNode - Label text (optional)
- `labelPosition`: 'before' | 'after' (default: 'after')
- `className`: string - Additional CSS classes

**Usage**:
```tsx
import { Checkbox } from '@/components/ui/checkbox'

<Checkbox 
  checked={isSubscribed} 
  onChange={(e) => setSubscribed(e.target.checked)}
  label="Subscribe to newsletter"
/>
<Checkbox 
  checked={isIndeterminate} 
  indeterminate 
  label="Partial selection"
/>
```

#### Combobox
Text input with dropdown suggestions.

**Props**:
- `value`: string - Current value
- `className`: string - Additional CSS classes
- `placeholder`: string - Input placeholder (default: "")
- `options`: Array<string> | Array<{ label: string; value: string }> - Suggestion options
- `onChange`: (value: string) => void - Value change handler
- `onSelect`: (option: string | { label: string; value: string }) => void - Option selection handler
- `disabled`: boolean - Whether input is disabled (default: false)
- `filterable`: boolean - Whether to filter options based on input (default: true)
- `matchFromStart`: boolean - Whether to match from start of string (default: true)
- `freeSolo`: boolean - Whether to allow values not in options (default: false)
- `autoComplete`: boolean - Whether to auto-complete with first match (default: false)
- `listboxClassName`: string - CSS classes for listbox
- `listboxMaxHeight`: number | string - Maximum height of listbox (default: 200px)
- `loading`: boolean - Whether to show loading state (default: false)
- `loadingSpinner`: React.ReactNode - Custom loading spinner (optional)
- `open`: boolean - Whether listbox is open (default: false)
- `onOpenChange`: (open: boolean) => void - Open/close handler
- `anchorEl`: HTMLElement - Anchor element for positioning (optional)
- `className`: string - Additional CSS classes

**Usage**:
```tsx
import { Combobox } from '@/components/ui/combobox'

<Combobox
  placeholder="Search repositories..."
  options={['web-app', 'api-service', 'database', 'cache']}
  onChange={(value) => setSearch(value)}
  onSelect={(value) => setSelected(value)}
  freeSolo
/>
```

#### Form
Container for related inputs.

**Props**:
- `children`: React.ReactNode - Form content
- `className`: string - Additional CSS classes
- `onSubmit`: (event: React.FormEvent) => void - Submit handler
- `method`: 'get' | 'post' (default: 'post')
- `action`: string - Form action URL (optional)
- `noValidate`: boolean - Whether to disable HTML5 validation (default: false)
- `role`: string - ARIA role (optional)
- `id`: string - Form ID (optional)

**Usage**:
```tsx
import { Form } from '@/components/ui/form'
import { Field } from '@/components/ui/form'
import { Label } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'

<Form onSubmit={handleSubmit}>
  <Field>
    <Label>Email address</Label>
    <Input type="email" placeholder="you@example.com" />
  </Field>
  
  <Field>
    <Label>Password</Label>
    <Input type="password" placeholder="••••••••" />
  </Field>
  
  <Field className="mb-4">
    <Button type="submit">Sign in</Button>
  </Field>
</Form>
```

#### Input
Text entry field.

**Props**:
- `value`: string - Current value
- `className`: string - Additional CSS classes
- `placeholder`: string - Input placeholder (default: "")
- `type`: 'text' | 'password' | 'email' | 'number' | 'tel' | 'url' | 'search' | 'date' | 'time' | 'datetime-local' | 'month' | 'week' (default: 'text')
- `disabled`: boolean - Whether input is disabled (default: false)
- `readOnly`: boolean - Whether input is read-only (default: false)
- `required`: boolean - Whether input is required (default: false)
- `min`: number | string - Minimum value (for number/date types)
- `max`: number | string - Maximum value (for number/date types)
- `step`: number | string - Step value (for number types)
- `pattern`: string - Regex pattern for validation
- `onChange`: (event: React.ChangeEvent<HTMLInputElement>) => void - Change handler
- `onBlur`: (event: React.FocusEvent<HTMLInputElement>) => void - Blur handler
- `onFocus`: (event: React.FocusEvent<HTMLInputElement>) => void - Focus handler
- `autoComplete`: string - Autocomplete attribute (default: "off")
- `autoFocus`: boolean - Whether to auto-focus (default: false)
- `inputMode`: string - Input mode attribute (optional)
- `showPassword`: boolean - Whether to show password toggle (for password type)
- `inputClassName`: string - CSS classes for input element
- `wrapperClassName`: string - CSS classes for wrapper element
- `label`: React.ReactNode - Label text (optional)
- `labelPosition`: 'before' | 'after' | 'above' (default: 'after')
- `helpText`: React.ReactNode - Help text (optional)
- `validation`: React.ReactNode - Validation message (optional)
- `className`: string - Additional CSS classes

**Usage**:
```tsx
import { Input } from '@/components/ui/input'

<Input 
  type="text" 
  placeholder="Enter your name" 
  label="Name"
/>
<Input 
  type="email" 
  placeholder="you@example.com" 
  label="Email"
  required
/>
<Input 
  type="password" 
  placeholder="••••••••" 
  label="Password"
  showPassword
/>
<Input 
  type="number" 
  placeholder="Enter age" 
  label="Age"
  min={0}
  max={150}
  step={1}
/>
```

#### Label
Text label for form elements.

**Props**:
- `children`: React.ReactNode - Label text
- `className`: string - Additional CSS classes
- `htmlFor`: string - ID of associated input element
- `disabled`: boolean - Whether associated input is disabled
- `required`: boolean - Whether associated input is required
- `position`: 'before' | 'after' | 'above' (default: 'after')

**Usage**:
```tsx
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'

<Label htmlFor="email" required>
  Email address
</Label>
<Input id="email" type="email" placeholder="you@example.com" />
```

#### Radio Group
Mutually exclusive selection.

**Props**:
- `value`: string | number - Selected value
- `options`: Array<{ label: string; value: string | number }> - Radio options
- `className`: string - Additional CSS classes
- `disabled`: boolean - Whether radio group is disabled (default: false)
- `onChange`: (value: string | number) => void - Change handler
- `layout`: 'horizontal' | 'vertical' (default: 'vertical')
- `spacing`: number or string - Space between options (default: 4px)
- `label`: React.ReactNode - Group label (optional)
- `labelPosition`: 'before' | 'after' | 'above' (default: 'after')
- `className`: string - Additional CSS classes

**Usage**:
```tsx
import { RadioGroup } from '@/components/ui/radio-group'

<RadioGroup
  value={size}
  options={[
    { label: 'Small', value: 's' },
    { label: 'Medium', value: 'm' },
    { label: 'Large', value: 'l' }
  ]}
  onChange={(value) => setSize(value)}
  label="Size"
  layout="horizontal"
/>
```

#### Select
Dropdown selection.

**Props**:
- `value`: string | number - Selected value
- `options`: Array<{ label: string; value: string | number }> - Select options
- `className`: string - Additional CSS classes
- `disabled`: boolean - Whether select is disabled (default: false)
- `onChange`: (value: string | number) => void - Change handler
- `placeholder`: string - Placeholder when no value selected (default: "Select...")
- `autoWidth`: boolean - Whether to adjust width based on options (default: false)
- `menuMaxHeight`: number | string - Maximum height of dropdown menu (default: 200px)
- `loading`: boolean - Whether to show loading state (default: false)
- `loadingSpinner`: React.ReactNode - Custom loading spinner (optional)
- `showArrow`: boolean - Whether to show dropdown arrow (default: true)
- `className`: string - Additional CSS classes

**Usage**:
```tsx
import { Select } from '@/components/ui/select'

<Select
  value={role}
  options={[
    { label: 'Admin', value: 'admin' },
    { label: 'User', value: 'user' },
    { label: 'Viewer', value: 'viewer' }
  ]}
  onChange={(value) => setRole(value)}
  placeholder="Select role"
/>
```

#### Switch
Toggle for on/off state.

**Props**:
- `checked`: boolean - Whether switch is checked
- `className`: string - Additional CSS classes
- `disabled`: boolean - Whether switch is disabled (default: false)
- `onChange`: (checked: boolean) => void - Change handler
- `label`: React.ReactNode - Label text (optional)
- `labelPosition`: 'before' | 'after' | 'above' (default: 'after')
- `className`: string - Additional CSS classes

**Usage**:
```tsx
import { Switch } from '@/components/ui/switch'

<Switch 
  checked={isEnabled} 
  onChange={(checked) => setEnabled(checked)}
  label="Enable notifications"
/>
<Switch 
  checked={isDarkMode} 
  onChange={(checked) => setDarkMode(checked)}
  label="Dark mode"
/>
```

#### Textarea
Multi-line text input.

**Props**:
- `value`: string - Current value
- `className`: string - Additional CSS classes
- `placeholder`: string - Input placeholder (default: "")
- `disabled`: boolean - Whether textarea is disabled (default: false)
- `readOnly`: boolean - Whether textarea is read-only (default: false)
- `required`: boolean - Whether textarea is required (default: false)
- `minLength`: number - Minimum length
- `maxLength`: number - Maximum length
- `onChange`: (event: React.ChangeEvent<HTMLTextAreaElement>) => void - Change handler
- `onBlur`: (event: React.FocusEvent<HTMLTextAreaElement>) => void - Blur handler
- `onFocus`: (event: React.FocusEvent<HTMLTextAreaElement>) => void - Focus handler
- `autoComplete`: string - Autocomplete attribute (default: "off")
- `autoFocus`: boolean - Whether to auto-focus (default: false)
- `inputClassName`: string - CSS classes for textarea element
- `wrapperClassName`: string - CSS classes for wrapper element
- `label`: React.ReactNode - Label text (optional)
- `labelPosition`: 'before' | 'after' | 'above' (default: 'after')
- `helpText`: React.ReactNode - Help text (optional)
- `validation`: React.ReactNode - Validation message (optional)
- `className`: string - Additional CSS classes

**Usage**:
```tsx
import { Textarea } from '@/components/ui/textarea'

<Textarea
  placeholder="Enter your bio..."
  label="Bio"
  minLength={10}
  maxLength={500}
  rows={4}
/>
```

### 5. Navigation Components
Components for navigation and wayfinding.

#### Breadcrumb
Hierarchical navigation trail.

**Props**:
- `children`: React.ReactNode - Breadcrumb items
- `className`: string - Additional CSS classes
- `separator`: React.ReactNode - Separator between items (default: "/")
- `showSeparator`: boolean - Whether to show separator (default: true)
- `label`: React.ReactNode - Breadcrumb label (optional)
- `tag`: keyof JSX.IntrinsicElements - HTML element to render (default: "nav")
- `ariaLabel`: string - ARIA label (default: "breadcrumb")

**Usage**:
```tsx
import { Breadcrumb } from '@/components/ui/breadcrumb'
import { BreadcrumbItem } from '@/components/ui/breadcrumb'

<Breadcrumb label="You are here:" separator=" > ">
  <BreadcrumbItem href="/">Home</BreadcrumbItem>
  <BreadcrumbItem href="/repositories">Repositories</BreadcrumbItem>
  <BreadcrumbItem>web-app</BreadcrumbItem>
</Breadcrumb>
```

#### Link
Navigation to another resource.

**Props**:
- `children`: React.ReactNode - Link content
- `className`: string - Additional CSS classes
- `href`: string - URL to navigate to
- `target`: string - Target attribute (default: "_self")
- `rel`: string - Rel attribute (default: "noopener noreferrer")
- `replace`: boolean - Whether to replace current entry in history (default: false)
- `state`: any - State object to push to history (optional)
- `scroll`: boolean - Whether to scroll to top (default: false)
- `prefetch`: boolean - Whether to prefetch (default: false)
- `prefetchPriority`: number - Prefetch priority (default: 0)
- `scroll`: boolean - Whether to scroll to top (default: false)
- `className`: string - Additional CSS classes

**Usage**:
```tsx
import { Link } from '@/components/ui/link'

<Link href="/docs" target="_blank" rel="noopener noreferrer">
  Documentation
</Link>
<Link href="/settings" replace>
  Settings
</Link>
```

#### Menu
List of actions or options.

**Props**:
- `children`: React.ReactNode - Menu items
- `className`: string - Additional CSS classes
- `trigger`: React.ReactNode - Element that triggers the menu
- `alignment`: 'start' | 'center' | 'end' (default: 'end')
- `position`: 'top' | 'bottom' | 'left' | 'right' (default: 'bottom')
- `offset`: number - Distance from trigger (default: 0)
- `skipTrigger`: boolean - Whether to skip trigger element when calculating position (default: false)
- `portalPadding`: number - Padding when using portal (default: 5)
- `sideOffset`: number - Offset from side (default: 0)
- `sidePadding`: number - Padding from side (default: 0)
- `collisionPadding`: number - Padding when colliding with boundaries (default: 5)
- `collisionBoundary`: string - Boundary for collision detection (default: "clippingAncestors")
- `loop`: boolean - Whether to loop focus (default: false)
- `closeOnScroll`: boolean - Whether to close on scroll (default: false)
- `closeOnEscape`: boolean - Whether to close on escape key (default: true)
- `closeOnPointerDown`: boolean - Whether to close on pointer down outside (default: true)
- `closeOnBlur`: boolean - Whether to close on blur (default: true)
- `disablePointerDownEvent`: boolean - Whether to disable pointer down event on trigger (default: false)
- `disableOutsidePointerDownEvent`: boolean - Whether to disable outside pointer down event (default: false)
- `disableBlurEvent`: boolean - Whether to disable blur event on trigger (default: false)
- `portalContainer`: HTMLElement - Container for portal (optional)
- `disablePortal`: boolean - Whether to disable portal (default: false)
- `disableLayers`: boolean - Whether to disable layers (default: false)
- `static`: boolean - Whether to position statically (default: false)
- `className`: string - Additional CSS classes

**Usage**:
```tsx
import { Menu } from '@/components/ui/menu'
import { MenuTrigger } from '@/components/ui/menu'
import { MenuContent } from '@/components/ui/menu'
import { MenuItem } from '@/components/ui/menu'

<Menu>
  <MenuTrigger asChild>
    <Button variant="outline" size="icon">
      <MoreHorizontal size={20} />
    </Button>
  </MenuTrigger>
  <MenuContent align="end" sideOffset={4}>
    <MenuItem>Edit</MenuItem>
    <MenuItem>Delete</MenuItem>
    <MenuItem>Duplicate</MenuItem>
  </MenuContent>
</Menu>
```

#### Pagination
Navigation through paginated content.

**Props**:
- `page`: number - Current page (default: 1)
- `pageCount`: number - Total number of pages
- `className`: string - Additional CSS classes
- `variant`: 'default' | 'primary' | 'secondary' | 'destructive' | 'warning' | 'success' (default: 'default')
- `showFirstLast`: boolean - Whether to show first/last page buttons (default: true)
- `showPrevNext`: boolean - Whether to show previous/next buttons (default: true)
- `showPageNumbers`: boolean - Whether to show page numbers (default: true)
- `hideSinglePage`: boolean - Whether to hide pagination when only one page (default: true)
- `onPageChange`: (page: number) => void - Page change handler
- `showSizeChanger`: boolean - Whether to show page size changer (default: false)
- `pageSize`: number - Current page size (default: 10)
- `pageSizeOptions`: Array<number> - Available page sizes (default: [10, 25, 50, 100])
- `onPageSizeChange`: (size: number) => void - Page size change handler
- `showQuickJumper`: boolean - Whether to show quick jumper (default: false)
- `pageSizeClassName`: string - CSS classes for page size selector
- `pageNumberClassName`: string - CSS classes for page number buttons
- `prevIcon`: React.ReactNode - Custom previous icon (optional)
- `nextIcon`: React.ReactNode - Custom next icon (optional)
- `firstIcon`: React.ReactNode - Custom first icon (optional)
- `lastIcon`: React.ReactNode - Custom last icon (optional)
- `className`: string - Additional CSS classes

**Usage**:
```tsx
import { Pagination } from '@/components/ui/pagination'

<Pagination 
  page={currentPage} 
  pageCount={totalPages} 
  onPageChange={handlePageChange}
  showSizeChanger
  pageSizeOptions={[10, 25, 50, 100]}
/>
```

#### Tabs
Organized content sections.

**Props**:
- `children`: React.ReactNode - Tab content
- `className`: string - Additional CSS classes
- `defaultValue`: string | number - Initially selected tab (default: first tab)
- `value`: string | number - Currently selected tab (controlled)
- `onValueChange`: (value: string | number) => void - Tab change handler
- `orientation`: 'horizontal' | 'vertical' (default: 'horizontal')
- `verticalPosition`: 'top' | 'bottom' (default: 'top')
- `manual`: boolean - Whether to disable automatic value change (default: false)
- `lazy`: boolean - Whether to lazily load tab content (default: false)
- `destroyOnHide`: boolean - Whether to destroy tab content when hidden (default: false)
- `renderTabContent`: (tab: string | number) => React.ReactNode - Custom tab content renderer
- `tabPosition`: 'start' | 'end' | 'top' | 'bottom' (default: 'end' for horizontal, 'bottom' for vertical)
- `domRef`: HTMLElement - DOM reference for positioning (optional)
- `className`: string - Additional CSS classes

**Usage**:
```tsx
import { Tabs } from '@/components/ui/tabs'
import { TabsList } from '@/components/ui/tabs'
import { TabsTrigger } from '@/components/ui/tabs'
import { TabsContent } from '@/components/ui/tabs'

<Tabs defaultValue="overview" orientation="horizontal">
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
```

#### Stepper
Multi-step process indicator.

**Props**:
- `children`: React.ReactNode - Step content
- `className`: string - Additional CSS classes
- `variant`: 'default' | 'primary' | 'secondary' | 'destructive' | 'warning' | 'success' (default: 'default')
- `orientation`: 'horizontal' | 'vertical' (default: 'horizontal')
- `active`: number - Currently active step (0-indexed, default: 0)
- `steps`: Array<{ label: string; description?: string; icon?: React.ReactNode }> - Step definitions
- `showStepNumbers`: boolean - Whether to show step numbers (default: true)
- `showStepDescription`: boolean - Whether to show step descriptions (default: true)
- `connectSteps`: boolean - Whether to connect steps with lines (default: true)
- `size`: 'xs' | 'sm' | 'md' | 'lg' | 'xl' (default: 'md')
- `className`: string - Additional CSS classes

**Usage**:
```tsx
import { Stepper } from '@/components/ui/stepper'

<Stepper 
  active={currentStep}
  steps={[
    { label: 'Account Info', description: 'Enter your personal details' },
    { label: 'Security', description: 'Set up your password and 2FA' },
    { label: 'Preferences', description: 'Choose your notification settings' },
    { label: 'Complete', description: 'Your account is ready to use' }
  ]}
  showStepDescription
  connectSteps
/>
```

### 6. Overlay Components
Components that appear over other content.

#### Dialog
Modal window requiring user action.

**Props**:
- `children`: React.ReactNode - Dialog content
- `className`: string - Additional CSS classes
- `trigger`: React.ReactNode - Element that triggers the dialog
- `open`: boolean - Whether dialog is open (controlled)
- `defaultOpen`: boolean - Whether dialog is open by default (uncontrolled)
- `onOpenChange`: (open: boolean) => void - Open/close handler
- `alignment`: 'center' | 'top' | 'bottom' | 'left' | 'right' (default: 'center')
- `position`: 'center' | 'top' | 'bottom' | 'left' | 'right' (default: 'center')
- `offset`: number - Distance from trigger (default: 0)
- `skipTrigger`: boolean - Whether to skip trigger element when calculating position (default: false)
- `portalPadding`: number - Padding when using portal (default: 5)
- `sideOffset`: number - Offset from side (default: 0)
- `sidePadding`: number - Padding from side (default: 0)
- `collisionPadding`: number - Padding when colliding with boundaries (default: 5)
- `collisionBoundary`: string - Boundary for collision detection (default: "clippingAncestors")
- `loop`: boolean - Whether to loop focus (default: false)
- `closeOnScroll`: boolean - Whether to close on scroll (default: false)
- `closeOnEscape`: boolean - Whether to close on escape key (default: true)
- `closeOnPointerDown`: boolean - Whether to close on pointer down outside (default: true)
- `closeOnBlur`: boolean - Whether to close on blur (default: true)
- `disablePointerDownEvent`: boolean - Whether to disable pointer down event on trigger (default: false)
- `disableOutsidePointerDownEvent`: boolean - Whether to disable outside pointer down event (default: false)
- `disableBlurEvent`: boolean - Whether to disable blur event on trigger (default: false)
- `portalContainer`: HTMLElement - Container for portal (optional)
- `disablePortal`: boolean - Whether to disable portal (default: false)
- `disableLayers`: boolean - Whether to disable layers (default: false)
- `static`: boolean - Whether to position statically (default: false)
- `className`: string - Additional CSS classes

**Usage**:
```tsx
import { Dialog } from '@/components/ui/dialog'
import { DialogTrigger } from '@/components/ui/dialog'
import { DialogContent } from '@/components/ui/dialog'
import { DialogHeader } from '@/components/ui/dialog'
import { DialogTitle } from '@/components/ui/dialog'
import { DialogDescription } from '@/components/ui/dialog'
import { DialogFooter } from '@/components/ui/dialog'

const [open, setOpen] = useState(false)

<>
  <Button onClick={() => setOpen(true)}>
    Delete Repository
  </Button>
  
  <Dialog open={open} onOpenChange={setOpen}>
    <DialogTrigger asChild>
      <Button variant="outline" onClick={() => setOpen(false)}>
        Cancel
      </Button>
    </DialogTrigger>
    <DialogContent className="space-y-4">
      <DialogHeader className="mb-4">
        <DialogTitle className="text-lg font-semibold">
          Delete Repository
        </DialogTitle>
        <DialogDescription className="text-sm text-gray-500">
          Are you sure you want to delete this repository? This action cannot be undone.
        </DialogDescription>
      </DialogHeader>
      <DialogFooter className="flex justify-end space-x-3">
        <Button variant="outline" onClick={() => setOpen(false)}>
          Cancel
        </Button>
        <Button variant="destructive" onClick={() => {
          // Handle deletion
          setOpen(false)
        }}>
          Delete Repository
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</>
```

#### Drawer
Slide-in panel from screen edge.

**Props**:
- `children`: React.ReactNode - Drawer content
- `className`: string - Additional CSS classes
- `trigger`: React.ReactNode - Element that triggers the drawer
- `open`: boolean - Whether drawer is open (controlled)
- `defaultOpen`: boolean - Whether drawer is open by default (uncontrolled)
- `onOpenChange`: (open: boolean) => void - Open/close handler
- `alignment`: 'start' | 'center' | 'end' | 'left' | 'right' (default: 'end')
- `position`: 'start' | 'end' | 'left' | 'right' (default: 'end')
- `offset`: number - Distance from trigger (default: 0)
- `skipTrigger`: boolean - Whether to skip trigger element when calculating position (default: false)
- `portalPadding`: number - Padding when using portal (default: 5)
- `sideOffset`: number - Offset from side (default: 0)
- `sidePadding`: number - Padding from side (default: 0)
- `collisionPadding`: number - Padding when colliding with boundaries (default: 5)
- `collisionBoundary`: string - Boundary for collision detection (default: "clippingAncestors")
- `loop`: boolean - Whether to loop focus (default: false)
- `closeOnScroll`: boolean - Whether to close on scroll (default: false)
- `closeOnEscape`: boolean - Whether to close on escape key (default: true)
- `closeOnPointerDown`: boolean - Whether to close on pointer down outside (default: true)
- `closeOnBlur`: boolean - Whether to close on blur (default: true)
- `disablePointerDownEvent`: boolean - Whether to disable pointer down event on trigger (default: false)
- `disableOutsidePointerDownEvent`: boolean - Whether to disable outside pointer down event (default: false)
- `disableBlurEvent`: boolean - Whether to disable blur event on trigger (default: false)
- `portalContainer`: HTMLElement - Container for portal (optional)
- `disablePortal`: boolean - Whether to disable portal (default: false)
- `disableLayers`: boolean - Whether to disable layers (default: false)
- `static`: boolean - Whether to position statically (default: false)
- `width`: number | string - Drawer width (default: 300px)
- `height`: number | string - Drawer height (default: 100vh)
- `className`: string - Additional CSS classes

**Usage**:
```tsx
import { Drawer } from '@/components/ui/drawer'
import { DrawerTrigger } from '@/components/ui/drawer'
import { DrawerContent } from '@/components/ui/drawer'

const [open, setOpen] = useState(false)

<>
  <Button onClick={() => setOpen(true)}>
    Open Drawer
  </Button>
  
  <Drawer open={open} onOpenChange={setOpen} width={300}>
    <DrawerTrigger asChild>
      <Button variant="outline" onClick={() => setOpen(false)}>
        Close
      </Button>
    </DrawerTrigger>
    <DrawerContent className="p-4">
      <h2 className="text-lg font-semibold mb-4">
        Drawer Title
      </h2>
      <p className="text-gray-600">
        This is the drawer content.
      </p>
    </DrawerContent>
  </Drawer>
</>
```

#### Popover
Temporary content anchored to an element.

**Props**:
- `children`: React.ReactNode - Popover content
- `className`: string - Additional CSS classes
- `trigger`: React.ReactNode - Element that triggers the popover
- `open`: boolean - Whether popover is open (controlled)
- `defaultOpen`: boolean - Whether popover is open by default (uncontrolled)
- `onOpenChange`: (open: boolean) => void - Open/close handler
- `alignment`: 'start' | 'center' | 'end' | 'top' | 'bottom' | 'left' | 'right' (default: 'end')
- `position`: 'start' | 'end' | 'top' | 'bottom' | 'left' | 'right' (default: 'bottom')
- `offset`: number - Distance from trigger (default: 0)
- `skipTrigger`: boolean - Whether to skip trigger element when calculating position (default: false)
- `portalPadding`: number - Padding when using portal (default: 5)
- `sideOffset`: number - Offset from side (default: 0)
- `sidePadding`: number - Padding from side (default: 0)
- `collisionPadding`: number - Padding when colliding with boundaries (default: 5)
- `collisionBoundary`: string - Boundary for collision detection (default: "clippingAncestors")
- `loop`: boolean - Whether to loop focus (default: false)
- `closeOnScroll`: boolean - Whether to close on scroll (default: false)
- `closeOnEscape`: boolean - Whether to close on escape key (default: true)
- `closeOnPointerDown`: boolean - Whether to close on pointer down outside (default: true)
- `closeOnBlur`: boolean - Whether to close on blur (default: true)
- `disablePointerDownEvent`: boolean - Whether to disable pointer down event on trigger (default: false)
- `disableOutsidePointerDownEvent`: boolean - Whether to disable outside pointer down event (default: false)
- `disableBlurEvent`: boolean - Whether to disable blur event on trigger (default: false)
- `portalContainer`: HTMLElement - Container for portal (optional)
- `disablePortal`: boolean - Whether to disable portal (default: false)
- `disableLayers`: boolean - Whether to disable layers (default: false)
- `static`: boolean - Whether to position statically (default: false)
- `className`: string - Additional CSS classes

**Usage**:
```tsx
import { Popover } from '@/components/ui/popover'
import { PopoverTrigger } from '@/components/ui/popover'
import { PopoverContent } from '@/components/ui/popover'

const [open, setOpen] = useState(false)

<>
  <Button ref={triggerRef}>
    Hover me
  </Button>
  
  <Popover 
    open={open} 
    onOpenChange={setOpen} 
    anchorEl={triggerRef.current}
    placement="bottom"
  >
    <PopoverTrigger asChild>
      <Button variant="outline" onClick={() => setOpen(false)}>
        Close
      </Button>
    </PopoverTrigger>
    <PopoverContent className="p-4">
      <h2 className="text-lg font-semibold mb-4">
        Popover Title
      </h2>
      <p className="text-gray-600">
        This is the popover content.
      </p>
    </PopoverContent>
  </Popover>
</>
```

#### Tooltip
Brief contextual information.

**Props**:
- `children`: React.ReactNode - Content that triggers the tooltip
- `content`: React.ReactNode - Tooltip content
- `className`: string - Additional CSS classes
- `contentClassName`: string - CSS classes for tooltip content
- `placement`: 'top' | 'top-start' | 'top-end' | 'bottom' | 'bottom-start' | 'bottom-end' | 'left' | 'left-start' | 'left-end' | 'right' | 'right-start' | 'right-end' (default: 'top')
- `offset`: number - Distance from trigger (default: 4)
- `delay`: number - Delay in ms before showing (default: 0)
- `hideDelay`: number - Delay in ms before hiding (default: 0)
- `arrow`: boolean - Whether to show arrow (default: true)
- `className`: string - Additional CSS classes

**Usage**:
```tsx
import { Tooltip } from '@/components/ui/tooltip'

<Tooltip content="View detailed repository information">
  <Button variant="outline" size="icon">
    <Info size={20} />
  </Button>
</Tooltip>
```

### 7. Data Visualization Components
Components for displaying data visually.

#### Chart
Generic chart container.

**Props**:
- `children`: React.ReactNode - Chart content
- `className`: string - Additional CSS classes
- `type`: 'line' | 'bar' | 'area' | 'pie' | 'donut' | 'scatter' | 'heatmap' | 'radar' (default: 'line')
- `data`: any - Chart data (format depends on type)
- `options`: any - Chart options (format depends on library)
- `library`: 'recharts' | 'chart.js' | 'victory' | 'd3' (default: 'recharts')
- `width`: number | string - Chart width (default: 100%)
- `height`: number | string - Chart height (default: 300px)
- `animation`: boolean - Whether to animate (default: true)
- `animationDuration`: number - Animation duration in ms (default: 1000)
- `animationEasing`: string - Animation easing (default: 'easeInOutQuad')
- `refreshDelay`: number - Delay before refreshing data (default: 0)
- `redraw`: boolean - Whether to redraw on resize (default: true)
- `style`: any - Custom CSS styles (optional)
- `className`: string - Additional CSS classes
- `style`: any - Custom CSS styles (optional)
- `className`: string - Additional CSS classes

**Usage**:
```tsx
import { Chart } from '@/components/ui/chart'

<Chart 
  type="line" 
  data={lineData} 
  options={lineOptions} 
  width={800} 
  height={400}
  animation
/>
```

#### Graph
Node-link visualization.

**Props**:
- `children`: React.ReactNode - Graph content
- `className`: string - Additional CSS classes
- `nodes`: Array<{ id: string; label: string; x: number; y: number; data?: any }> - Graph nodes
- `edges`: Array<{ source: string; target: string; label?: string; data?: any }> - Graph edges
- `layout`: 'force' | 'hierarchical' | 'circular' | 'random' (default: 'force')
- `physics`: boolean - Whether to enable physics simulation (default: true)
- `stabilization`: boolean - Whether to stabilize before rendering (default: true)
- `zoom`: boolean - Whether to enable zoom (default: true)
- `drag`: boolean - Whether to enable drag (default: true)
- `select`: boolean - Whether to enable node selection (default: true)
- `highlight`: boolean - Whether to enable node highlighting (default: true)
- `width`: number | string - Graph width (default: 100%)
- `height`: number | string - Graph height (default: 400px)
- `background`: string - Background color (default: '#ffffff')
- `nodeSize`: number - Node size (default: 20)
- `nodeColor`: string - Node color (default: '#3b82f6')
- `nodeLabelColor`: string - Node label color (default: '#ffffff')
- `edgeWidth`: number - Edge width (default: 2)
- `edgeColor`: string - Edge color (default: '#6b7280')
- `edgeLabelColor`: string - Edge label color (default: '#6b7280')
- `className`: string - Additional CSS classes

**Usage**:
```tsx
import { Graph } from '@/components/ui/graph'

<Graph 
  nodes={[
    { id: '1', label: 'Service A', x: 100, y: 100 },
    { id: '2', label: 'Service B', x: 300, y: 100 },
    { id: '3', label: 'Service C', x: 200, y: 250 }
  ]}
  edges={[
    { source: '1', target: '2', label: 'API' },
    { source: '2', target: '3', label: 'Database' },
    { source: '3', target: '1', label: 'Events' }
  ]}
  layout="force"
  width={400}
  height={300}
/>
```

#### Map
Geographic data display.

**Props**:
- `children`: React.ReactNode - Map content
- `className`: string - Additional CSS classes
- `provider`: 'mapbox' | 'google' | 'leaflet' | 'openlayers' (default: 'leaflet')
- `accessToken`: string - Access token for provider (if required)
- `center`: [number, number] - [latitude, longitude] (default: [0, 0])
- `zoom`: number - Zoom level (default: 2)
- `minZoom`: number - Minimum zoom level (default: 0)
- `maxZoom`: number - Maximum zoom level (default: 18)
- `style`: string - Map style (default: 'default')
- `layers`: Array<any> - Map layers (format depends on provider)
- `markers`: Array<{ position: [number, number]; label: string; data?: any }> - Map markers
- `polylines`: Array<{ positions: Array<[number, number]>; color: string; weight: number }> - Map polylines
- `polygons`: Array<{ positions: Array<[number, number]>; color: string; weight: number }> - Map polygons
- `bounds`: [[number, number], [number, number]] - [[minLat, minLng], [maxLat, maxLng]]
- `zoomControl`: boolean - Whether to show zoom control (default: true)
- `attributionControl`: boolean - Whether to show attribution control (default: true)
- `scrollWheelZoom`: boolean - Whether to enable scroll wheel zoom (default: true)
- `doubleClickZoom`: boolean - Whether to enable double click zoom (default: true)
- `dragging`: boolean - Whether to enable dragging (default: true)
- `touchZoom`: boolean - Whether to enable touch zoom (default: true)
- `boxZoom`: boolean - Whether to enable box zoom (default: false)
- `keyboard`: boolean - Whether to enable keyboard navigation (default: true)
- `className`: string - Additional CSS classes

**Usage**:
```tsx
import { Map } from '@/components/ui/map'

<Map 
  provider="leaflet"
  center={[40.7128, -74.0060]} 
  zoom={12}
  style="streets"
  markers={[
    { position: [40.7128, -74.0060], label: "New York City" },
    { position: [34.0522, -118.2437], label: "Los Angeles" }
  ]}
  width={600}
  height={400}
/>
```

#### Metric
Single important value with context.

**Props**:
- `value`: React.ReactNode - Main value (number, string, etc.)
- `label`: React.ReactNode - Label/description
- `className`: string - Additional CSS classes
- `variant`: 'default' | 'primary' | 'secondary' | 'destructive' | 'warning' | 'success' (default: 'default')
- `size`: 'xs' | 'sm' | 'md' | 'lg' | 'xl' (default: 'md')
- `prefix`: React.ReactNode - Content before value (optional)
- `suffix`: React.ReactNode - Content after value (optional)
- `trend`: React.ReactNode - Trend indicator (optional)
- `trendUp`: boolean - Whether trend is positive (default: false)
- `trendDown`: boolean - Whether trend is negative (default: false)
- `icon`: React.ReactNode - Icon to show (optional)
- `iconPosition`: 'before' | 'after' (default: 'after')
- `className`: string - Additional CSS classes

**Usage**:
```tsx
import { Metric } from '@/components/ui/metric'

<Metric 
  value="1,248" 
  label="Total Images"
  prefix={<span className="text-sm text-gray-500">/</span>}
  suffix={<span className="text-sm text-gray-500">images</span>}
  trend={<span className="text-sm text-green-600">+12%</span>}
  trendUp
/>
<Metric 
  value="0.78" 
  label="Average Trust Score"
  prefix={<span className="text-sm text-gray-500">/</span>}
  suffix={<span className="text-sm text-gray-500">score</span>}
  trend={<span className="text-sm text-green-600">+0.05</span>}
  trendUp
  variant="success"
/>
```

#### Sparkline
Simple line chart for trends.

**Props**:
- `data`: number[] - Data points
- `className`: string - Additional CSS classes
- `color`: string - Line color (default: '#3b82f6')
- `fill`: boolean - Whether to fill under line (default: false)
- `fillColor`: string - Fill color (default: 'currentColor')
- `width`: number | string - Chart width (default: 100%)
- `height`: number | string - Chart height (default: 20px)
- `lineWidth`: number - Line width (default: 2)
- `pointSize`: number - Point size (default: 0)
- `pointColor`: string - Point color (default: 'currentColor')
- `showPoints`: boolean - Whether to show points (default: false)
- `smooth`: boolean - Whether to smooth line (default: false)
- `className": "1,248",
    "label": "Total Images",
    "prefix": <span className="text-sm text-gray-500">/</span>,
    "suffix": <span className="text-sm text-gray-500">images</span>,
    "trend": <span className="text-sm text-green-600">+12%</span>,
    "trendUp": true
  }}
/>
</ArtifactCard>
</div>
</div>
</div>
</CardContent>
<CardFooter className="pt-4">
  <Button variant="outline" size="icon" onClick={() => {
    // Handle pull
  }}>
    <Download size={20} />
  </Button>
  <Button variant="outline" size="icon" onClick={() => {
    // Handle scan
  }}>
    <Zap size={20} />
  </Button>
  <DropdownMenu>
    <DropdownMenuTrigger asChild>
      <Button variant="outline" size="icon">
        <MoreHorizontal size={20} />
      </Button>
    </DropdownMenuTrigger>
    <DropdownMenuContent align="end" sideOffset={4}>
      <DropdownMenuItem>Settings</DropdownMenuItem>
      <DropdownMenuItem>Manage Access</DropdownMenuItem>
      <DropdownMenuItem>Delete Repository</DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>
</CardFooter>
</Card>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
</div