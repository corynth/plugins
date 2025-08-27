# Terraform Plugin

Infrastructure as Code operations using Terraform for Corynth workflows.

## Actions

### init
Initialize Terraform working directory.

**Parameters:**
- `working_dir` (string, optional): Directory containing Terraform configuration (default: ".")
- `backend` (boolean, optional): Initialize backend configuration (default: true)

**Returns:**
- `success` (boolean): Whether initialization succeeded
- `output` (string): Command output

**Example:**
```hcl
step "terraform_init" {
  plugin = "terraform"
  action = "init"
  params = {
    working_dir = "./infrastructure"
    backend = true
  }
}
```

### plan
Create a Terraform execution plan.

**Parameters:**
- `working_dir` (string, optional): Directory containing Terraform configuration (default: ".")
- `var_file` (string, optional): Path to variables file
- `out` (string, optional): Path to save plan file

**Returns:**
- `changes` (number): Number of resources to change
- `success` (boolean): Whether planning succeeded
- `output` (string): Plan output

**Example:**
```hcl
step "terraform_plan" {
  plugin = "terraform"
  action = "plan"
  params = {
    working_dir = "./infrastructure"
    var_file = "production.tfvars"
    out = "plan.out"
  }
}
```

### apply
Apply Terraform configuration.

**Parameters:**
- `working_dir` (string, optional): Directory containing Terraform configuration (default: ".")
- `plan_file` (string, optional): Path to plan file
- `auto_approve` (boolean, optional): Skip interactive approval (default: false)

**Returns:**
- `applied` (number): Number of resources applied
- `success` (boolean): Whether apply succeeded
- `output` (string): Apply output

**Example:**
```hcl
step "terraform_apply" {
  plugin = "terraform"
  action = "apply"
  params = {
    working_dir = "./infrastructure"
    plan_file = "plan.out"
    auto_approve = true
  }
}
```

### destroy
Destroy Terraform-managed infrastructure.

**Parameters:**
- `working_dir` (string, optional): Directory containing Terraform configuration (default: ".")
- `auto_approve` (boolean, optional): Skip interactive approval (default: false)

**Returns:**
- `destroyed` (number): Number of resources destroyed
- `success` (boolean): Whether destroy succeeded

**Example:**
```hcl
step "terraform_destroy" {
  plugin = "terraform"
  action = "destroy"
  params = {
    working_dir = "./infrastructure"
    auto_approve = true
  }
}
```

## Installation

```bash
corynth plugin install terraform
```

The plugin will be compiled from source and installed automatically.

## Note

This plugin currently returns mock data for demonstration. In a production environment, it would execute actual Terraform commands.