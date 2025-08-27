# JSON Processor Plugin

JSON parsing, manipulation, and validation for Corynth workflows.

## Actions

### parse
Parse a JSON string into structured data.

**Parameters:**
- `json` (string, required): JSON string to parse

**Returns:**
- `data` (object): Parsed JSON data
- `valid` (boolean): Whether the JSON was valid
- `error` (string): Error message if parsing failed

**Example:**
```hcl
step "parse_json" {
  plugin = "json-processor"
  action = "parse"
  params = {
    json = "{\"name\": \"John\", \"age\": 30}"
  }
}
```

### query
Query JSON data using simple path syntax.

**Parameters:**
- `data` (object, required): JSON data to query (or JSON string)
- `path` (string, required): JSON path (e.g., "user.name", "items[0].title")

**Returns:**
- `result` (object): Query result
- `found` (boolean): Whether the path was found

**Example:**
```hcl
step "get_user_name" {
  plugin = "json-processor"
  action = "query"
  params = {
    data = step.parse_json.outputs.data
    path = "user.name"
  }
}
```

### validate
Validate JSON structure and check for required fields.

**Parameters:**
- `json` (string, required): JSON string to validate
- `required_fields` (array, optional): List of required field paths

**Returns:**
- `valid` (boolean): Whether the JSON is valid
- `errors` (array): Validation errors

**Example:**
```hcl
step "validate_user" {
  plugin = "json-processor"
  action = "validate"
  params = {
    json = "{\"name\": \"John\", \"email\": \"john@example.com\"}"
    required_fields = ["name", "email"]
  }
}
```

### transform
Transform JSON data by mapping fields to new structure.

**Parameters:**
- `data` (object, required): JSON data to transform
- `mappings` (object, required): Field mappings (old_path: new_path)

**Returns:**
- `result` (object): Transformed data

**Example:**
```hcl
step "transform_user" {
  plugin = "json-processor"
  action = "transform"
  params = {
    data = step.parse_json.outputs.data
    mappings = {
      "name" = "full_name"
      "age" = "user_age"
      "contact.email" = "email_address"
    }
  }
}
```

## Path Syntax

The query and mapping system supports:
- Simple field access: `name`, `user.name`
- Array indexing: `items[0]`, `users[1].name`
- Nested objects: `contact.address.street`

## Installation

```bash
corynth plugin install json-processor
```

The plugin will be compiled from source and installed automatically.