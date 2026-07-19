# Kyros Developer Platform

## Overview
This document describes the developer-facing aspects of the Kyros platform, including the Command Line Interface (CLI), Software Development Kits (SDKs), GitHub Actions, Kubernetes Custom Resource Definitions (CRDs), Terraform Provider, API versioning, and the extension model. It aims to provide developers with the tools and resources needed to build, deploy, and extend applications on Kyros.

## Command Line Interface (CLI)

### Purpose
The Kyros CLI provides a command-line interface for interacting with the Kyros platform, enabling automation, scripting, and efficient management of resources.

### Features
- **Authentication**: Login, logout, token management.
- **Resource Management**: Create, read, update, delete repositories, namespaces, users, etc.
- **Image Operations**: Push, pull, scan, tag images.
- **Trust Score**: View and manage trust scores.
- **Policy Management**: Create, evaluate, delete policies.
- **Webhook Management**: Create, test, delete webhooks.
- **Configuration**: View and update system settings.
- **Debugging**: Logs, metrics, and troubleshooting commands.

### Installation
```bash
# Download the latest release
curl -L https://github.com/kyros-project/kyros-cli/releases/latest/download/kyros-linux-amd64 -o kyros
chmod +x kyros
sudo mv kyros /usr/local/bin/

# Verify installation
kyros version
```

### Usage Examples
```bash
# Login
kyros login --server https://kyros.example.com --username alice --password secret

# Create a namespace
kyros namespace create --name production --description "Production environment"

# Create a repository
kyros repository create --name web-app --namespace production --description "Web application images"

# Push an image (using Docker CLI)
docker tag my-app:latest kyros.example.com/production/web-app:latest
docker push kyros.example.com/production/web-app:latest

# Check trust score
kyros trust-score get --repository production/web-app --tag latest

# Create a webhook
kyros webhook create --name "CI Trigger" --url https://ci.example.com/webhook --events artifact.pushed,trustscore.updated

# Logout
kyros logout
```

### Extensibility
The CLI is designed to be extensible via plugins. Developers can add new commands by implementing the `Command` interface and registering them with the CLI framework.

## Software Development Kits (SDKs)

### Official SDKs
Kyros provides official SDKs for popular programming languages to facilitate integration with the platform.

#### Go SDK
- **Package**: `github.com/kyros-project/kyros-go-sdk`
- **Features**: Full API coverage, authentication helpers, resource management.
- **Example**:
  ```go
  package main

  import (
      "context"
      "fmt"
      "log"

      kyros "github.com/kyros-project/kyros-go-sdk"
  )

  func main() {
      // Create a client
      client, err := kyros.NewClient("https://kyros.example.com", kyros.WithToken("your-token"))
      if err != nil {
          log.Fatal(err)
      }

      // Create a repository
      repo, err := client.Repositories.Create(context.Background(), kyros.Repository{
          Name:        "my-app",
          NamespaceID: "namespace-uuid",
          Description: "My application",
      })
      if err != nil {
          log.Fatal(err)
      }
      fmt.Printf("Created repository: %s\n", repo.ID)
  }
  ```

#### Python SDK
- **Package**: `kyros-py-sdk` (available on PyPI)
- **Features**: Async and sync clients, comprehensive resource management.
- **Example**:
  ```python
  from kyros_sdk import KyrosClient

  client = KyrosClient("https://kyros.example.com", token="your-token")

  # Create a repository
  repo = client.repositories.create(
      name="my-app",
      namespace_id="namespace-uuid",
      description="My application"
  )
  print(f"Created repository: {repo.id}")
  ```

#### JavaScript/TypeScript SDK
- **Package**: `@kyros/client` (available on npm)
- **Features**: TypeScript definitions, browser and Node.js support.
- **Example**:
  ```typescript
  import { KyrosClient } from '@kyros/client';

  const client = new KyrosClient('https://kyros.example.com', { token: 'your-token' });

  // Create a repository
  const repo = await client.repositories.create({
    name: 'my-app',
    namespaceId: 'namespace-uuid',
    description: 'My application'
  });

  console.log(`Created repository: ${repo.id}`);
  ```

### Community SDKs
Community-maintained SDKs are available for other languages (e available for other languages (Java, .NET, Ruby, etc.) and can be found in the Kyros ecosystem.

## GitHub Actions

Kyros provides official GitHub Actions to integrate Kyros into CI/CD workflows.

### Available Actions
1. **kyros-login**: Authenticate to a Kyros registry.
2. **kyros-push**: Push a Docker image to Kyros.
3. **kyros-scan**: Scan an image for vulnerabilities using Kyros trust score.
4. **kyros-deploy**: Deploy a Kubernetes manifest using Kyros (if integrated with GitOps).
5. **kyros-webhook-trigger**: Trigger a Kyros webhook.

### Example Workflow
```yaml
name: CI/CD Pipeline

on:
  push:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3

    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v2

    - name: Build and push Docker image
      uses: kyros-actions/kyros-push@v1
      with:
        kyros-server: ${{ secrets.KYROS_SERVER }}
        kyros-username: ${{ secrets.KYROS_USERNAME }}
        kyros-password: ${{ secrets.KYROS_PASSWORD }}
        image-name: my-app
        image-tag: ${{ github.sha }}

    - name: Scan image
      uses: kyros-actions/kyros-scan@v1
      with:
        kyros-server: ${{ secrets.KYROS_SERVER }}
        kyros-username: ${{ secrets.KYROS_USERNAME }}
        kyros-password: ${{ secrets.KYROS_PASSWORD }}
        image: my-app:${{ github.sha }}

    - name: Deploy to Kubernetes
      uses: kyros-actions/kyros-deploy@v1
      with:
        kubeconfig: ${{ secrets.KUBECONFIG }}
        manifest: ./deploy.yaml
        image: my-app:${{ github.sha }}
```

## Kubernetes Custom Resource Definitions (CRDs)

Kyros provides a Kubernetes Operator that manages Kyros resources via custom resources.

### Custom Resources
1. **KyrosCluster**: Represents a Kyros platform deployment.
2. **KyrosRepository**: Represents a repository in Kyros.
3. **KyrosNamespace**: Represents a namespace in Kyros.
4. **KyrosTrustPolicy**: Represents a trust score policy.
5. **KyrosWebhook**: Represents a webhook subscription.
6. **KyrosUser**: Represents a user in Kyros.

### Example: KyrosRepository CRD
```yaml
apiVersion: kyros.example.com/v1
kind: KyrosRepository
metadata:
  name: web-app
  namespace: kyros-system
spec:
  name: web-app
  namespaceRef:
    name: production
  description: "Web application container images"
  visibility: private
  accessControl:
    - role: admin
      subjects:
        - kind: Group
          name: platform-admins
    - role: developer
      subjects:
        - kind: Group
          name: dev-team
```

### Installation
The Kyros Operator can be installed via Helm:
```bash
helm repo add kyros https://charts.kyros.example.com
helm install kyros-operator kyros/kyros-operator
```

## Terraform Provider

The Kyros Terraform provider allows infrastructure-as-code management of Kyros resources.

### Usage
```hcl
terraform {
  required_providers {
    kyros = {
      source  = "kyros-project/kyros"
      version = "~> 0.1"
    }
  }
}

provider "kyros" {
  server   = "https://kyros.example.com"
  username = var.kyros_username
  password = var.kyros_password
}

resource "kyros_namespace" "production" {
  name        = "production"
  description = "Production environment"
}

resource "kyros_repository" "web_app" {
  name         = "web-app"
  namespace_id = kyros_namespace.production.id
  description  = "Web application images"
  visibility   = "private"
}
```

### Resources
- `kyros_namespace`
- `kyros_repository`
- `kyros_user`
- `kyros_group`
- `kyros_role`
- `kyros_policy`
- `kyros_webhook`
- `kyros_trust_score` (data source)

## API Versioning

Kyros follows a strict API versioning policy to ensure backward compatibility.

### Versioning Scheme
- **URL Path**: `/api/v1/`, `/api/v2/`, etc.
- **Backward Compatibility**: 
  - Within a major version (v1), all changes are backward compatible.
  - Breaking changes only occur in major version bumps.
- **Deprecation**: 
  - Deprecated endpoints are marked with a `Warning` header and a sunset date.
  - Deprecated endpoints are removed after a minimum of 6 months.

### Example
```http
GET /api/v1/repositories
# Returns list of repositories in v1 format

GET /api/v2/repositories
# Returns list of repositories in v2 format (if available)
```

### Response Headers
- `Deprecation`: Date when the endpoint will be deprecated.
- `Sunset`: Date when the endpoint will be removed.
- `Link`: Link to the newer version (if applicable).
- `Warning`: 299 - "Deprecated endpoint; use /api/v2/ instead"

## Extension Model

Kyros is designed to be extensible through multiple mechanisms:

### 1. Plugins (see PLUGIN_SDK.md)
- Add new functionality via dynamically loaded plugins.
- Supported plugin types: scanners, storage, notifications, authentication, trust engine, policy engine, analytics, UI, CLI.

### 2. Webhooks
- Subscribe to events and receive HTTP callbacks.
- Enable integration with external systems (CI/CD, ticketing, chatops).

### 3. API
- RESTful API for programmatic access to all platform features.
- Webhooks for event-driven integrations.

### 4. Kubernetes Operator
- Manage Kyros resources via custom resources.
- Enable GitOps workflows for Kyros configuration.

### 5. Terraform Provider
- Manage Kyros resources via infrastructure-as-code.

### 6. CLI Plugins
- Extend the Kyros CLI with new commands.

### 7. UI Plugins (Future)
- Extend the web interface with new pages, widgets, and components (planned).

## Best Practices for Developers

### 1. Use the SDKs
- Prefer using the official SDKs over direct API calls for better abstraction and error handling.

### 2. Handle Errors Gracefully
- Check for error conditions and handle them appropriately.
- Use exponential backoff for retries on transient errors.

### 3. Respect Rate Limits
- Adhere to the API rate limits (1000 requests/hour for authenticated users).
- Implement retry-after handling for 429 responses.

### 4. Secure Credentials
- Never hardcode credentials in source code.
- Use environment variables, secret managers, or secure credential stores.

### 5. Follow Semantic Versioning
- When building integrations, specify compatible Kyros versions in your dependencies.

### 6. Test Against a Staging Environment
- Always test integrations against a non-production Kyros instance first.

### 7. Monitor Your Integration
- Log errors and metrics from your integration.
- Set up alerts for failure rates or latency spikes.

## Resources
- **Documentation**: https://docs.kyros.example.com
- **API Reference**: https://api.kyros.example.com/v1/docs
- **GitHub Repository**: https://github.com/kyros-project/kyros
- **Community Forum**: https://forum.kyros.example.com
- **Issue Tracker**: https://github.com/kyros-project/kyros/issues

## Conclusion
The Kyros Developer Platform provides a comprehensive set of tools and resources for developers to build, deploy, and extend applications on Kyros. By leveraging the CLI, SDKs, GitHub Actions, Kubernetes CRDs, Terraform Provider, and API, developers can integrate Kyros seamlessly into their workflows and create powerful extensions to the platform.