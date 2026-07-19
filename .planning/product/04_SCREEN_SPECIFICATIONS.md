# Kyros Screen Specifications

## Overview
This document details the specifications for each screen in the Kyros platform, including purpose, components, states, interactions, and accessibility considerations.

## 1. Dashboard Screen

### Purpose
Provides an at-a-glance view of system health, key metrics, recent activity, and actionable insights for platform operators and administrators.

### Components
- **Header**: 
  - Application logo and name
  - Global search bar
  - User avatar with dropdown (profile, notifications, settings, logout)
- **Metrics Row** (4 configurable cards):
  - Each card: label, value, trend indicator (up/down/neutral), icon, optional sparkline
  - Default metrics: Total Repositories, Total Images Scanned, Average Trust Score, Security Alerts (24h)
- **Activity Feed** (left column, 2/3 width):
  - Filter controls (event types, time range)
  - List of recent events (pushes, pulls, trust score changes, security events, system events)
  - Each event item: icon, timestamp, description, actor, related resource link
- **Quick Insights** (right column, 1/3 width):
  - Top 5 Repositories by Activity (bar chart)
  - Trust Score Distribution (donut chart)
  - Vulnerability Trend (sparkline with tooltip)
  - System Health Status (service status indicators)

### Empty States
- **No Activity**: "No recent activity to display. Start by pushing an image or inviting a team member."
- **No Data for Insights**: "Insufficient data for insights. Begin using Kyros to populate metrics."

### Loading States
- Skeleton placeholders for metrics cards and charts
- Activity feed shows loading spinner with "Loading recent activity..."

### Error States
- Individual component error: "Unable to load [component]. Please try again later."
- Page-level error: "Failed to load dashboard. [Retry] [Report Issue]"

### Bulk Actions
- Not applicable (dashboard is read-only overview)

### Keyboard Shortcuts
- `/`: Focus global search
- `g d`: Navigate to dashboard
- `?`: Show keyboard shortcuts modal

### Accessibility
- All interactive elements have accessible names
- Color contrast meets WCAG AA
- Charts include accessible data tables as fallback
- Live regions for dynamic content updates

## 2. Repository List Screen

### Purpose
Enables browsing, searching, and managing repositories within namespaces.

### Components
- **Header**:
  - Page title: "Repositories"
  - Namespace selector (dropdown)
  - Search bar with placeholder "Search repositories..."
  - Filter button (opens filter sidebar)
  - "New Repository" button (primary action)
- **Filter Sidebar** (collapsible):
  - Visibility: Public, Private, Protected, Inherited
  - Sort by: Name (A-Z), Name (Z-A), Updated (newest first), Updated (oldest first), Trust Score (high-low), Trust Score (low-high)
  - Tags: Has tags, No tags
  - Size: Min/max range inputs
- **Results Grid/List Toggle**:
  - Grid view (default): Repository cards
  - List view: Compact table with columns
- **Repository Card** (grid view):
  - Repository name and namespace
  - Truncated description (tooltip on hover)
  - Tag count badge
  - Last updated timestamp
  - Trust score visualization (mini-bar or badge)
  - Security status indicator (icon + tooltip)
  - Action menu (pull, scan, settings)
- **Repository Table** (list view):
  - Checkbox for selection
  - Repository name (link to detail)
  - Namespace
  - Description
  - Tag count
  - Last updated
  - Trust score (numeric + color)
  - Security issues (count by severity)
  - Actions (pull, scan, settings, delete)
- **Bulk Action Bar** (appears when items selected):
  - Selected count indicator
  - Actions: Delete selected, Change visibility, Trigger scan, Export list

### Empty States
- **No Repositories**: "No repositories found in this namespace. Create your first repository to get started."
- **No Matches**: "No repositories match your current filters. Try adjusting your search or filters."

### Loading States
- Skeleton cards or table rows while loading
- Footer message: "Loading repositories..."

### Error States
- Inline error: "Failed to load repositories. [Retry]"
- Per-item error: Error icon with tooltip for specific repository

### Bulk Actions
- Delete Selected: Confirmation dialog with list of repositories to delete
- Change Visibility: Dropdown to select new visibility (public/private/protected)
- Trigger Scan: Confirmation to initiate security scan on all selected repositories
- Export List: Download CSV of selected repositories

### Keyboard Shortcuts
- `g r`: Navigate to repositories
- `/`: Focus search
- `f`: Open filter sidebar
- `n`: New repository (focuses name field in modal)
- `Enter`: Open selected repository (when focused on row/card)
- `Space`: Toggle selection (when focused on row/card)
- `Delete`: Delete selected items (after confirmation)

### Accessibility
- All form fields have associated labels
- Keyboard navigable flow: header → filters → results → bulk actions
- ARIA live region for announcement of action results
- Sufficient touch targets for mobile (≥44x44px)

## 3. Repository Detail Screen

### Purpose
Displays comprehensive information about a specific repository and its tags.

### Components
- **Header**:
  - Breadcrumb: Home > Repositories > [namespace] > [repository name]
  - Repository name and namespace (prominent)
  - Description (expandable if truncated)
  - Action buttons: Pull, Scan, Settings (dropdown: Manage Access, Delete)
- **Tabs**:
  - **Overview** (default):
    - Statistics panel: Total tags, total size, push/pull counts (last 30 days)
    - Trust score trend (sparkline over time)
    - Recent activity (last 5 events)
  - **Tags**:
    - Search bar with placeholder "Search tags..."
    - Tag grid/list toggle
    - Tag card (grid): tag name, digest (truncated), size, pushed date, trust score badge, vulnerability count indicators, action menu (pull, delete)
    - Tag table (list): checkbox, tag name, digest, size, pushed, trust score, vulnerabilities, actions
    - Bulk action bar (when items selected): Delete selected, Retag selected, Scan selected
  - **Security**:
    - Vulnerability summary: counts by severity (critical, high, medium, low)
    - SBOM status: present/missing, format, last generated
    - Signature status: verified/unsigned/failed, details
    - Policy compliance: pass/fail/warn, applicable policies
    - Detailed vulnerability table: filter by severity, fix availability, package name
  - **Analytics**:
    - Push/pull volume over time (line chart)
    - Storage usage breakdown (pie chart: by tag, by layer type)
    - Geographic distribution (if available): world map with pull heatmap
    - Usage by user/team (bar chart)

### Empty States
- **No Tags**: "This repository has no tags. Push your first image to get started."
- **No Vulnerabilities**: "No vulnerabilities found in this repository's images."
- **No SBOM**: "No SBOM generated for images in this repository. Enable SBOM generation in settings."

### Loading States
- Tab-specific skeleton loaders
- Placeholder text: "Loading tags..." etc.

### Error States
- Per-section error: "Failed to load [section]. [Retry]"
- Global error: Unable to load repository details.

### Bulk Actions (Tags Tab)
- Delete Selected: Confirmation with list of tags to delete
- Retag Selected: Dialog to specify new tag pattern (e.g., add prefix/suffix)
- Scan Selected: Confirmation to initiate security scan on selected tags

### Keyboard Shortcuts
- `g t`: Navigate to tags tab
- `g s`: Navigate to security tab
- `g a`: Navigate to analytics tab
- `/`: Focus search in current tab
- `Enter`: Open selected tag detail (when focused on card/row)
- `Space`: Toggle selection (when focused on tag item)
- `c`: Copy digest to clipboard (when focused on tag)

### Accessibility
- Tab panel follows ARIA tabs pattern
- All charts have accessible data tables as fallback
- Color is not sole means of conveying information (e.g., trust score badge includes tooltip with numeric value)
- Focus management when opening/closing modals

## 4. Image/Artifact Detail Screen

### Purpose
Shows detailed information about a specific image/artifact including layers, manifest, and security analysis.

### Components
- **Header**:
  - Breadcrumb: Home > Repositories > [namespace] > [repository name] > Tags > [tag name]
  - Image reference: [repository]:[tag] @ [digest]
  - Action buttons: Pull, Scan, Retag, Delete
- **Tabs**:
  - **Overview**:
    - Basic metadata: size, created, pushed, uploaded by
    - Trust score with detailed breakdown (radar chart or bar chart)
    - Layer count and size distribution (treemap or bar chart)
    - Base image information (if detectable)
  - **Layers**:
    - Expandable list of layers (sorted by order in manifest)
    - Each layer: index, size, created by (Dockerfile instruction), digest, security scan results (if available)
    - Visual size breakdown (stacked bar chart)
    - Option to show/hide negligible layers (<1MB)
  - **Manifest**:
    - Syntax-highlighted JSON view of image manifest
    - Sections: schemaVersion, mediaType, config, layers
    - Copy button for each section
    - Option to view as table
  - **Security**:
    - Vulnerability details: filterable table (severity, package, version, fixed version, link to advisory)
    - Affected packages summary
    - SBOM tab: if available, view/download SBOM in SPDX/CycloneDX format
    - Signature verification: details, key ID, verification status
    - License information: detected licenses, compliance status
    - Policy evaluation results: passed/failed/warned policies with details
  - **Dependencies**:
    - Dependency tree/graph view (collapsible)
    - License summary table
    - Outdated dependencies list
    - Known vulnerable dependencies (cross-referenced with vulnerability database)

### Empty States
- **No Layers**: "This image has no layers. This may indicate a manifest issue."
- **No Vulnerabilities**: "No vulnerabilities found in this image."
- **No SBOM**: "SBOM not available for this image. Enable SBOM generation to view components and licenses."
- **No Dependencies**: "No dependencies detected. This may be a minimal or scratch image."

### Loading States
- Section-specific skeleton loaders
- Placeholder: "Analyzing layers..." etc.

### Error States
- Per-section error with retry option
- Manifest error: "Unable to parse image manifest. The image may be corrupted or in an unsupported format."

### Keyboard Shortcuts
- `g l`: Navigate to layers tab
- `g m`: Navigate to manifest tab
- `g s`: Navigate to security tab
- `g d`: Navigate to dependencies tab
- `Enter`: Expand/collapse layer (when focused on layer header)
- `c`: Copy digest to clipboard (when focused on header)
- `+/-`: Increase/decrease font size in manifest view (when focused)

### Accessibility
- Code blocks have accessible labels and support screen reader navigation
- Tree views follow ARIA treeitem pattern
- Color contrast in charts meets WCAG AA
- All interactive elements are keyboard operable

## 5. Trust Score Dashboard

### Purpose
Monitors and analyzes trust scores across the registry to identify trends and outliers.

### Components
- **Header**:
  - Page title: "Trust Score Dashboard"
  - Controls: time range picker (presets: 1h, 6h, 24h, 7d, 30d, custom), refresh button, export button
- **Summary Metrics** (row of cards):
  - Average trust score distribution (histogram):
    - X-axis: score ranges (0.0-0.29, 0.3-0.49, etc.)
    - Y-axis: count of artifacts
    - Color-coded by trust level
    - Tooltip shows exact count and percentage on hover
    - Brush selection to filter detail table
  - Trend over time (line chart):
    - Daily average trust score
    - Optional: percentile bands (25th, 50th, 75th)
    - Tooltip shows date, value, percentile
  - Heatmap (optional): repository vs. time, color intensity = trust score
  - Scatter plot (optional): trust score vs. pull count, size = image size, color = vulnerability count
- **Detailed List**:
  - Table of artifacts with:
    - Repository:Tag (link to detail)
    - Trust score (bar visualization)
    - Trust level (badge with color)
    - Last scanned (timestamp)
    - Vulnerability count (breakdown by severity)
    - Actions: rescan, view details
  - Table features: sorting, column selection, pagination, row selection for bulk actions
  - Bulk action bar: Rescan selected, Export selected

### Empty States
- **No Data**: "No trust score data available for the selected time range. Push an image to begin scoring."
- **No Matches**: "No artifacts match your current filters. Try adjusting your criteria."

### Loading States
- Skeleton charts and table rows
- Message: "Calculating trust scores..."

### Error States
- Chart error: "Unable to load trust score distribution. [Retry]"
- Table error: "Failed to load artifact list. [Retry]"

### Bulk Actions
- Rescan Selected: Confirmation to trigger trust score recalculation on selected artifacts
- Export Selected: Download CSV of selected artifacts with trust score details

### Keyboard Shortcuts
- `g t`: Focus time range picker
- `e`: Focus export button
- `s`: Focus search in detail table
- `Enter`: Open selected artifact detail (when focused on row)
- `Space`: Toggle row selection (when focused on row)
- `c`: Copy trust score to clipboard (when focused on score cell)

### Accessibility
- Charts include accessible data tables as fallback
- Color is not sole means of conveying information (trust level includes text label)
- All interactive elements are keyboard operable
- Sufficient contrast for all text and non-text elements

## 6. Security Center

### Purpose
Centralized view of security posture, vulnerabilities, compliance, and threats.

### Components
- **Header**:
  - Page title: "Security Center"
  - Security status indicator (overall score: green/yellow/red)
  - Time range picker
  - Export button
- **Overview Cards**:
  - Critical Vulnerabilities (24h): count, trend sparkline
  - High Vulnerabilities (24h): count, trend sparkline
  - Policy Violations (24h): count, trend sparkline
  - Images Requiring Attention: count (untrusted or low trust score)
- **Tabs**:
  - **Vulnerabilities**:
    - Filter bar: severity, fix availability, repository, package name, CVE ID
    - Trend over time (line chart: new vulnerabilities per day)
    - Top affected packages (bar chart)
    - Expiration tracking: temporary mitigations with end dates
    - Detailed table: CVE ID, severity, package, version, fixed version, affected images, fix status, actions
    - Bulk actions: Apply temporary mitigation, Notify owners, Export list
  - **Compliance**:
    - Policy compliance status: percentage of images passing all policies
    - Failed evaluations by policy (bar chart)
    - Trend over time (line chart)
    - Remediation tracking: open/closed/resolved issues
    - Detailed table: policy, artifact, result, details, evaluated at, actions
    - Bulk actions: Create ticket, Export list
  - **Configuration**:
    - Security policy management: list of policies with toggle, edit, delete
    - Scanner configuration: enable/disable scanners, adjust settings (frequency, depth)
    - Notification settings: event types, frequency, channels
    - Audit log access: link to audit viewer
  - **Threat Intelligence** (optional):
    - External threat feeds: status, last updated
    - Indicator of compromise (IOC) management: add/edit/delete IOCs
    - Threat hunting queries: saved queries, run on demand

### Empty States
- **No Vulnerabilities**: "No vulnerabilities found in the selected time range. Your images are currently clean!"
- **No Policy Violations**: "No policy violations detected. All images comply with configured security policies."
- **No Threats**: "No threats detected from external feeds."

### Loading States
- Section-specific skeleton loaders
- Placeholder: "Scanning for vulnerabilities..." etc.

### Error States
- Per-tab error with retry option
- Global error: "Unable to load security data. [Retry]"

### Bulk Actions
- Vulnerabilities Tab:
  - Apply Temporary Mitigation: Dialog to select mitigation type and duration
  - Notify Owners: Opens modal to compose message to selected image owners
  - Export List: Download CSV of selected vulnerabilities
- Compliance Tab:
  - Create Ticket: Opens integration with ticketing system (Jira, etc.)
  - Export List: Download CSV of selected compliance issues

### Keyboard Shortcuts
- `g v`: Navigate to vulnerabilities tab
- `g c`: Navigate to compliance tab
- `g n`: Navigate to configuration tab
- `g t`: Focus time range picker
- `f`: Open filter sidebar (in vulnerabilities/compliance tabs)
- `Enter`: Open selected item detail (when focused on row)
- `Space`: Toggle row selection (when focused on row)
- `e`: Focus export button

### Accessibility
- Color is not sole means of conveying information (status includes icon and text)
- All charts have accessible data tables as fallback
- Tables are navigable via keyboard with row/column headers
- Modal dialogs follow ARIA dialog pattern
- Sufficient contrast for status indicators

## 7. User Management Screen

### Purpose
Manage users, groups, roles, and permissions for access control.

### Components
- **Header**:
  - Page title: "Users & Teams"
  - New user button
  - Import users button (from CSV or directory sync)
  - Search bar with placeholder "Search users..."
  - Filter button (opens filter sidebar)
- **Tabs**:
  - **Users**:
    - Search and filter controls
    - User list table:
      - Checkbox for selection
      - Avatar, full name, username, email
      - Roles (badges)
      - Last active
      - Status (active/disabled/locked)
      - Actions: edit, disable, reset MFA, view activity
    - Bulk action bar: Enable selected, Disable selected, Reset MFA for selected, Delete selected
  - **Groups**:
    - Search and filter controls
    - Group list table:
      - Group name, description
      - Member count
      - Created at
      - Actions: edit, view members, delete
    - Bulk action bar: Delete selected
    - Group detail modal: members tab, permissions tab, add/remove members
  - **Roles**:
    - Search and filter controls
    - Role list table:
      - Role name, description
      - Permission count
      - Scope (realm/client)
      - Actions: edit, view permissions, delete
    - Bulk action bar: Delete selected
    - Role detail modal: permissions tab, assign users/groups tab
  - **Permissions**:
    - Search and filter controls
    - Permission catalog table:
      - Permission name, description
      - Assigned to (roles count)
      - Actions: view usage, assign to role
    - Permission detail modal: description, assigned roles, usage examples

### Empty States
- **No Users**: "No users found. Invite your team members to get started."
- **No Groups**: "No groups created. Create a group to organize users."
- **No Roles**: "No custom roles defined. Create a role to define specific permissions."
- **No Permissions**: "No permissions available. This should not occur - contact support."

### Loading States
- Skeleton table rows
- Placeholder: "Loading users..."

### Error States
- Per-tab error with retry option
- Inline error for failed operations (e.g., "Failed to create user: [reason]")

### Bulk Actions
- Enable Selected: Confirmation to enable selected disabled users
- Disable Selected: Confirmation to disable selected active users
- Reset MFA for Selected: Confirmation to reset MFA for selected users
- Delete Selected: Confirmation to delete selected users (with option to transfer ownership)
- Delete Selected Groups: Confirmation to delete selected groups
- Delete Selected Roles: Confirmation to delete selected roles (reassign users first)

### Keyboard Shortcuts
- `g u`: Navigate to users tab
- `g g`: Navigate to groups tab
- `g r`: Navigate to roles tab
- `g p`: Navigate to permissions tab
- `/`: Focus search in current tab
- `n`: New user (focuses username field in modal)
- `Enter`: Open selected item detail (when focused on row)
- `Space`: Toggle row selection (when focused on row)
- `d`: Disable selected users (when in users tab with selection)
- `e`: Enable selected users (when in users tab with selection)

### Accessibility
- Form fields have associated labels
- Error messages are announced via ARIA live region
- Tables are navigable via keyboard with proper headers
- Modal dialogs trap focus and return focus to trigger on close
- Color is not sole means of conveying status (e.g., user status includes text label)

## 8. Settings Screen

### Purpose
Configures system-wide settings and preferences.

### Components
- **Sidebar Navigation** (collapsible):
  - General
  - Authentication
  - Authorization
  - Registry
  - Security
  - Observability
  - Integrations
  - License & Usage
- **Content Area**:
  - Form-based interface for selected category
  - Save changes button (top and bottom)
  - Reset to defaults button (section-specific)
  - Section description and help text

### Sections
#### General
- Instance name and description
- Contact information (email, phone)
- Language and timezone
- Maintenance window (start time, duration)
- Enable/disable self-registration
- Default namespace visibility

#### Authentication
- Identity provider selection (Keycloak, LDAP, SAML, etc.)
- Connection details (test connection button)
- Password policy: min length, require uppercase/lowercase/numbers/symbols, expiration, history
- MFA requirements: enforce for admins, all users, none
- Session settings: idle timeout, absolute timeout, remember me
- Single sign-out enabled/disabled

#### Authorization
- Default role for new users
- Permission inheritance rules
- OPA policy evaluation settings (cache TTL, etc.)
- Enable/disable policy-based access control
- Super user bypass (enable/disable)

#### Registry
- Default storage backend (MinIO, S3, etc.)
- Storage connection details (test connection)
- Retention policies: tag expiration, manifest retention, blob garbage collection schedule
- Image signing: enable/disable, default key management (Cosign/Fulcio)
- Mirror and proxy configuration: upstream registries, routing rules

#### Security
- Vulnerability scanner selection (Trivy, Grype, etc.)
- Scanner configuration: frequency, depth, severity thresholds
- Trust score weights: adjustable sliders for each factor (must sum to 1.0)
- Trust score thresholds: auto-block, manual review, auto-approve
- Policy management: link to policy builder
- Audit settings: retention period, export format, integrity checking

#### Observability
- Metrics collection: enable/disable, endpoint, interval
- Logging: level, format, external endpoint (Loki, Elasticsearch)
- Tracing: enable/disable, endpoint, sample rate
- Health checks: endpoints, timeouts, dependencies
- Alerting: notification channels, escalation policies, silencing rules

#### Integrations
- Webhook templates: predefined payloads for common systems (Slack, Teams, Jira)
- API key management: create, view, revoke keys
- Plugin marketplace: browse, install, configure extensions
- SSO provisioning: Just-in-time user provisioning from identity provider

#### License & Usage
- Current license: type, status, expiration date
- Install license button (file upload or key entry)
- Usage statistics: active users, storage consumption, bandwidth, API calls
- Usage alerts: thresholds and notification methods
- Export usage data: CSV, JSON

### Empty States
- Not applicable (settings are always present)

### Loading States
- Section-specific skeleton forms
- Placeholder: "Loading settings..."

### Error States
- Inline field validation errors (red border, message below)
- Section-level error: "Failed to load [section] settings. [Retry]"
- Save error: "Failed to save settings. Please correct the errors and try again."

### Bulk Actions
- Not applicable (settings are instance-wide)

### Keyboard Shortcuts
- `g s`: Navigate to settings
- `1`-`8`: Jump to section by number (1=General, 2=Authentication, etc.)
- `s`: Focus save button
- `r`: Focus reset button
- `Tab`: Navigate between form fields
- `Enter`: Submit form (when focus is on save button or in a field with no submit action)

### Accessibility
- Form labels are properly associated with inputs
- Error messages are announced and associated with fields
- All interactive elements are keyboard operable
- Color is not sole means of conveying error state (includes icon and text)
- Sufficient contrast for all text and background combinations

## 9. User Profile & Settings Screen

### Purpose
Allows users to manage their personal account, security, notifications, and activity.

### Components
- **Header**:
  - User avatar and name
  - Edit profile button
- **Tabs**:
  - **Profile**:
    - Avatar upload/change (with preview)
    - Full name, email, display name (editable fields)
    - Biography (textarea)
    - Timezone and language preferences (dropdowns)
  - **Security**:
    - Password change: current password, new password, confirm
    - MFA management: enable/disable, regenerate backup codes, trusted devices
    - Session management: list of active sessions (location, device, IP, last active, terminate button)
    - Connected applications: list of OAuth applications with permissions, revoke button
  - **Notifications**:
    - Email notification preferences: toggles for each event type with frequency options (immediate, daily digest, weekly)
    - In-app notification settings: sound, desktop notifications, notification lifetime
    - Webhook URLs for personal notifications: add/edit/delete URLs
  - **Activity Log**:
    - Filter controls: action type, date range
    - Timeline of user actions: icon, timestamp, description, resource link
    - Export activity log button

### Empty States
- **No Activity**: "No recent activity to display. Your actions will appear here as you use Kyros."
- **No Connected Applications**: "You haven't connected any external applications yet."
- **No Notification Webhooks**: "You haven't configured any webhook URLs for personal notifications."

### Loading States
- Skeleton sections
- Placeholder: "Loading profile..."

### Error States
- Per-field validation errors (e.g., "Passwords do not match")
- Section-level error: "Failed to load [section]. [Retry]"
- Save error: "Failed to save changes. Please correct the errors and try again."

### Bulk Actions
- Not applicable (user-specific)

### Keyboard Shortcuts
- `g p`: Navigate to profile tab
- `g s`: Navigate to security tab
- `g n`: Navigate to notifications tab
- `g a`: Navigate to activity tab
- `/`: Focus search in activity log
- `s`: Focus save button (in current tab)
- `Enter`: Save changes (when focus is on save button)
- `Escape`: Close edit mode (in profile tab)

### Accessibility
- Form labels are properly associated with inputs
- Error messages are announced and associated with fields
- Avatar upload area is keyboard operable and has clear instructions
- All interactive elements are keyboard operable
- Color is not sole means of conveying state (e.g., password strength includes text indicator)

## 10. Audit Viewer Screen

### Purpose
Browses, searches, and exports audit events for compliance and investigation.

### Components
- **Header**:
  - Page title: "Audit Log"
  - Time range picker (presets: 24h, 7d, 30d, 90d, 1y, custom)
  - Export button (dropdown: CSV, JSON, PDF)
  - Search bar with placeholder "Search audit events..."
  - Filter button (opens filter sidebar)
- **Filter Sidebar**:
  - Actor: user or service account (searchable)
  - Action type: dropdown of event categories (user management, repository changes, policy changes, etc.)
  - Resource type: dropdown (user, repository, policy, etc.)
  - Outcome: success, failure, both
  - IP address/CIDR range
  - User agent contains
- **Results Table**:
  - Checkbox for selection
  - Timestamp
  - Actor (name/username, type icon)
  - Action (verb + object)
  - Resource (name/type)
  - Outcome (icon + text)
  - IP address
  - Summary of changes (truncated, expandable to show diff)
  - Actions: view details
- **Detail Modal** (when clicking view details):
  - Complete event JSON (syntax highlighted)
  - Tabbed view: formatted fields, raw JSON, changes diff
  - Close button

### Empty States
- **No Events**: "No audit events found for the selected time range and filters."
- **No Matches**: "No audit events match your current filters. Try adjusting your criteria."

### Loading States
- Skeleton table rows
- Message: "Loading audit events..."

### Error States
- Inline error: "Failed to load audit events. [Retry]"
- Per-row error: Error icon with tooltip for specific event

### Bulk Actions
- Export Selected: Download selected events in chosen format
- Delete Selected: (if retention policy allows) Confirmation to delete selected audit events
- Notify on Match: Set up alert for similar future events

### Keyboard Shortcuts
- `g a`: Navigate to audit viewer
- `t`: Focus time range picker
- `/`: Focus search
- `f`: Open filter sidebar
- `Enter`: Open selected event detail (when focused on row)
- `Space`: Toggle row selection (when focused on row)
- `e`: Focus export button
- `d`: Delete selected (when selection exists and permitted)

### Accessibility
- Table is navigable via keyboard with proper headers
- Detail modal follows ARIA dialog pattern
- Date picker is keyboard operable
- All interactive elements are keyboard operable
- Color is not sole means of conveying outcome (includes icon and text)

## 11. Policy Builder Screen

### Purpose
Creates and manages Open Policy Agent (OPA) policies for trust scoring and access control.

### Components
- **Header**:
  - Page title: "Policy Builder"
  - New policy button
  - Import/export buttons
  - Search bar with placeholder "Search policies..."
  - Filter button (opens filter sidebar)
- **Policy List**:
  - Card view (default) or list view toggle
  - Policy card:
    - Name and description
    - Type: trust score, access control, etc.
    - Scope: global, namespace, repository
    - Enabled/disabled toggle
    - Last updated
    - Actions: edit, view, test, delete
  - Policy table (list view):
    - Checkbox for selection
    - Name
    - Type
    - Scope
    - Enabled
    - Last updated
    - Actions: edit, view, test, delete
- **Policy Editor** (when creating/editing):
  - Form fields:
    - Name (required, unique)
    - Description (optional)
    - Type: dropdown (trust score, access control, SBOM validation, etc.)
    - Scope: radio buttons (global, namespace, repository) with dependent fields for namespace/repository selection
    - Enabled: toggle
    - Rules: ACE editor with syntax highlighting for Rego
    - Test cases: section to define input JSON and expected output
  - Sidebar:
    - Input schema viewer (based on type and scope)
    - Output schema viewer
    - Test results panel
    - Help and examples accordion
  - Buttons: Cancel, Save changes, Test policy

### Empty States
- **No Policies**: "No policies defined. Create your first policy to start enforcing rules."
- **No Matches**: "No policies match your current filters. Try adjusting your criteria."

### Loading States
- Skeleton policy cards or table rows
- Placeholder: "Loading policies..."

### Error States
- Inline validation errors (e.g., "Policy name already exists")
- Editor errors: highlighted syntax errors with explanations
- Test failure: "Test failed: expected [X], got [Y]"
- Save error: "Failed to save policy. Please correct the errors and try again."

### Bulk Actions
- Delete Selected: Confirmation to delete selected policies
- Enable Selected: Confirmation to enable selected disabled policies
- Disable Selected: Confirmation to disable selected enabled policies
- Export Selected: Download selected policies as Rego files
- Test Selected: Run test suite on selected policies

### Keyboard Shortcuts
- `g p`: Navigate to policy builder
- `n`: New policy (focuses name field)
- `/`: Focus search in policy list
- `f`: Open filter sidebar
- `e`: Focus edit button (when policy selected)
- `v`: Focus view button (when policy selected)
- `t`: Focus test button (when policy selected)
- `d`: Focus delete button (when policy selected)
- `Enter`: Save changes (in editor when focus is on save button)
- `Escape`: Cancel edit (in editor)

### Accessibility
- Form labels are properly associated with inputs
- Editor is keyboard operable and has accessible label
- Syntax errors are announced via ARIA live region
- All interactive elements are keyboard operable
- Color is not sole means of conveying error state (includes icon and text)
- Sufficient contrast for all text and background combinations

## 12. Webhook Builder Screen

### Purpose
Creates and manages webhook subscriptions for event notifications.

### Components
- **Header**:
  - Page title: "Webhook Builder"
  - New webhook button
  - Import/export buttons
  - Search bar with placeholder "Search webhooks..."
  - Filter button (opens filter sidebar)
- **Webhook List**:
  - Card view (default) or list view toggle
  - Webhook card:
    - Name and description
    - URL (truncated with tooltip)
    - Events: list of subscribed event types
    - Enabled/disabled toggle
    - Last triggered
    - Failure count
    - Actions: edit, view, test, delete
  - Webhook table (list view):
    - Checkbox for selection
    - Name
    - URL
    - Events (count)
    - Enabled
    - Last triggered
    - Failure count
    - Actions: edit, view, test, delete
- **Webhook Editor** (when creating/editing):
  - Form fields:
    - Name (required, unique)
    - Description (optional)
    - URL (required, valid HTTP/HTTPS URL)
    - Events: multi-select checkboxes (with searchable dropdown)
    - Secret: optional field for HMAC-SHA256 validation (with show/hide toggle)
    - Format: radio buttons (JSON, form-urlencoded)
    - Headers: key-value table (add/remove rows)
    - Enabled: toggle
  - Sidebar:
    - Example payload viewer (based on selected events)
    - Test delivery section: 
      - Test event type dropdown
      - Send test button
      - Response status, headers, body (collapsible)
    - Help and examples accordion
  - Buttons: Cancel, Save changes, Test delivery

### Empty States
- **No Webhooks**: "No webhooks defined. Create your first webhook to start receiving notifications."
- **No Matches**: "No webhooks match your current filters. Try adjusting your criteria."

### Loading States
- Skeleton webhook cards or table rows
- Placeholder: "Loading webhooks..."

### Error States
- Inline validation errors (e.g., "Invalid URL format")
- Duplicate name error
- Test delivery failure: "Test delivery failed: [status code] [reason]"
- Save error: "Failed to save webhook. Please correct the errors and try again."

### Bulk Actions
- Delete Selected: Confirmation to delete selected webhooks
- Enable Selected: Confirmation to enable selected disabled webhooks
- Disable Selected: Confirmation to disable selected enabled webhooks
- Export Selected: Download selected webhooks as JSON configurations
- Test Selected: Send test event to selected webhooks

### Keyboard Shortcuts
- `g w`: Navigate to webhook builder
- `n`: New webhook (focuses name field)
- `/`: Focus search in webhook list
- `f`: Open filter sidebar
- `e`: Focus edit button (when webhook selected)
- `v`: Focus view button (when webhook selected)
- `t`: Focus test button (when webhook selected)
- `d`: Focus delete button (when webhook selected)
- `Enter`: Save changes (in editor when focus is on save button)
- `Escape`: Cancel edit (in editor)

### Accessibility
- Form labels are properly associated with inputs
- URL field has clear error messaging for invalid formats
- Secret field has show/hide toggle with accessible label
- All interactive elements are keyboard operable
- Color is not sole means of conveying state (e.g., enabled status includes toggle position and text)
- Sufficient contrast for all text and background combinations

## 13. Dashboard Widgets (Reusable Components)

### Purpose
Defines reusable dashboard widgets that can be configured and placed on custom dashboards.

### Components
- **Widget Types**:
  - **Metric Card**: 
    - Configurable: label, value source (metric query), trend calculation, threshold alerts
    - States: normal, warning, critical, no data
    - Interactions: click to drill down to related view
  - **Time Series Chart**:
    - Configurable: title, metric queries, time range, chart type (line, area, bar), stacking
    - Features: zoom, pan, tooltip, legend, drilldown on click
    - Threshold bands: configurable warning/critical zones
  - **Pie/Donut Chart**:
    - Configurable: title, metric query, dimension, measure
    - Features: tooltip, legend, click to filter
  - **Table**:
    - Configurable: title, columns (field, label, format, sortable, filterable), row actions
    - Features: sorting, filtering, pagination, column resizing, export
  - **List**:
    - Configurable: title, item template (with fields and actions), empty state
    - Features: infinite scroll, load more button, item selection for bulk actions
  - **Heatmap**:
    - Configurable: title, metric query, row dimension, column dimension, color scale
    - Features: tooltip, zoom, click to filter
  - **Scatter Plot**:
    - Configurable: title, X-axis query, Y-axis query, size query, color query
    - Features: tooltip, zoom, pan, click to filter, legend
  - **Status Grid**:
    - Configurable: title, items (name, status source, icon, tooltip)
    - Features: click to drill down, refresh interval
  - **Text Box**:
    - Configurable: title, content (static text or templated with variables)
    - Features: markdown support, link auto-detection

### Widget Configuration Modal
- **Header**: Widget type, title field
- **Body**:
  - Data source section: query builder or direct PromQL/SQL input
  - Display options: varies by widget type
  - Interactions: click actions (navigate to URL, run JavaScript, show modal)
  - Advanced: refresh interval, height, visibility conditions
- **Footer**: Cancel, Save buttons

### Empty States
- **No Data**: "No data available for the selected query and time range."
- **Query Error**: "Invalid query: [error message]"

### Loading States
- Skeleton placeholder matching widget shape
- Message: "Loading data..."

### Error States
- Display error: "Error loading data: [error message] [Retry]"
- Partial data: Show available data with warning banner

### Keyboard Shortcuts
- Not applicable (widgets are configured via modal)

### Accessibility
- All widgets are keyboard operable when focused
- Charts include accessible data tables as fallback
- Color is not sole means of conveying information
- Sufficient contrast for all text and non-text elements
- Screen reader friendly labels and descriptions

## Cross-Screen Consistency

### Navigation Patterns
- Global navigation: persistent sidebar with icons and labels
- Local navigation: tabs for sectioning content within a page
- Breadcrumbs: show hierarchical path for context-dependent pages
- Action buttons: primary (solid), secondary (outline), destructive (red)

### Form Patterns
- Required fields: marked with asterisk (*)
- Validation: inline, real-time where possible
- Error messages: specific, actionable, non-technical
- Help text: available via tooltip or icon
- Defaults: sensible defaults provided where applicable

### Data Presentation Patterns
- Tables: sortable, filterable, paginatible, column selectable
- Lists: infinite scroll or load more, selection for bulk actions
- Cards: consistent padding, border radius, shadow, hover state
- Badges: used for counts, status levels, tags
- Progress: used for completion percentages, upload/download progress

### Feedback Patterns
- Toast notifications: temporary, non-blocking, auto-dismiss
- Modal dialogs: blocking, require action, centered on screen
- Inline validation: immediate feedback on field validity
- Page-level banners: persistent until dismissed or condition resolved

### Error Handling Patterns
- Retry actions: available for transient failures
- Undo: available for recent destructive actions where possible
- Error boundaries: isolate component failures to prevent page crashes
- Fallback UI: show degraded functionality when possible

### Accessibility Patterns
- ARIA labels and roles for all interactive elements
- Keyboard navigable flow: logical tab order
- Focus management: returns focus to triggering element after modal/dialog close
- Skip navigation: "Skip to main content" link
- Landmark regions: header, navigation, main, complementary, footer
- Responsive design: adapts to different screen sizes and orientations

## Conclusion
These screen specifications provide a detailed blueprint for the Kyros platform's user interface. By adhering to these specifications, the development team ensures a consistent, accessible, and user-friendly experience that meets the needs of platform engineers, developers, security engineers, organization administrators, and system administrators.

Each screen is designed with clear purpose, intuitive components, appropriate states for loading/error/empty conditions, and thoughtful interactions. The specifications prioritize accessibility, usability, and maintainability while enabling powerful functionality for managing the software supply chain.