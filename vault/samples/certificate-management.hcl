workflow "certificate-management" {
  description = "Vault PKI certificate management and rotation"
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

  variable "domain_name" {
    type        = string
    default     = "example.com"
    description = "Domain name for certificates"
  }

  variable "service_name" {
    type        = string
    default     = "api"
    description = "Service name for certificate"
  }

  step "enable_pki_engine" {
    plugin = "vault"
    action = "mount"
    
    params = {
      addr  = var.vault_addr
      token = var.vault_token
      path  = "pki"
      type  = "pki"
      config = {
        default_lease_ttl = "768h"
        max_lease_ttl     = "8760h"
      }
    }
  }

  step "configure_ca" {
    plugin = "vault"
    action = "write"
    
    depends_on = ["enable_pki_engine"]
    
    params = {
      addr   = var.vault_addr
      token  = var.vault_token
      path   = "pki/root/generate/internal"
      data = {
        common_name = "Corynth Internal CA"
        ttl         = "8760h"
        key_bits    = "4096"
        country     = "US"
        organization = "Corynth"
      }
    }
  }

  step "configure_ca_urls" {
    plugin = "vault"
    action = "write"
    
    depends_on = ["configure_ca"]
    
    params = {
      addr   = var.vault_addr
      token  = var.vault_token
      path   = "pki/config/urls"
      data = {
        issuing_certificates    = "${var.vault_addr}/v1/pki/ca"
        crl_distribution_points = "${var.vault_addr}/v1/pki/crl"
      }
    }
  }

  step "create_pki_role" {
    plugin = "vault"
    action = "write"
    
    depends_on = ["configure_ca_urls"]
    
    params = {
      addr   = var.vault_addr
      token  = var.vault_token
      path   = "pki/roles/${var.service_name}-role"
      data = {
        allowed_domains    = [var.domain_name, "*.${var.domain_name}"]
        allow_subdomains   = true
        allow_glob_domains = false
        max_ttl           = "720h"
        key_bits          = "2048"
        key_type          = "rsa"
      }
    }
  }

  step "issue_certificate" {
    plugin = "vault"
    action = "write"
    
    depends_on = ["create_pki_role"]
    
    params = {
      addr   = var.vault_addr
      token  = var.vault_token
      path   = "pki/issue/${var.service_name}-role"
      data = {
        common_name = "${var.service_name}.${var.domain_name}"
        alt_names   = "www.${var.service_name}.${var.domain_name},${var.service_name}-staging.${var.domain_name}"
        ttl         = "168h"
      }
    }
  }

  step "save_certificate" {
    plugin = "file"
    action = "write"
    
    depends_on = ["issue_certificate"]
    
    params = {
      path    = "/tmp/${var.service_name}.${var.domain_name}.crt"
      content = "${issue_certificate.certificate}"
    }
  }

  step "save_private_key" {
    plugin = "file"
    action = "write"
    
    depends_on = ["save_certificate"]
    
    params = {
      path    = "/tmp/${var.service_name}.${var.domain_name}.key"
      content = "${issue_certificate.private_key}"
    }
  }

  step "save_ca_certificate" {
    plugin = "file"
    action = "write"
    
    depends_on = ["save_private_key"]
    
    params = {
      path    = "/tmp/ca.crt"
      content = "${issue_certificate.issuing_ca}"
    }
  }

  step "verify_certificate" {
    plugin = "shell"
    action = "exec"
    
    depends_on = ["save_ca_certificate"]
    
    params = {
      command = <<-EOF
        echo "=== Certificate Management Summary ==="
        echo "Service: ${var.service_name}.${var.domain_name}"
        echo "Certificate serial: ${issue_certificate.serial_number}"
        echo "Certificate expiry: ${issue_certificate.expiration}"
        echo ""
        echo "Files created:"
        ls -la /tmp/${var.service_name}.${var.domain_name}.* /tmp/ca.crt
        echo ""
        echo "Certificate details:"
        openssl x509 -in /tmp/${var.service_name}.${var.domain_name}.crt -text -noout | grep -E "(Subject:|DNS:|Not After)"
      EOF
    }
  }

  step "create_renewal_reminder" {
    plugin = "vault"
    action = "write"
    
    depends_on = ["verify_certificate"]
    
    params = {
      addr   = var.vault_addr
      token  = var.vault_token
      path   = "secret/data/certificates/${var.service_name}"
      data = {
        domain          = "${var.service_name}.${var.domain_name}"
        serial_number   = "${issue_certificate.serial_number}"
        issued_date     = "$(date)"
        expiry_date     = "${issue_certificate.expiration}"
        renewal_needed  = "$(date -d '+5 days' -d '${issue_certificate.expiration}')"
        managed_by      = "corynth-vault-plugin"
      }
    }
  }
}