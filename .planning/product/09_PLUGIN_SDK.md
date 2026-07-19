# Kyros Plugin SDK

## Overview
The Kyros Plugin SDK provides a standardized way to extend the functionality of the Kyros platform without modifying the core codebase. Plugins can be developed independently and loaded at runtime to add new features, integrate with external systems, or customize existing behavior.

## Design Principles
1. **Loose Coupling**: Plugins interact with Kyros through well-defined interfaces, minimizing dependencies on internal implementation details.
2. **Discoverability**: The platform can automatically discover and load plugins from a designated directory or registry.
3. **Versioning**: Plugins are versioned independently and can declare compatibility with specific Kyros versions.
4. **Safety**: Plugins run in isolated environments with restricted access to prevent harm to the core system.
5. **Lifecycle Management**: The platform manages the plugin lifecycle (loading, initialization, activation, deactivation, unloading).
6. **Configuration**: Plugins can be configured via the Kyros configuration system or their own configuration files.
7. **Observability**: Plugins emit metrics, logs, and traces that are integrated with the platform's observability system.
8. **Security**: Plugins are subject to the same security policies as core components, with additional sandboxing where appropriate.

## Plugin Types
Kyros supports several types of plugins, each extending a different aspect of the platform:

### 1. Scanner Plugins
Extend the vulnerability scanning capabilities by adding new scanners or modifying scanning behavior.

### 2. Storage Plugins
Add support for new blob storage backends (e.g., Azure Blob Storage, Google Cloud Storage, on-premises NAS).

### 3. Notification Plugins
Add new notification channels (e.g., Slack, Microsoft Teams, PagerDuty, custom webhooks).

### 4. Authentication Plugins
Add new authentication methods (e.g., LDAP, SAML, OAuth providers) or integrate with external identity systems.

### 5. Trust Engine Plugins
Extend the trust score engine with new analysis modules, scoring algorithms, or policy engines.

### 6. Policy Engine Plugins
Add new policy languages or extend the existing OPA policy engine with custom functions.

### 7. Analytics Plugins
Add new metrics, dashboards, or analytical capabilities.

### 8. Storage Driver Plugins
Extend the registry's storage driver interface to support new storage technologies.

### 9. UI Plugins
Extend the web interface with new pages, widgets, or custom components.

### 10. CLI Plugins
Extend the command-line interface with new commands.

## Plugin Interface

### Base Plugin Interface
All plugins must implement the base `Plugin` interface:

```go
type Plugin interface {
    // ID returns the unique identifier for the plugin
    ID() string
    
    // Name returns the human-readable name of the plugin
    Name() string
    
    // Version returns the plugin version (semver)
    Version() string
    
    // Description returns a brief description of the plugin
    Description() string
    
    // KyrosVersion returns the minimum Kyros version required
    KyrosVersion() string
    
    // Init initializes the plugin with the provided configuration
    Init(config *Config) error
    
    // Start starts the plugin (after all plugins have been initialized)
    Start() error
    
    // Stop stops the plugin gracefully
    Stop() error
    
    // Destroy cleans up any resources held by the plugin
    Destroy() error
}
```

### Configuration
Plugins receive a configuration object during initialization. The configuration is specific to the plugin type and is validated by the plugin.

```go
type Config struct {
    // Raw configuration as provided by the user (YAML, JSON, etc.)
    Raw map[string]interface{}
    
    // Get retrieves a configuration value by key path
    Get(path string) (interface{}, error)
    
    // GetString retrieves a string configuration value
    GetString(path string) (string, error)
    
    // GetInt retrieves an integer configuration value
    GetInt(path string) (int, error)
    
    // GetBool retrieves a boolean configuration value
    GetBool(path string) (bool, error)
    
    // GetSlice retrieves a slice configuration value
    GetSlice(path string) ([]interface{}, error)
    
    // GetMap retrieves a map configuration value
    GetMap(path string) (map[string]interface{}, error)
}
```

## Plugin Lifecycle

### 1. Discovery
Plugins are discovered from:
- **File System**: A designated plugins directory (e.g., `./plugins/`)
- **Registry**: A plugin registry (e.g., GitHub releases, custom plugin registry)
- **Built-in**: Plugins bundled with the Kyros distribution

### 2. Loading
The plugin loader:
1. Scans the discovery sources for plugin manifests
2. Validates plugin compatibility with the current Kyros version
3. Loads the plugin binary or source code
4. Creates an instance of the plugin

### 3. Initialization
Each plugin's `Init` method is called with its configuration. Plugins should:
- Validate their configuration
- Set up internal state
- Prepare resources (but not start processing)
- Return an error if initialization fails

### 4. Starting
After all plugins are initialized, the `Start` method is called for each plugin. Plugins should:
- Begin their main processing loop
- Establish connections to external services
- Start listening for events
- Return an error if startup fails

### 5. Running
During normal operation, plugins:
- Process events from the event stream (if applicable)
- Respond to API requests (if they expose endpoints)
- Perform background tasks
- Emit metrics, logs, and traces

### 6. Stopping
When the platform is shutting down or a plugin is being disabled, the `Stop` method is called. Plugins should:
- Stop accepting new work
- Gracefully shut down processing loops
- Close external connections
- Wait for ongoing work to complete (within a timeout)
- Return an error if stopping fails

### 7. Destroying
After stopping, the `Destroy` method is called to clean up resources. Plugins should:
- Release all held resources
- Close any remaining connections
- Clean up temporary files
- Return an error if destruction fails

## Plugin Manifest
Each plugin must include a manifest file (`plugin.yaml` or `plugin.json`) that provides metadata about the plugin.

### Example Manifest (YAML)
```yaml
id: com.example.kyros.plugin.trivy-scanner
name: Trivy Vulnerability Scanner
version: 1.2.0
description: Integrates Trivy vulnerability scanner with Kyros trust score engine
kyrosVersion: ">=1.0.0"
main: github.com/example/kyros-plugin-trivy-scanner
config:
  - name: scan-depth
    type: string
    description: Depth of scanning (full, quick)
    default: "full"
    required: false
  - name: ignore-unfixed
    type: boolean
    description: Ignore unfixed vulnerabilities
    default: false
    required: false
dependencies:
  - name: trivy
    version: ">=0.20.0"
    type: executable
```

### Manifest Fields
- **id**: Unique identifier (reverse DNS notation recommended)
- **name**: Human-readable name
- **version**: Plugin version (semver)
- **description**: Brief description
- **kyrosVersion**: Minimum Kyros version required (semver constraint)
- **main**: Entry point (Go package path, executable path, etc.)
- **config**: Configuration schema (array of config items)
- **dependencies**: External dependencies (executables, libraries, etc.)
- **author**: Plugin author (optional)
- **homepage**: Plugin homepage URL (optional)
- **license**: Plugin license (optional)

## Plugin Isolation

### Process Isolation
Plugins run in separate processes from the core Kyros services to prevent crashes and security issues from affecting the core system.

#### Communication Mechanisms
- **gRPC**: Primary communication mechanism for plugin-core interaction
- **Shared Memory**: For high-performance data transfer (optional)
- **Event Stream**: Plugins can publish and consume events via NATS JetStream
- **Shared Database**: Plugins can access the PostgreSQL database (with appropriate permissions)

### Sandboxing
Plugins operate with restricted privileges:
- **Filesystem Access**: Limited to plugin directory and temporary directories
- **Network Access**: Configurable allow-list for outbound connections
- **System Calls**: Restricted via seccomp, AppArmor, or similar mechanisms
- **Resources**: CPU and memory limits enforced via cgroups or container runtime

### Trust Levels
Plugins can be assigned different trust levels:
- **Trusted**: Bundled with Kyros, full access
- **Verified**: Signed by trusted authority, moderate restrictions
- **Untrusted**: From unknown sources, maximum restrictions

## Plugin Development Guidelines

### Language Support
Plugins can be developed in any language that can communicate via gRPC or consume/produce events via NATS. Official SDKs are provided for:
- **Go**: Official Kyros Plugin SDK
- **Python**: Community-supported SDK
- **JavaScript/TypeScript**: Community-supported SDK
- **Java**: Community-supported SDK
- **Rust**: Community-supported SDK

### Go Plugin SDK
The official Go SDK provides helper interfaces and base implementations.

#### Base Plugin Implementation
```go
type BasePlugin struct {
    ID   string
    Name string
    Version string
    Description string
    KyrosVersion string
    config *Config
}

func (p *BasePlugin) ID() string { return p.ID }
func (p *BasePlugin) Name() string { return p.Name }
func (p *BasePlugin) Version() string { return p.Version }
func (p *BasePlugin) Description() string { return p.Description }
func (p *BasePlugin) KyrosVersion() string { return p.KyrosVersion }
func (p *BasePlugin) Init(config *Config) error {
    p.config = config
    return nil
}
func (p *BasePlugin) Start() error { return nil }
func (p *BasePlugin) Stop() error { return nil }
func (p *BasePlugin) Destroy() error { return nil }
```

#### Scanner Plugin Example
```go
type TrivyScannerPlugin struct {
    BasePlugin
    client *trivy.Client
    config *TrivyScannerConfig
}

func NewTrivyScannerPlugin() *TrivyScannerPlugin {
    return &TrivyScannerPlugin{
        BasePlugin: BasePlugin{
            ID:          "com.example.kyros.plugin.trivy-scanner",
            Name:        "Trivy Vulnerability Scanner",
            Version:     "1.2.0",
            Description: "Integrates Trivy vulnerability scanner with Kyros trust score engine",
            KyrosVersion: ">=1.0.0",
        },
    }
}

func (p *TrivyScannerPlugin) Init(config *Config) error {
    // Parse configuration
    p.config = &TrivyScannerConfig{
        ScanDepth:      config.GetString("scan-depth"),
        IgnoreUnfixed:  config.GetBool("ignore-unfixed"),
    }
    
    // Initialize Trivy client
    var err error
    p.client, err = trivy.NewClient(trivy.Config{
        // ... Trivy-specific configuration
    })
    return err
}

func (p *TrivyScannerPlugin) Start() error {
    // Subscribe to artifact.pushed events
    go p.processEvents()
    return nil
}

func (p *TrivyScannerPlugin) processEvents() {
    // Implementation omitted for brevity
}

func (p *TrivyScannerPlugin) Stop() error {
    // Close Trivy client, stop event processing
    return nil
}

func (p *TrivyScannerPlugin) Destroy() error {
    // Clean up resources
    return nil
}
```

### Event-Based Plugins
Plugins that primarily consume and produce events can use the event SDK:

```go
type EventPlugin struct {
    BasePlugin
    eventConsumer events.Consumer
    eventPublisher events.Publisher
}

func (p *EventPlugin) Init(config *Config) error {
    // Initialize event consumer and publisher
    var err error
    p.eventConsumer, err = events.NewConsumer(config.GetString("event-stream"))
    if err != nil { return err }
    
    p.eventPublisher, err = events.NewPublisher()
    if err != nil { return err }
    
    return nil
}

func (p *EventPlugin) Start() error {
    // Start consuming events
    go p.consumeEvents()
    return nil
}

func (p *EventPlugin) consumeEvents() {
    for event := range p.eventConsumer.Events() {
        switch event.Type {
        case "registry.artifact.pushed":
            p.handleArtifactPushed(event)
        // ... other event types
        }
    }
}

func (p *EventPlugin) handleArtifactPushed(event *events.Event) {
    // Process the event and potentially publish new events
    // Example: publish trustscore.calculated event
}
```

## Plugin Configuration

### Configuration Sources
Plugins can obtain configuration from:
1. **Plugin Manifest**: Default values defined in the manifest
2. **Kyros Configuration**: Values from the main Kyros configuration file (under `plugins.<plugin-id>`)
3. **Environment Variables**: Override values from environment variables (e.g., `PLUGIN_<ID>_CONFIG_SCAN_DEPTH=full`)
4. **Plugin Configuration File**: A separate configuration file for the plugin (e.g., `plugins/<id>/config.yaml`)

### Configuration Precedence
1. Environment variables (highest precedence)
2. Plugin configuration file
3. Kyros configuration
4. Plugin manifest defaults (lowest precedence)

### Configuration Validation
Plugins should validate their configuration during the `Init` method and return an error if invalid.

## Plugin Security

### Authentication
Plugins authenticate to Kyros services using:
- **Service Accounts**: Each plugin gets a unique service account with minimal required permissions
- **JWT Tokens**: Plugins obtain JWT tokens from the Auth Service using client credentials grant
- **mTLS**: Optional mutual TLS for plugin-core communication

### Authorization
Plugins are subject to the same authorization checks as core services:
- **RBAC**: Role-based access control for Kyros resources
- **PBAC**: Policy-based access control via OPA for complex decisions
- **Resource Scoping**: Plugins can only access resources they are explicitly permitted to access

### Sandboxing
Plugins run in isolated environments:
- **Filesystem**: Read-only access to plugin directory, write access to temporary directory
- **Network**: Configurable outbound connections (allow-list or deny-list)
- **System Resources**: CPU and memory limits enforced
- **Capabilities**: Linux capabilities dropped to minimum required set

## Plugin Distribution

### Plugin Registry
Kyros maintains a plugin registry (similar to Docker Hub) where plugins can be published and discovered.

#### Plugin Metadata
Each plugin in the registry includes:
- Metadata from the plugin manifest
- Download checksums (SHA256)
- Supported Kyros versions
- Dependencies (with version constraints)
- Documentation and examples
- Issue tracker and source repository links

### Installation Methods
1. **Manual**: Download plugin binary and place in plugins directory
2. **Command Line**: `kyros plugin install <plugin-id>[@version]`
3. **Configuration**: Specify plugins to load in Kyros configuration
4. **Orchestration**: Automated deployment via Kubernetes Operator or Helm chart

### Plugin Updates
- **Version Checking**: Platform checks for plugin updates on startup (configurable)
- **Automatic Updates**: Optional automatic update of plugins to latest compatible version
- **Manual Updates**: Administrators can trigger plugin updates via CLI or API
- **Rollback**: Ability to revert to previous plugin version

## Plugin Examples

### Example 1: Custom Storage Plugin
Adds support for a new blob storage backend (e.g., Ceph RADOS).

```go
type CephStoragePlugin struct {
    BasePlugin
    client *rados.Client
    pool   string
}

func (p *CephStoragePlugin) Init(config *Config) error {
    // Parse configuration
    // Initialize Ceph client
    // Connect to cluster
    // Open pool
    return nil
}

func (p *CephStoragePlugin) Start() error {
    // Register storage driver with registry service
    registry.RegisterStorageDriver("ceph", p)
    return nil
}

// Implement storage.StorageDriver interface
func (p *CephStoragePlugin) Get(ctx context.Context, key string) ([]byte, error) {
    // Read object from Ceph pool
    return p.client.Get(p.pool, key)
}

func (p *CephStoragePlugin) Put(ctx context.Context, key string, value []byte) error {
    // Write object to Ceph pool
    return p.client.Put(p.pool, key, value)
}

// ... other StorageDriver methods
```

### Example 2: Notification Plugin
Adds a new notification channel (e.g., Microsoft Teams).

```go
type TeamsNotificationPlugin struct {
    BasePlugin
    client *teams.Client
    config *TeamsConfig
}

func (p *TeamsNotificationPlugin) Init(config *Config) error {
    // Parse configuration
    // Initialize Teams client with webhook URL
    return nil
}

func (p *TeamsNotificationPlugin) Start() error {
    // Subscribe to webhook.delivery.* events
    go p.processEvents()
    return nil
}

func (p *TeamsNotificationPlugin) handleWebhookDeliveryFailed(event *events.Event) {
    // Send notification to Teams channel
    message := teams.Message{
        Text: fmt.Sprintf("Webhook delivery failed: %s", event.Data["webhook"].(string)),
    }
    p.client.SendMessage(message)
    return nil
}

// ... other event handlers
```

### Example 3: Trust Engine Plugin
Adds a new analysis module to the trust score engine (e.g., license compliance scanner).

```go
type LicenseCompliancePlugin struct {
    BasePlugin
    licenseEngine *license.Engine
}

func (p *LicenseCompliancePlugin) Init(config *Config) error {
    // Parse configuration (approved licenses, restricted licenses)
    // Initialize license engine
    return nil
}

func (p *TrustScorePlugin) Start() error {
    // Register analysis module with trust score service
    trustscore.RegisterAnalyzer("license-compliance", p)
    return nil
}

// Implement analyzer.Analyzer interface
func (p *LicenseCompliancePlugin) Analyze(ctx context.Context, artifact *trustscore.Artifact) (interface{}, error) {
    // Scan artifact for license compliance
    return p.licenseEngine.Scan(artifact), nil
}

func (p *LicenseCompliancePlugin) Score(result interface{}) (float64, error) {
    // Convert analysis result to score (0.0-1.0)
    licenseScore := result.(*license.Result).ComplianceScore()
    return licenseScore, nil
}

func (p *LicenseCompliancePlugin) Weight() float64 {
    return 0.10 // 10% weight
}

func (p *LicenseCompliancePlugin) Description() string {
    return "License compliance analyzer"
}

func (p *LicenseCompliancePlugin) IsRequired() bool {
    return false // Optional analyzer
}
```

## Plugin Management API

Kyros provides an API for managing plugins:

### List Plugins
```http
GET /api/v1/plugins
```
Returns a list of installed plugins with their status (loaded, active, error, etc.)

### Get Plugin Details
```http
GET /api/v1/plugins/{plugin-id}
```
Returns detailed information about a specific plugin.

### Load Plugin
```http
POST /api/v1/plugins/{plugin-id}/load
```
Loads and initializes the specified plugin.

### Unload Plugin
```http
POST /api/v1/plugins/{plugin-id}/unload
```
Stops and unloads the specified plugin.

### Enable Plugin
```http
POST /api/v1/plugins/{plugin-id}/enable
```
Enables a loaded plugin (starts it if not already started).

### Disable Plugin
```http
POST /api/v1/plugins/{plugin-id}/disable
```
Disables a loaded plugin (stops it but keeps it loaded).

### Reload Plugin
```http
POST /api/v1/plugins/{plugin-id}/reload
```
Reloads the plugin (unload, then load).

### Get Plugin Configuration
```http
GET /api/v1/plugins/{plugin-id}/config
```
Returns the current configuration of the plugin.

### Update Plugin Configuration
```http
PATCH /api/v1/plugins/{plugin-id}/config
```
Updates the plugin's configuration.

## Plugin Observability

### Metrics
Plugins should export Prometheus metrics for monitoring:
- **Plugin-specific metrics**: Prefixed with `kyros_plugin_<plugin-id>_`
- **Standard metrics**: Plugin uptime, error rates, processing latency, etc.
- **Business metrics**: Specific to plugin function (e.g., scans per hour for scanner plugins)

### Logging
Plugins should use structured logging with trace IDs:
- **Logger**: Provided by the Kyros logging SDK
- **Fields**: Include plugin ID, version, and context
- **Levels**: Standard log levels (debug, info, warn, error, fatal)
- **Sampling**: Configurable sampling rates

### Tracing
Plugins should participate in distributed tracing:
- **Instrumentation**: Use OpenTelemetry SDK
- **Context Propagation**: Extract trace context from incoming events/requests
- **Span Creation**: Create spans for significant operations
- **Attributes**: Add relevant attributes for filtering and analysis

### Health Checks
Plugins should implement health check endpoints:
- **Liveness**: `/healthz` - plugin is running
- **Readiness**: `/readyz` - plugin is ready to serve traffic
- **Startup**: `/startupz` - plugin has completed initialization

## Plugin Testing

### Unit Testing
- Test plugin initialization with various configurations
- Test individual plugin functions in isolation
- Mock dependencies (external services, Kyros APIs)

### Integration Testing
- Test plugin interaction with Kyros services
- Use test doubles for external dependencies
- Verify event production and consumption
- Verify API endpoint behavior (if applicable)

### End-to-End Testing
- Test plugin in a full Kyros deployment
- Verify end-to-end workflows
- Test plugin updates and rollbacks
- Test plugin failure scenarios and recovery

## Plugin Distribution and Marketplace

### Official Plugins
Kyros maintains a set of official plugins that are:
- Fully tested and supported
- Distributed with Kyros releases
- Available in the official plugin registry

### Community Plugins
The community can contribute plugins to the Kyros Plugin Marketplace:
- **Submission Process**: Submit plugin to marketplace for review
- **Review Process**: Check for security, quality, and compatibility
- **Publication**: Approved plugins are published to the marketplace
- **Versioning**: Plugin authors manage their own versions
- **Deprecation**: Authors can deprecate old versions

### Enterprise Plugins
Commercial vendors can distribute proprietary plugins:
- **Verification**: Additional security and IP checks
- **Support**: Optional support agreements
- **Distribution**: Through private registries or direct download

## Conclusion
The Kyros Plugin SDK provides a robust and extensible mechanism for enhancing the platform's functionality. By adhering to the guidelines and interfaces outlined here, developers can create plugins that are secure, maintainable, and seamlessly integrated with the Kyros platform.

The plugin system enables Kyros to adapt to diverse environments and requirements while maintaining a stable core. Through careful design of isolation, versioning, and lifecycle management, the platform can safely evolve with the ecosystem of plugins that surround it.

Regular review and updates to the Plugin SDK are recommended to ensure it continues to meet the needs of plugin developers and platform administrators.