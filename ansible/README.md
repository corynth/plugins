# Ansible Plugin

Ansible automation plugin for Corynth that provides configuration management, deployments, and infrastructure orchestration capabilities.

## Features

- Execute Ansible playbooks with full parameter support
- Run ad-hoc commands across infrastructure
- Inventory management and host discovery
- Ansible Vault integration for secrets
- Real Ansible CLI integration
- Comprehensive error handling

## Prerequisites

- Ansible installed and available in PATH
- Python with required Ansible modules
- SSH access to target hosts (for remote execution)
- Proper inventory configuration

## Installation

```bash
cd ansible
go build -o corynth-plugin-ansible main.go
```

## Actions

### playbook

Execute an Ansible playbook with full parameter support.

**Parameters:**
- `playbook_path` (required): Path to playbook YAML file
- `inventory` (optional): Inventory file path or host list (default: "localhost,")
- `extra_vars` (optional): Extra variables as key-value object
- `tags` (optional): Comma-separated list of tags to run
- `limit` (optional): Limit execution to specific hosts
- `check_mode` (optional): Run in check mode/dry-run (default: false)

**Returns:**
- `status`: Execution status (success, changed, failed)
- `output`: Full Ansible output
- `changed`: Whether any changes were made
- `failed`: Whether execution failed
- `exit_code`: Command exit code

**Example:**
```hcl
step "deploy_application" {
  plugin = "ansible"
  action = "playbook"
  params = {
    playbook_path = "./playbooks/deploy.yml"
    inventory = "./inventory/production"
    extra_vars = {
      app_version = "1.2.3"
      environment = "production"
      rollback_enabled = true
    }
    tags = "deploy,configure"
    limit = "webservers"
  }
}
```

### adhoc

Run ad-hoc Ansible commands for quick operations.

**Parameters:**
- `hosts` (required): Target hosts pattern
- `module` (required): Ansible module to execute
- `args` (optional): Module arguments
- `inventory` (optional): Inventory file or host list
- `become` (optional): Use privilege escalation (default: false)

**Returns:**
- `status`: Command execution status
- `output`: Command output

**Example:**
```hcl
step "restart_services" {
  plugin = "ansible"
  action = "adhoc"
  params = {
    hosts = "webservers"
    module = "service"
    args = "name=nginx state=restarted"
    inventory = "./inventory/production"
    become = true
  }
}
```

### inventory_list

List and inspect Ansible inventory.

**Parameters:**
- `inventory` (required): Inventory file path
- `group` (optional): Specific group to inspect

**Returns:**
- `hosts`: Array of host names
- `inventory`: Full inventory data structure

**Example:**
```hcl
step "check_inventory" {
  plugin = "ansible"
  action = "inventory_list"
  params = {
    inventory = "./inventory/staging"
    group = "databases"
  }
}
```

### vault_decrypt

Decrypt Ansible Vault files for processing.

**Parameters:**
- `vault_file` (required): Path to encrypted vault file
- `vault_password` (required): Vault password

**Returns:**
- `status`: Decryption status
- `content`: Decrypted file content

**Example:**
```hcl
step "decrypt_secrets" {
  plugin = "ansible"
  action = "vault_decrypt"
  params = {
    vault_file = "./vars/secrets.yml"
    vault_password = "${vault_password}"
  }
}
```

## Sample Workflows

### Infrastructure Deployment
```hcl
workflow "infrastructure_deployment" {
  description = "Complete infrastructure deployment with Ansible"
  
  step "check_connectivity" {
    plugin = "ansible"
    action = "adhoc"
    params = {
      hosts = "all"
      module = "ping"
      inventory = "./inventory/production"
    }
  }
  
  step "deploy_base_config" {
    plugin = "ansible"
    action = "playbook"
    depends_on = ["check_connectivity"]
    params = {
      playbook_path = "./playbooks/base-config.yml"
      inventory = "./inventory/production"
      extra_vars = {
        environment = "production"
        config_version = "2.1.0"
      }
    }
  }
  
  step "deploy_application" {
    plugin = "ansible"
    action = "playbook"
    depends_on = ["deploy_base_config"]
    params = {
      playbook_path = "./playbooks/app-deploy.yml"
      inventory = "./inventory/production"
      extra_vars = {
        app_version = "${app_version}"
        database_host = "${db_host}"
      }
      tags = "deploy,configure,start"
    }
  }
  
  step "verify_deployment" {
    plugin = "ansible"
    action = "adhoc"
    depends_on = ["deploy_application"]
    params = {
      hosts = "webservers"
      module = "uri"
      args = "url=http://localhost/health"
      inventory = "./inventory/production"
    }
  }
}
```

### Configuration Management
```hcl
workflow "configuration_drift_check" {
  description = "Check and remediate configuration drift"
  
  step "check_drift" {
    plugin = "ansible"
    action = "playbook"
    params = {
      playbook_path = "./playbooks/system-config.yml"
      inventory = "./inventory/all"
      check_mode = true
      extra_vars = {
        check_only = true
      }
    }
  }
  
  step "remediate_if_drift" {
    plugin = "ansible"
    action = "playbook"
    depends_on = ["check_drift"]
    params = {
      playbook_path = "./playbooks/system-config.yml"
      inventory = "./inventory/all"
      extra_vars = {
        apply_changes = "${check_drift.changed}"
      }
    }
  }
}
```

### Security Patching
```hcl
workflow "security_patching" {
  description = "Automated security patch deployment"
  
  step "check_updates" {
    plugin = "ansible"
    action = "adhoc"
    params = {
      hosts = "all"
      module = "package"
      args = "name=* state=latest"
      inventory = "./inventory/all"
      check_mode = true
      become = true
    }
  }
  
  step "apply_security_patches" {
    plugin = "ansible"
    action = "playbook"
    depends_on = ["check_updates"]
    params = {
      playbook_path = "./playbooks/security-patches.yml"
      inventory = "./inventory/all"
      extra_vars = {
        patch_type = "security"
        reboot_required = true
      }
      tags = "security,patches"
    }
  }
  
  step "verify_patch_status" {
    plugin = "ansible"
    action = "adhoc"
    depends_on = ["apply_security_patches"]
    params = {
      hosts = "all"
      module = "command"
      args = "cat /var/run/reboot-required"
      inventory = "./inventory/all"
    }
  }
}
```

## Command Line Arguments

The plugin translates parameters to standard Ansible CLI arguments:

- `inventory` → `-i inventory`
- `extra_vars` → `--extra-vars '{"key":"value"}'`
- `tags` → `--tags tag1,tag2`
- `limit` → `--limit pattern`
- `check_mode` → `--check`
- `become` → `--become`

## Output Parsing

The plugin parses Ansible output to provide structured results:
- **Status**: Determines success/changed/failed from output
- **Changes**: Detects if any tasks made changes
- **Exit Codes**: Captures command exit codes for error handling
- **Output**: Provides full Ansible output for debugging

## Error Handling

Comprehensive error handling for:
- Missing Ansible installation
- Invalid playbook paths
- Network connectivity issues
- Authentication failures
- Syntax errors in playbooks
- Host unreachability

## Security Features

- Vault password handling with temporary files
- SSH key authentication support
- Privilege escalation controls
- Inventory validation
- Secure parameter passing

## Dependencies

- Ansible (ansible, ansible-playbook commands)
- Python with Ansible modules
- SSH client for remote operations

## License

MIT License