# Calculator Plugin

Mathematical calculations and unit conversions for Corynth workflows.

## Actions

### calculate
Perform basic mathematical calculations.

**Parameters:**
- `expression` (string, required): Mathematical expression to evaluate
- `precision` (number, optional): Number of decimal places (default: 2)

**Returns:**
- `result` (number): Calculation result
- `expression` (string): Original expression

**Example:**
```hcl
step "math" {
  plugin = "calculator"
  action = "calculate"
  params = {
    expression = "10 * 5 + 3"
    precision = 2
  }
}
```

### convert
Convert between different units.

**Parameters:**
- `value` (number, required): Value to convert
- `from` (string, required): Source unit
- `to` (string, required): Target unit

**Supported Units:**
- Temperature: `celsius`, `fahrenheit`
- Length: `meters`, `feet`

**Returns:**
- `result` (number): Converted value
- `unit` (string): Target unit

**Example:**
```hcl
step "convert_temp" {
  plugin = "calculator"
  action = "convert"
  params = {
    value = 25
    from = "celsius"
    to = "fahrenheit"
  }
}
```

## Installation

```bash
corynth plugin install calculator
```

The plugin will be compiled from source and installed automatically.