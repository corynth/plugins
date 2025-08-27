workflow "cache-management" {
  description = "Redis cache operations and data management"
  version     = "1.0.0"

  variable "redis_host" {
    type        = string
    default     = "localhost"
    description = "Redis server host"
  }

  variable "redis_port" {
    type        = number
    default     = 6379
    description = "Redis server port"
  }

  variable "app_prefix" {
    type        = string
    default     = "myapp"
    description = "Application prefix for cache keys"
  }

  step "set_user_session" {
    plugin = "redis"
    action = "set"
    
    params = {
      host = var.redis_host
      port = var.redis_port
      key = "${var.app_prefix}:session:user123"
      value = "{\"user_id\": 123, \"username\": \"john_doe\", \"login_time\": \"$(date)\"}"
      ttl = 3600
    }
  }

  step "set_app_config" {
    plugin = "redis"
    action = "set"
    
    depends_on = ["set_user_session"]
    
    params = {
      host = var.redis_host
      port = var.redis_port
      key = "${var.app_prefix}:config:features"
      value = "{\"feature_flags\": {\"new_ui\": true, \"beta_api\": false}, \"version\": \"2.1.0\"}"
      ttl = 7200
    }
  }

  step "cache_api_response" {
    plugin = "redis"
    action = "set"
    
    depends_on = ["set_app_config"]
    
    params = {
      host = var.redis_host
      port = var.redis_port
      key = "${var.app_prefix}:cache:popular_posts"
      value = "[{\"id\": 1, \"title\": \"Post 1\"}, {\"id\": 2, \"title\": \"Post 2\"}]"
      ttl = 900
    }
  }

  step "retrieve_user_session" {
    plugin = "redis"
    action = "get"
    
    depends_on = ["cache_api_response"]
    
    params = {
      host = var.redis_host
      port = var.redis_port
      key = "${var.app_prefix}:session:user123"
    }
  }

  step "verify_cache_data" {
    plugin = "shell"
    action = "exec"
    
    depends_on = ["retrieve_user_session"]
    
    params = {
      command = "echo 'Session data retrieved: ${retrieve_user_session.value}' && echo 'Cache operations completed successfully'"
    }
  }
}