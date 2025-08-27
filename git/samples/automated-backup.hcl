workflow "automated-backup" {
  description = "Automated git repository backup and versioning"
  version     = "1.0.0"

  variable "repo_path" {
    type        = string
    default     = "."
    description = "Path to git repository"
  }

  variable "backup_message" {
    type        = string
    default     = "Automated backup via Corynth"
    description = "Commit message for backup"
  }

  step "check_repo_status" {
    plugin = "git"
    action = "status"
    
    params = {
      path = var.repo_path
    }
  }

  step "backup_changes" {
    plugin = "git"
    action = "commit"
    
    depends_on = ["check_repo_status"]
    
    params = {
      path = var.repo_path
      message = "${var.backup_message} - $(date)"
      add_all = true
    }
  }

  step "verify_backup" {
    plugin = "shell"
    action = "exec"
    
    depends_on = ["backup_changes"]
    
    params = {
      command = "echo 'Backup completed. Commit: ${backup_changes.commit}' && echo 'Repository status: ${check_repo_status.clean ? 'Clean' : 'Has changes'}'"
    }
  }
}