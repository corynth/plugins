workflow "ai-powered-documentation" {
  description = "Generate technical documentation using AI"
  version = "1.0.0"

  step "analyze_kubernetes_setup" {
    plugin = "kubernetes"
    action = "get"
    
    params = {
      resource = "deployments"
      namespace = "production"
    }
  }

  step "list_available_models" {
    plugin = "llm"
    action = "list_models"
    depends_on = ["analyze_kubernetes_setup"]
    
    params = {
      provider = "ollama"
      base_url = "http://localhost:11434"
    }
  }

  step "generate_architecture_overview" {
    plugin = "llm"
    action = "generate"
    depends_on = ["list_available_models"]
    
    params = {
      prompt = "Based on this Kubernetes deployment information, write a comprehensive architecture overview:\n\nDeployments: ${analyze_kubernetes_setup.count}\nNamespace: production\n\nInclude sections on:\n1. System Architecture\n2. Component Overview\n3. Deployment Strategy\n4. Scaling Considerations"
      provider = "ollama"
      model = "codellama:13b"
      temperature = 0.3
      max_tokens = 2000
      system_prompt = "You are a senior DevOps architect. Write clear, detailed technical documentation that would help both developers and operations teams understand the system."
    }
  }

  step "generate_troubleshooting_guide" {
    plugin = "llm"
    action = "chat"
    depends_on = ["generate_architecture_overview"]
    
    params = {
      messages = [
        {
          role = "user"
          content = "Now create a troubleshooting guide for this Kubernetes setup. Include common issues, diagnostic commands, and resolution steps."
        }
      ]
      provider = "anthropic"
      model = "claude-3-sonnet-20240229"
      temperature = 0.2
      max_tokens = 1500
      api_key = "${vault.read('secret/anthropic').data.api_key}"
    }
  }

  step "create_deployment_runbook" {
    plugin = "llm"
    action = "generate"
    depends_on = ["generate_troubleshooting_guide"]
    
    params = {
      prompt = "Create a step-by-step deployment runbook for this Kubernetes environment. Include pre-deployment checks, deployment procedures, and post-deployment validation."
      provider = "openai"
      model = "gpt-4"
      temperature = 0.1
      max_tokens = 1200
      system_prompt = "You are creating operational procedures. Be precise, actionable, and include safety checks."
      api_key = "${vault.read('secret/openai').data.api_key}"
    }
  }

  step "compare_ai_responses" {
    plugin = "llm"
    action = "generate"
    depends_on = ["create_deployment_runbook"]
    
    params = {
      prompt = "Analyze and compare these three AI-generated documents for consistency and completeness:\n\n1. Architecture Overview (Ollama):\n${generate_architecture_overview.response}\n\n2. Troubleshooting Guide (Claude):\n${generate_troubleshooting_guide.response}\n\n3. Deployment Runbook (GPT-4):\n${create_deployment_runbook.response}\n\nProvide a brief analysis of coverage and suggest improvements."
      provider = "anthropic"
      model = "claude-3-opus-20240229"
      temperature = 0.4
      max_tokens = 800
      system_prompt = "You are a technical documentation reviewer. Focus on consistency, completeness, and practical utility."
      api_key = "${vault.read('secret/anthropic').data.api_key}"
    }
  }

  step "generate_complete_documentation" {
    plugin = "reporting"
    action = "create_report"
    depends_on = ["compare_ai_responses"]
    
    params = {
      title = "Kubernetes Production Environment Documentation"
      format = "markdown"
      output_file = "./docs/k8s-production-guide.md"
      template = "technical"
      metadata = {
        author = "AI Documentation Generator"
        date = "$(date +%Y-%m-%d)"
        version = "1.0"
        generated_by = "Corynth AI Workflow"
      }
      sections = [
        {
          heading = "System Architecture"
          level = 2
          content = "${generate_architecture_overview.response}"
        },
        {
          heading = "Troubleshooting Guide"
          level = 2
          content = "${generate_troubleshooting_guide.response}"
        },
        {
          heading = "Deployment Runbook"
          level = 2
          content = "${create_deployment_runbook.response}"
        },
        {
          heading = "AI Analysis Summary"
          level = 2
          content = "${compare_ai_responses.response}"
        },
        {
          heading = "Generated Metadata"
          level = 2
          table = {
            headers = ["Component", "AI Provider", "Model", "Tokens Used", "Temperature"]
            rows = [
              ["Architecture", "${generate_architecture_overview.provider}", "${generate_architecture_overview.model}", "${generate_architecture_overview.tokens_used}", "0.3"],
              ["Troubleshooting", "${generate_troubleshooting_guide.provider}", "claude-3-sonnet", "${generate_troubleshooting_guide.tokens_used}", "0.2"],
              ["Runbook", "${create_deployment_runbook.provider}", "${create_deployment_runbook.model}", "${create_deployment_runbook.tokens_used}", "0.1"],
              ["Analysis", "${compare_ai_responses.provider}", "${compare_ai_responses.model}", "${compare_ai_responses.tokens_used}", "0.4"]
            ]
          }
        }
      ]
    }
  }

  step "create_embeddings_for_search" {
    plugin = "llm"
    action = "embed"
    depends_on = ["generate_complete_documentation"]
    
    params = {
      text = "${generate_complete_documentation.content}"
      provider = "ollama"
      model = "nomic-embed-text"
      base_url = "http://localhost:11434"
    }
  }

  step "save_embeddings" {
    plugin = "file"
    action = "write"
    depends_on = ["create_embeddings_for_search"]
    
    params = {
      path = "./docs/k8s-production-guide.embeddings.json"
      content = "${json(create_embeddings_for_search)}"
    }
  }

  step "display_results" {
    plugin = "reporting"
    action = "display"
    depends_on = ["save_embeddings"]
    
    params = {
      content = "🤖 AI Documentation Generation Complete!\n\n📄 Generated Documentation:\n- Architecture Overview (${generate_architecture_overview.tokens_used} tokens)\n- Troubleshooting Guide (${generate_troubleshooting_guide.tokens_used} tokens)\n- Deployment Runbook (${create_deployment_runbook.tokens_used} tokens)\n- AI Analysis (${compare_ai_responses.tokens_used} tokens)\n\n📊 Total Tokens Used: ${generate_architecture_overview.tokens_used + generate_troubleshooting_guide.tokens_used + create_deployment_runbook.tokens_used + compare_ai_responses.tokens_used}\n\n💾 Files Created:\n- ${generate_complete_documentation.file_path}\n- ${save_embeddings.path}\n\n🔍 Embeddings: ${create_embeddings_for_search.dimensions} dimensions"
      format = "text"
    }
  }
}