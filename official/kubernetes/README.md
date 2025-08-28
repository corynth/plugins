# Kubernetes Plugin

A Go implementation of the Kubernetes plugin for Corynth workflows, providing comprehensive Kubernetes cluster management and resource operations.

## Features

The plugin supports all major Kubernetes operations via `kubectl` commands:

- **apply**: Apply Kubernetes manifests from files or inline YAML
- **get**: Retrieve Kubernetes resources with filtering and formatting
- **describe**: Get detailed resource descriptions
- **scale**: Scale deployments and replica sets
- **logs**: Fetch pod logs with filtering options
- **exec**: Execute commands in running pods
- **port_forward**: Forward local ports to pods (basic implementation)
- **delete**: Delete Kubernetes resources by name, file, or selector

## Requirements

- Go 1.21+ for compilation
- kubectl CLI tool installed and configured
- Access to a Kubernetes cluster

## Usage

The plugin automatically compiles the Go source when invoked through the bash wrapper script. No manual compilation is needed.

### Examples

See `samples/basic-operations.hcl` for comprehensive usage examples.

## Actions

### apply
Apply Kubernetes manifests to the cluster.
- **manifest**: YAML manifest content (string)
- **file**: Path to manifest file (string)  
- **namespace**: Target namespace (string, optional)
- **dry_run**: Perform dry run only (boolean, default: false)

### get
Retrieve Kubernetes resources with optional filtering.
- **resource**: Resource type like 'pods', 'services' (string, required)
- **name**: Specific resource name (string, optional)
- **namespace**: Target namespace (string, optional)
- **all_namespaces**: Search all namespaces (boolean, default: false)
- **selector**: Label selector (string, optional)
- **output**: Output format (string, default: 'json')

### describe
Get detailed information about specific resources.
- **resource**: Resource type (string, required)
- **name**: Resource name (string, required)
- **namespace**: Target namespace (string, optional)

### scale
Scale deployments or replica sets.
- **resource**: Resource type like 'deployment' (string, required)
- **name**: Resource name (string, required)
- **replicas**: Number of replicas (number, required)
- **namespace**: Target namespace (string, optional)

### logs
Fetch logs from pods.
- **pod**: Pod name (string, required)
- **container**: Container name (string, optional)
- **namespace**: Target namespace (string, optional)
- **tail**: Number of lines to show (number, optional)
- **follow**: Follow log output (boolean, default: false)
- **previous**: Get logs from previous container (boolean, default: false)

### exec
Execute commands in running pods.
- **pod**: Pod name (string, required)
- **container**: Container name (string, optional)
- **command**: Command to execute (string, required)
- **namespace**: Target namespace (string, optional)

### port_forward
Forward local ports to pods.
- **pod**: Pod name (string, required)
- **port_mapping**: Port mapping like '8080:80' (string, required)
- **namespace**: Target namespace (string, optional)

### delete
Delete Kubernetes resources.
- **resource**: Resource type (string, required)
- **name**: Resource name (string, optional)
- **file**: Manifest file to delete (string, optional)
- **namespace**: Target namespace (string, optional)
- **selector**: Label selector (string, optional)
- **force**: Force deletion (boolean, default: false)

## Implementation Notes

- Uses `kubectl` CLI commands via Go's `os/exec` package
- Proper JSON parsing of kubectl output when using JSON format
- Comprehensive error handling and validation
- Maintains compatibility with the original Python plugin interface
- Automatic compilation through bash wrapper script
- Full support for namespace operations and selectors