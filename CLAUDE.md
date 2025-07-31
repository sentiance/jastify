# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Jastify is a Go CLI tool that converts Datadog JSON configurations (dashboards and monitors) to Terraform HCL format. It enables infrastructure-as-code management of Datadog resources by parsing JSON exports from the Datadog UI and generating corresponding `datadog_dashboard` and `datadog_monitor` Terraform resources.

## Architecture

The application follows a simple converter architecture:

1. **Main Entry Point** (`main.go`): Handles CLI argument parsing, JSON input (from file or stdin), and delegates to appropriate converter based on JSON structure
2. **Converter Package** (`converter/`): Contains the core conversion logic with separate modules for dashboards and monitors
3. **Type Detection**: Automatically detects resource type - monitors have a `name` field, dashboards do not

### Key Components

- **Dashboard Converter** (`converter/dashboard_converter.go`): Handles conversion of Datadog dashboard JSON to Terraform `datadog_dashboard` resources
- **Monitor Converter** (`converter/monitor-converter.go`): Handles conversion of Datadog monitor JSON to Terraform `datadog_monitor` resources  
- **Utilities** (`converter/utils.go`): Shared conversion utilities including type-safe JSON parsing, Terraform literal formatting, and block generation
- **Definition Maps**: Each converter uses definition maps that specify how JSON fields map to Terraform resource attributes

### Conversion Pattern

The converters use a function map pattern where each JSON field has a corresponding `stringFunc` that defines how to convert that field to Terraform HCL:

```go
var DASHBOARD = map[string]stringFunc{
    "title": stringGen("title"),
    "widgets": func(v any) string { /* complex conversion logic */ },
}
```

## Development Commands

### Building and Running

```bash
# Build the binary
go build -o jastify .

# Run without installing
go run . dashboard.json

# Install globally  
go install github.com/fgm/jastify@latest
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with golden file updates (after making changes)
go test -update

# Run specific test package
go test ./converter/
```

### Golden File Testing

The project uses [goldie v2](https://github.com/sebdah/goldie) for golden file testing. Test data is stored in `converter/testdata/` with corresponding `.golden` files containing expected Terraform output.

When modifying conversion logic:
1. Run tests to see failures
2. Review changes carefully  
3. Run `go test -update` to update golden files
4. Verify only intended changes are included

## Usage Patterns

### Input Methods
- **File input**: `jastify dashboard.json` (resource named from filename)
- **Stdin input**: `jastify` then paste JSON + Ctrl-D (resource named `dashboard_1` or `monitor_1`)

### Output Formatting
Raw output should be formatted with Terraform:
```bash
jastify dashboard.json | terraform fmt > dashboard.tf
```

## Resource Type Detection

The converter automatically detects resource type by checking for the presence of a `name` field:
- **Monitor**: JSON contains `name` field → generates `datadog_monitor` resource
- **Dashboard**: JSON lacks `name` field → generates `datadog_dashboard` resource

## Dependencies

- **Core**: Standard library only for main functionality
- **Testing**: `github.com/sebdah/goldie/v2` for golden file testing
- **Testing**: `github.com/google/go-cmp` for test comparisons

## Error Handling

The conversion process uses a definition-based approach where unknown JSON fields result in clear error messages indicating which field couldn't be converted. This helps identify when Datadog introduces new fields that need mapping support.