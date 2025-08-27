workflow "secret-management" {
  description = "HashiCorp Vault secret management and rotation"
  version     = "1.0.0"

  variable "vault_addr" {
    type        = string
    default     = "http://localhost:8200"
    description = "Vault server address"
  }

  variable "vault_token" {
    type        = string
    default     = "dev-token"
    description = "Vault authentication token"
  }

  variable "app_name" {
    type        = string
    default     = "webapp"
    description = "Application name for secret paths"
  }

  variable "environment" {
    type        = string
    default     = "development"
    description = "Environment name"
  }

  step "enable_kv_engine" {
    plugin = "vault"
    action = "mount"
    
    params = {
      addr  = var.vault_addr
      token = var.vault_token
      path  = "secret"
      type  = "kv-v2"
    }
  }

  step "store_database_credentials" {
    plugin = "vault"
    action = "write"
    
    depends_on = ["enable_kv_engine"]
    
    params = {
      addr   = var.vault_addr
      token  = var.vault_token
      path   = "secret/data/${var.app_name}/${var.environment}/database"
      data = {
        username = "db_user"
        password = "supersecret123"
        host     = "db.internal.com"
        port     = "5432"
        database = "${var.app_name}_${var.environment}"
      }
    }
  }

  step "store_api_keys" {
    plugin = "vault"
    action = "write"
    
    depends_on = ["store_database_credentials"]
    
    params = {
      addr   = var.vault_addr
      token  = var.vault_token
      path   = "secret/data/${var.app_name}/${var.environment}/api-keys"
      data = {
        stripe_secret_key    = "sk_test_123456789"
        sendgrid_api_key     = "SG.123456789"
        jwt_secret           = "jwt-super-secret-key"
        encryption_key       = "32-char-encryption-key-here"
      }
    }
  }

  step "store_service_credentials" {
    plugin = "vault"
    action = "write"
    
    depends_on = ["store_api_keys"]
    
    params = {
      addr   = var.vault_addr
      token  = var.vault_token
      path   = "secret/data/${var.app_name}/${var.environment}/services"
      data = {
        redis_password       = "redis-password-123"
        elasticsearch_user   = "elastic"
        elasticsearch_pass   = "elastic-password"
        s3_access_key        = "AKIAIOSFODNN7EXAMPLE"
        s3_secret_key        = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
      }
    }
  }

  step "retrieve_database_config" {
    plugin = "vault"
    action = "read"
    
    depends_on = ["store_service_credentials"]
    
    params = {
      addr   = var.vault_addr
      token  = var.vault_token
      path   = "secret/data/${var.app_name}/${var.environment}/database"
    }
  }

  step "list_secret_paths" {
    plugin = "vault"
    action = "list"
    
    depends_on = ["retrieve_database_config"]
    
    params = {
      addr   = var.vault_addr
      token  = var.vault_token
      path   = "secret/metadata/${var.app_name}/${var.environment}"
    }
  }

  step "create_secret_policy" {
    plugin = "vault"
    action = "policy"
    
    depends_on = ["list_secret_paths"]
    
    params = {
      addr   = var.vault_addr
      token  = var.vault_token
      name   = "${var.app_name}-${var.environment}-policy"
      policy = <<-EOF
        # Allow reading secrets for specific app and environment
        path "secret/data/${var.app_name}/${var.environment}/*" {
          capabilities = ["read"]
        }
        
        # Allow listing secret paths
        path "secret/metadata/${var.app_name}/${var.environment}/*" {
          capabilities = ["list"]
        }
        
        # Allow updating specific secrets
        path "secret/data/${var.app_name}/${var.environment}/api-keys" {
          capabilities = ["create", "update"]
        }
      EOF
    }
  }

  step "verify_secret_access" {
    plugin = "shell"
    action = "exec"
    
    depends_on = ["create_secret_policy"]
    
    params = {
      command = "echo '=== Vault Secret Management Summary ===' && echo 'Application: ${var.app_name}' && echo 'Environment: ${var.environment}' && echo 'Database host: ${retrieve_database_config.data.host}' && echo 'Secrets stored: ${list_secret_paths.count}' && echo 'Policy created: ${var.app_name}-${var.environment}-policy'"
    }
  }
}