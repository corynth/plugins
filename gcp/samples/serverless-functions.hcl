workflow "gcp-serverless-functions" {
  description = "Deploy and manage Cloud Functions for serverless processing"
  version = "1.0.0"

  step "download_source_code" {
    plugin = "gcp"
    action = "storage_download"
    
    params = {
      bucket = "source-code"
      object_name = "functions/data-processor.zip"
      file_path = "./functions/processor.zip"
      project = "my-project"
    }
  }

  step "deploy_http_function" {
    plugin = "gcp"
    action = "functions_deploy"
    depends_on = ["download_source_code"]
    
    params = {
      name = "api-endpoint"
      source = "./functions/api"
      entry_point = "handleRequest"
      runtime = "nodejs18"
      trigger = "http"
      region = "us-central1"
      memory = "512MB"
      timeout = "60s"
      env_vars = {
        DATABASE_URL = "postgresql://..."
        API_VERSION = "v2"
        LOG_LEVEL = "info"
      }
      project = "my-project"
    }
  }

  step "test_function" {
    plugin = "gcp"
    action = "functions_invoke"
    depends_on = ["deploy_http_function"]
    
    params = {
      name = "api-endpoint"
      data = {
        action = "process"
        payload = {
          user_id = "12345"
          request_type = "data_export"
        }
      }
      region = "us-central1"
      project = "my-project"
    }
  }

  step "deploy_pubsub_function" {
    plugin = "gcp"
    action = "functions_deploy"
    depends_on = ["test_function"]
    
    params = {
      name = "event-processor"
      source = "./functions/events"
      entry_point = "processEvent"
      runtime = "python39"
      trigger = "pubsub"
      topic = "event-stream"
      region = "us-central1"
      memory = "1GB"
      timeout = "300s"
      project = "my-project"
    }
  }

  step "list_input_data" {
    plugin = "gcp"
    action = "storage_list"
    depends_on = ["deploy_pubsub_function"]
    
    params = {
      bucket = "data-input"
      prefix = "batch/"
      project = "my-project"
    }
  }

  step "process_batch" {
    plugin = "gcp"
    action = "functions_invoke"
    depends_on = ["list_input_data"]
    
    params = {
      name = "event-processor"
      data = {
        batch_id = "batch-$(date +%Y%m%d)"
        input_count = "${list_input_data.count}"
        processing_mode = "parallel"
      }
      region = "us-central1"
      project = "my-project"
    }
  }
}