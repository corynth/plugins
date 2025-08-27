workflow "repository-sync" {
  description = "Clone and synchronize multiple git repositories"
  version     = "1.0.0"

  variable "base_path" {
    type        = string
    default     = "/tmp/repos"
    description = "Base directory for cloned repositories"
  }

  step "setup_workspace" {
    plugin = "shell"
    action = "exec"
    
    params = {
      command = "mkdir -p ${var.base_path} && cd ${var.base_path}"
    }
  }

  step "clone_primary_repo" {
    plugin = "git"
    action = "clone"
    
    depends_on = ["setup_workspace"]
    
    params = {
      url = "https://github.com/corynth/corynth-dist"
      path = "${var.base_path}/corynth-dist"
      branch = "main"
    }
  }

  step "clone_plugins_repo" {
    plugin = "git"
    action = "clone"
    
    depends_on = ["setup_workspace"]
    
    params = {
      url = "https://github.com/justynroberts/corynth-plugin-sources"
      path = "${var.base_path}/plugin-sources"
      branch = "main"
    }
  }

  step "verify_repos" {
    plugin = "shell"
    action = "exec"
    
    depends_on = ["clone_primary_repo", "clone_plugins_repo"]
    
    params = {
      command = "ls -la ${var.base_path}/ && echo 'Repository sync completed'"
    }
  }

  step "get_repo_info" {
    plugin = "git"
    action = "status"
    
    depends_on = ["verify_repos"]
    
    params = {
      path = "${var.base_path}/corynth-dist"
    }
  }
}