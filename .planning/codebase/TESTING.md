# Testing Strategy

> Mapped: 2026-07-20

## Test Framework
- **Framework**: terratest (Go) — `test/`
- **Config**: `go.mod` / `go.sum` in `test/`
- **Additional**: Terratest with testify for assertions

## Test Structure
| Test Type | Location | Naming | Scope |
|-----------|----------|--------|-------|
| Unit (Terraform) | `test/unit/` | `*_test.go` | Module input validation, locals computation |
| Integration | `test/integration/` | `*_test.go` | Full module deployment to Azure |
| E2E | `test/e2e/` | `*_test.go` | End-to-end environment deployment |
| Contract | N/A | N/A | Not implemented |
| Policy/OPA | `policies/` | `*.rego` | OPA policies for compliance |

## Test Patterns
| Pattern | Example Location |
|---------|------------------|
| Fixtures | `test/fixtures/` — Terraform fixtures for testing |
| Mocking | Terratest `terraform.InitAndApply` with test fixtures |
| Assertions | `testify/assert`, `testify/require` — `assert.Equal`, `require.NoError` |
| Parallel execution | `t.Parallel()` in test functions |
| Test fixtures | `test/fixtures/<module>/main.tf` — minimal module usage |
| Terratest options | `terraform.Options` with `Vars`, `EnvVars`, `RetryableTerraformErrors` |
| Retry logic | `retry.Do` with `terraform.Apply` retries for eventual consistency |
| Unique names | `fmt.Sprintf("test-%s", random.UniqueId())` for unique resource names |

## Test File Organization
```
test/
├── unit/
│   ├── network_test.go
│   ├── compute_test.go
│   └── storage_test.go
├── integration/
│   ├── network_integration_test.go
│   └── bootstrap_integration_test.go
├── e2e/
│   └── infra_e2e_test.go
├── fixtures/
│   ├── network/
│   │   └── main.tf
│   ├── compute/
│   │   └── main.tf
│   └── bootstrap/
│       └── main.tf
├── go.mod
├── go.sum
└── testutil/
    └── helpers.go
```

## Test Patterns in Code

### Unit Test Pattern (test/unit/network_test.go)
```go
func TestNetworkModule(t *testing.T) {
    t.Parallel()
    
    // Test variable validation
    // Test locals computation
    // Test output values
}
```

### Integration Test Pattern (test/integration/network_integration_test.go)
```go
func TestNetworkModuleIntegration(t *testing.T) {
    t.Parallel()
    
    opts := &terraform.Options{
        TerraformDir: "../fixtures/network",
        Vars: map[string]interface{}{
            "location": "eastus",
            "environment": "test",
        },
        EnvVars: map[string]string{
            "ARM_SUBSCRIPTION_ID": os.Getenv("ARM_SUBSCRIPTION_ID"),
        },
        RetryableTerraformErrors: map[string]bool{
            "ResourceGroupNotFound": true,
        },
    }
    
    defer terraform.Destroy(t, opts)
    terraform.InitAndApply(t, opts)
    
    // Assert outputs
    vnetID := terraform.Output(t, opts, "vnet_id")
    assert.NotEmpty(t, vnetID)
}
```

### Terratest Helpers (test/testutil/helpers.go)
```go
func GetTestOptions(t *testing.T, terraformDir string, vars map[string]interface{}) *terraform.Options
func UniqueResourceName(prefix string) string
func GetAzureSubscriptionID(t *testing.T) string
```

## Coverage
- **Target**: Not explicitly enforced (Terraform/terratest doesn't have traditional coverage)
- **Tool**: `go test -cover` for Go test coverage
- **Config**: `go test -coverprofile=coverage.out ./test/...`

## CI Integration
| Stage | Command | Config Location |
|-------|---------|-----------------|
| Lint | `terraform fmt -check -recursive` | `.github/workflows/ci.yaml` |
| Lint | `terraform validate` | `.github/workflows/ci.yaml` |
| Lint | `tflint` | `.github/workflows/ci.yaml` |
| Lint | `golangci-lint run ./test/...` | `.github/workflows/ci.yaml` |
| Unit | `go test ./test/unit/... -v -parallel 4` | `.github/workflows/ci.yaml` |
| Integration | `go test ./test/integration/... -v -timeout 30m` | `.github/workflows/ci.yaml` |
| E2E | `go test ./test/e2e/... -v -timeout 60m` | `.github/workflows/ci.yaml` |
| Policy | `opa eval -i tfplan.json -d policies/ "data.terraform.deny"` | `.github/workflows/ci.yaml` |

## Test Data Management
| Strategy | Location |
|----------|----------|
| Fixtures | `test/fixtures/<module>/main.tf` — minimal module usage for testing |
| Test accounts | Azure subscription via `ARM_SUBSCRIPTION_ID` env var (GitHub secrets) |
| Cleanup | `defer terraform.Destroy(t, opts)` in every integration test |
| Unique names | `testutil.UniqueResourceName("prefix")` for resource naming |
| Retry logic | `RetryableTerraformErrors` for transient Azure errors |

## Test Dependencies (test/go.mod)
```go
module github.com/<org>/homelab/test

require (
    github.com/gruntwork-io/terratest/modules/terraform v0.58.0
    github.com/stretchr/testify v1.9.0
    github.com/gruntwork-io/terratest/modules/random v0.58.0
    github.com/gruntwork-io/terratest/modules/retry v0.58.0
)
```

## Policy Testing (OPA)
- **Policies**: `policies/*.rego` — Rego policies for Terraform plan validation
- **Test**: `opa test policies/`
- **CI**: `opa eval -i tfplan.json -d policies/ "data.terraform.deny[x]"`

## Test Naming Conventions
| Test Type | Pattern |
|-----------|---------|
| Unit test function | `Test<ModuleName>_<Behavior>` |
| Integration test function | `Test<ModuleName>Integration_<Scenario>` |
| E2E test function | `TestE2E_<Environment>_<Scenario>` |
| Benchmark | `Benchmark<ModuleName>_<Operation>` |

## Running Tests Locally
```bash
# Unit tests
cd test && go test ./unit/... -v -parallel 4

# Integration tests (requires Azure credentials)
cd test && go test ./integration/... -v -timeout 30m

# E2E tests (requires Azure credentials)
cd test && go test ./e2e/... -v -timeout 60m

# All tests with coverage
cd test && go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Test Environment Requirements
- **Azure Subscription**: Required for integration/E2E tests
- **Service Principal**: With Contributor role on subscription
- **Environment Variables**:
  - `ARM_SUBSCRIPTION_ID`
  - `ARM_TENANT_ID`
  - `ARM_CLIENT_ID`
  - `ARM_CLIENT_SECRET`
- **Go**: >= 1.21
- **Terraform**: >= 1.5.0
- **Azure CLI**: For authentication helper

---

*Testing analysis: 2026-07-20*