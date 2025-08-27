workflow "aws-serverless-processing" {
  description = "Serverless data processing pipeline using Lambda and S3"
  version = "1.0.0"

  step "list_input_data" {
    plugin = "aws"
    action = "s3_list"
    
    params = {
      bucket = "raw-data-input"
      prefix = "batch-processing/pending/"
      region = "us-west-2"
    }
  }

  step "download_processing_config" {
    plugin = "aws"
    action = "s3_download"
    depends_on = ["list_input_data"]
    
    params = {
      bucket = "configuration-bucket"
      key = "processing/data-transform-config.json"
      file_path = "./config/transform-config.json"
      region = "us-west-2"
    }
  }

  step "invoke_data_processor" {
    plugin = "aws"
    action = "lambda_invoke"
    depends_on = ["download_processing_config"]
    
    params = {
      function_name = "data-transformation-processor"
      region = "us-west-2"
      invocation_type = "RequestResponse"
      payload = {
        input_bucket = "raw-data-input"
        input_prefix = "batch-processing/pending/"
        output_bucket = "processed-data-output"
        output_prefix = "batch-processing/completed/"
        config_file = "./config/transform-config.json"
        batch_size = 1000
        parallel_workers = 10
      }
    }
  }

  step "verify_processing_results" {
    plugin = "aws"
    action = "s3_list"
    depends_on = ["invoke_data_processor"]
    
    params = {
      bucket = "processed-data-output"
      prefix = "batch-processing/completed/"
      region = "us-west-2"
    }
  }

  step "upload_processing_report" {
    plugin = "aws"
    action = "s3_upload"
    depends_on = ["verify_processing_results"]
    
    params = {
      file_path = "./reports/processing-summary.json"
      bucket = "processing-reports"
      key = "daily-reports/$(date +%Y/%m/%d)/processing-summary.json"
      content_type = "application/json"
      region = "us-west-2"
      metadata = {
        input_objects = "${list_input_data.count}"
        output_objects = "${verify_processing_results.count}"
        processor_status = "${invoke_data_processor.status_code}"
        processing_date = "$(date -Iseconds)"
      }
    }
  }

  step "notify_completion" {
    plugin = "aws"
    action = "lambda_invoke"
    depends_on = ["upload_processing_report"]
    
    params = {
      function_name = "processing-notification"
      region = "us-west-2"
      payload = {
        report_location = "${upload_processing_report.url}"
        input_count = "${list_input_data.count}"
        output_count = "${verify_processing_results.count}"
        status = "completed"
        notification_channels = ["slack", "email"]
      }
    }
  }
}