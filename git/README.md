# Git Plugin

## Overview

The Git Plugin provides comprehensive Git version control operations for Corynth workflows, enabling repository management, code synchronization, and automated backup workflows.

## Features

- Repository cloning and initialization
- Branch and tag management
- Commit operations with automated messaging
- Push and pull synchronization
- Status checking and diff operations
- Remote repository management
- Authentication support (SSH, HTTPS, tokens)

## Actions

### clone
Clones a Git repository

**Parameters:**
- `url` (string, required): Repository URL (HTTPS or SSH)
- `path` (string, required): Local path for cloning
- `branch` (string, optional): Specific branch to clone (default: "main")
- `depth` (number, optional): Shallow clone depth
- `auth_token` (string, optional): Authentication token for private repos

**Returns:**
- `path`: Local repository path
- `branch`: Cloned branch
- `commit_hash`: Latest commit hash

### status
Checks repository status

**Parameters:**
- `path` (string, required): Repository path

**Returns:**
- `clean`: Boolean indicating if working tree is clean
- `modified_files`: List of modified files
- `untracked_files`: List of untracked files
- `branch`: Current branch name

### commit
Creates a commit with changes

**Parameters:**
- `path` (string, required): Repository path
- `message` (string, required): Commit message
- `files` (list, optional): Specific files to commit (default: all changes)
- `author_name` (string, optional): Commit author name
- `author_email` (string, optional): Commit author email

**Returns:**
- `commit_hash`: New commit hash
- `files_committed`: Number of files committed
- `message`: Commit message used

### push
Pushes commits to remote repository

**Parameters:**
- `path` (string, required): Repository path
- `remote` (string, optional): Remote name (default: "origin")
- `branch` (string, optional): Branch to push (default: current branch)
- `force` (bool, optional): Force push (default: false)

**Returns:**
- `pushed_commits`: Number of commits pushed
- `remote`: Remote repository pushed to
- `branch`: Branch pushed

### pull
Pulls changes from remote repository

**Parameters:**
- `path` (string, required): Repository path
- `remote` (string, optional): Remote name (default: "origin")
- `branch` (string, optional): Branch to pull (default: current branch)

**Returns:**
- `updated`: Boolean indicating if changes were pulled
- `commits_pulled`: Number of new commits
- `files_changed`: Number of files changed

## Usage Examples

### Basic Repository Operations
```hcl
step "clone_repository" {
  plugin = "git"
  action = "clone"
  params = {
    url    = "https://github.com/company/application.git"
    path   = "/tmp/app-source"
    branch = "main"
  }
}

step "check_status" {
  plugin = "git"
  action = "status"
  depends_on = ["clone_repository"]
  params = {
    path = "/tmp/app-source"
  }
}
```

### Automated Backup Workflow
```hcl
step "backup_configuration" {
  plugin = "file"
  action = "copy"
  params = {
    source      = "/etc/nginx/"
    destination = "/backup-repo/nginx/"
    recursive   = true
  }
}

step "commit_backup" {
  plugin = "git"
  action = "commit"
  depends_on = ["backup_configuration"]
  params = {
    path         = "/backup-repo"
    message      = "Automated backup - ${formatdate('YYYY-MM-DD HH:mm', timestamp())}"
    author_name  = "Corynth Automation"
    author_email = "automation@company.com"
  }
}

step "push_backup" {
  plugin = "git"
  action = "push"
  depends_on = ["commit_backup"]
  params = {
    path   = "/backup-repo"
    remote = "origin"
    branch = "main"
  }
}
```

### Code Deployment Pipeline
```hcl
step "clone_source" {
  plugin = "git"
  action = "clone"
  params = {
    url        = "https://github.com/company/app.git"
    path       = "/tmp/deployment-source"
    branch     = var.deploy_branch
    auth_token = var.github_token
  }
}

step "build_application" {
  plugin = "shell"
  action = "script"
  depends_on = ["clone_source"]
  params = {
    working_dir = "/tmp/deployment-source"
    script = <<-EOF
      #!/bin/bash
      npm install
      npm run build
      npm test
    EOF
  }
}

step "tag_release" {
  plugin = "shell"
  action = "exec"
  depends_on = ["build_application"]
  params = {
    working_dir = "/tmp/deployment-source"
    command     = "git tag -a v${var.version} -m 'Release version ${var.version}'"
  }
}

step "push_tag" {
  plugin = "git"
  action = "push"
  depends_on = ["tag_release"]
  params = {
    path   = "/tmp/deployment-source"
    remote = "origin"
    branch = "v${var.version}"
  }
}
```

### Multi-Repository Synchronization
```hcl
variable "repositories" {
  type = list(object({
    name = string
    url  = string
    path = string
  }))
  default = [
    {
      name = "backend"
      url  = "https://github.com/company/backend.git"
      path = "/sync/backend"
    },
    {
      name = "frontend"
      url  = "https://github.com/company/frontend.git"
      path = "/sync/frontend"
    }
  ]
}

# Clone multiple repositories
step "clone_backend" {
  plugin = "git"
  action = "clone"
  params = {
    url  = "https://github.com/company/backend.git"
    path = "/sync/backend"
  }
}

step "clone_frontend" {
  plugin = "git"
  action = "clone"
  params = {
    url  = "https://github.com/company/frontend.git"
    path = "/sync/frontend"
  }
}

# Pull latest changes
step "update_backend" {
  plugin = "git"
  action = "pull"
  depends_on = ["clone_backend"]
  params = {
    path = "/sync/backend"
  }
}

step "update_frontend" {
  plugin = "git"
  action = "pull"
  depends_on = ["clone_frontend"]
  params = {
    path = "/sync/frontend"
  }
}
```

## Advanced Usage

### Branch Management
```hcl
step "create_feature_branch" {
  plugin = "shell"
  action = "exec"
  params = {
    working_dir = "/repo"
    command     = "git checkout -b feature/${var.feature_name}"
  }
}

step "make_changes" {
  plugin = "file"
  action = "write"
  depends_on = ["create_feature_branch"]
  params = {
    path    = "/repo/feature-config.yaml"
    content = "feature: ${var.feature_name}\nenabled: true"
  }
}

step "commit_feature" {
  plugin = "git"
  action = "commit"
  depends_on = ["make_changes"]
  params = {
    path    = "/repo"
    message = "Add feature: ${var.feature_name}"
    files   = ["feature-config.yaml"]
  }
}
```

### Configuration Management
```hcl
step "clone_config_repo" {
  plugin = "git"
  action = "clone"
  params = {
    url  = "https://github.com/company/config.git"
    path = "/tmp/config"
  }
}

step "update_environment_config" {
  plugin = "file"
  action = "template"
  depends_on = ["clone_config_repo"]
  params = {
    template = "/tmp/config/templates/app.yaml.tpl"
    output   = "/tmp/config/environments/${var.environment}/app.yaml"
    variables = {
      environment    = var.environment
      database_host  = var.db_host
      redis_host     = var.redis_host
      api_version    = var.api_version
    }
  }
}

step "commit_config_update" {
  plugin = "git"
  action = "commit"
  depends_on = ["update_environment_config"]
  params = {
    path    = "/tmp/config"
    message = "Update ${var.environment} configuration for API ${var.api_version}"
    files   = ["environments/${var.environment}/app.yaml"]
  }
}

step "deploy_config" {
  plugin = "git"
  action = "push"
  depends_on = ["commit_config_update"]
  params = {
    path = "/tmp/config"
  }
}
```

### Automated Documentation Updates
```hcl
step "clone_docs_repo" {
  plugin = "git"
  action = "clone"
  params = {
    url  = "https://github.com/company/docs.git"
    path = "/tmp/docs"
  }
}

step "generate_api_docs" {
  plugin = "shell"
  action = "script"
  depends_on = ["clone_docs_repo"]
  params = {
    working_dir = "/tmp/docs"
    script = <<-EOF
      #!/bin/bash
      # Generate API documentation
      swagger-codegen generate \
        -i /app/swagger.yaml \
        -l html2 \
        -o api/
      
      # Update changelog
      echo "## Version ${var.api_version} - $(date)" >> CHANGELOG.md
      echo "${var.release_notes}" >> CHANGELOG.md
      echo "" >> CHANGELOG.md
    EOF
  }
}

step "commit_docs_update" {
  plugin = "git"
  action = "commit"
  depends_on = ["generate_api_docs"]
  params = {
    path         = "/tmp/docs"
    message      = "Update documentation for API version ${var.api_version}"
    author_name  = "API Documentation Bot"
    author_email = "docs@company.com"
  }
}

step "publish_docs" {
  plugin = "git"
  action = "push"
  depends_on = ["commit_docs_update"]
  params = {
    path = "/tmp/docs"
  }
}
```

### Disaster Recovery Repository Sync
```hcl
step "clone_primary_repo" {
  plugin = "git"
  action = "clone"
  params = {
    url  = "https://github.com/company/critical-app.git"
    path = "/tmp/primary"
  }
}

step "pull_latest_changes" {
  plugin = "git"
  action = "pull"
  depends_on = ["clone_primary_repo"]
  params = {
    path = "/tmp/primary"
  }
}

step "clone_backup_repo" {
  plugin = "git"
  action = "clone"
  params = {
    url  = "https://backup-git.company.com/critical-app.git"
    path = "/tmp/backup"
  }
}

step "sync_to_backup" {
  plugin = "shell"
  action = "script"
  depends_on = ["pull_latest_changes", "clone_backup_repo"]
  params = {
    script = <<-EOF
      #!/bin/bash
      cd /tmp/primary
      
      # Get latest commit info
      LATEST_COMMIT=$(git rev-parse HEAD)
      COMMIT_MESSAGE=$(git log -1 --pretty=%B)
      
      # Copy changes to backup repo
      cd /tmp/backup
      git remote add primary /tmp/primary
      git fetch primary
      git merge primary/main
      
      echo "Synced commit: $LATEST_COMMIT"
    EOF
  }
}

step "push_to_backup" {
  plugin = "git"
  action = "push"
  depends_on = ["sync_to_backup"]
  params = {
    path = "/tmp/backup"
  }
}
```

## Authentication

### SSH Key Authentication
```hcl
step "setup_ssh_key" {
  plugin = "shell"
  action = "script"
  params = {
    script = <<-EOF
      #!/bin/bash
      # Setup SSH key for Git operations
      mkdir -p ~/.ssh
      echo "${var.ssh_private_key}" > ~/.ssh/id_rsa
      chmod 600 ~/.ssh/id_rsa
      
      # Add GitHub to known hosts
      ssh-keyscan github.com >> ~/.ssh/known_hosts
    EOF
  }
}

step "clone_private_repo" {
  plugin = "git"
  action = "clone"
  depends_on = ["setup_ssh_key"]
  params = {
    url  = "git@github.com:company/private-repo.git"
    path = "/tmp/private"
  }
}
```

### Token-Based Authentication
```hcl
step "clone_with_token" {
  plugin = "git"
  action = "clone"
  params = {
    url        = "https://github.com/company/private-repo.git"
    path       = "/tmp/private"
    auth_token = var.github_token
  }
}
```

## Error Handling

### Repository State Validation
```hcl
step "check_repo_status" {
  plugin = "git"
  action = "status"
  params = {
    path = "/repo"
  }
}

step "handle_dirty_repo" {
  plugin = "shell"
  action = "exec"
  depends_on = ["check_repo_status"]
  condition  = "!${check_repo_status.clean}"
  params = {
    working_dir = "/repo"
    command     = "git stash save 'Automated stash before deployment'"
  }
}

step "safe_pull" {
  plugin = "git"
  action = "pull"
  depends_on = ["handle_dirty_repo"]
  params = {
    path = "/repo"
  }
}
```

### Merge Conflict Resolution
```hcl
step "attempt_pull" {
  plugin = "git"
  action = "pull"
  params = {
    path = "/repo"
  }
}

step "handle_merge_conflicts" {
  plugin = "shell"
  action = "script"
  depends_on = ["attempt_pull"]
  condition  = "${attempt_pull.error != ''}"
  params = {
    working_dir = "/repo"
    script = <<-EOF
      #!/bin/bash
      # Reset to known good state on merge conflicts
      git reset --hard HEAD
      git clean -fd
      
      # Try pull again
      git pull origin main
    EOF
  }
}
```

## Best Practices

1. **Always check repository status** before making changes
2. **Use meaningful commit messages** with context
3. **Set up proper authentication** for private repositories
4. **Handle merge conflicts gracefully** with automated resolution
5. **Use branches** for feature development and isolation
6. **Tag releases** for version management
7. **Clean up temporary repositories** after operations
8. **Implement backup strategies** for critical repositories

## Security Considerations

1. **Never commit secrets** or sensitive information
2. **Use secure authentication methods** (SSH keys, tokens)
3. **Validate repository URLs** to prevent malicious clones
4. **Set proper file permissions** for SSH keys and tokens
5. **Use read-only tokens** when possible
6. **Implement access controls** for repository operations

## Troubleshooting

### Common Issues

**Authentication Failed**
```
Error: authentication failed
```
- Verify credentials (SSH key, token, username/password)
- Check repository permissions
- Ensure correct authentication method

**Repository Not Found**
```
Error: repository not found
```
- Verify repository URL
- Check repository permissions
- Ensure repository exists

**Merge Conflicts**
```
Error: merge conflict detected
```
- Implement conflict resolution strategy
- Use automated merge tools
- Reset to known good state if needed

**Network Connectivity**
```
Error: failed to connect to repository
```
- Check network connectivity
- Verify DNS resolution
- Check firewall rules

### Debug Mode

Enable detailed logging:
```bash
export CORYNTH_DEBUG=true
export GIT_TRACE=1
corynth run workflow.hcl
```

## Sample Workflows

See the `/samples` directory for complete workflow examples:
- `automated-backup.hcl`: Automated configuration backup to Git
- `repository-sync.hcl`: Multi-repository synchronization workflow