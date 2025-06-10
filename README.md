# Golidation

A powerful and flexible validation library for Go applications. Golidation provides a fluent interface for validating data with support for multiple languages (English and Persian).

## Installation

```bash
go get github.com/KarimpourMehrab/golidation
```

## Basic Usage

```go
import "github.com/KarimpourMehrab/golidation/package/validator"

// Create a new validator instance
v := validator.Attribute("email").Is("test@example.com")

// Set language (English or Persian)
v.Lang("en") // or v.Lang("fa")

// Add validation rules
v.Email()

// Check for errors
if err := v.Error(); err != nil {
    // Handle error
}
```

## Available Validation Methods

### Basic Validation Rules

- `Accepted()` - Field must be "yes", "on", "1", "true", or true
- `AcceptedIf(other string, otherValue interface{})` - Field must be accepted when another field has a specific value
- `ActiveURL()` - Field must be a valid URL
- `After(compareDate time.Time)` - Field must be a date after the given date
- `AfterOrEqual(compareDate time.Time)` - Field must be a date after or equal to the given date
- `Alpha()` - Field must contain only alphabetic characters
- `AlphaDash()` - Field must contain only letters, numbers, and dashes
- `AlphaNum()` - Field must contain only letters and numbers
- `Array()` - Field must be an array
- `Before(compareDate time.Time)` - Field must be a date before the given date
- `BeforeOrEqual(compareDate time.Time)` - Field must be a date before or equal to the given date
- `Boolean()` - Field must be a boolean value
- `Confirmed(confirmValue interface{})` - Field must match its confirmation value
- `Date()` - Field must be a valid date
- `DateEquals(compareDate time.Time)` - Field must be a date equal to the given date
- `DateFormat(format string)` - Field must match the given date format
- `Declined()` - Field must be "no", "off", "0", "false", or false
- `DeclinedIf(otherValue interface{})` - Field must be declined when another field has a specific value
- `Different(otherValue interface{})` - Field must be different from another field
- `Digits(digits int)` - Field must be numeric and have an exact length
- `DigitsBetween(min int, max int)` - Field must be numeric and have a length between min and max
- `Dimensions(minWidth, minHeight, maxWidth, maxHeight int)` - Field must be an image with specific dimensions
- `Distinct()` - Field must be an array with unique values
- `Email()` - Field must be a valid email address
- `EndsWith(values []string)` - Field must end with one of the given values
- `Exists(validValues []interface{})` - Field must exist in the given array
- `Filled()` - Field must not be empty
- `Image()` - Field must be an image file
- `In(validValues []interface{})` - Field must be included in the given array
- `Integer()` - Field must be an integer
- `InArray(arr []interface{})` - Field must exist in the given array
- `IP()` - Field must be a valid IP address
- `IPv4()` - Field must be a valid IPv4 address
- `IPv6()` - Field must be a valid IPv6 address
- `JSON()` - Field must be a valid JSON string
- `Mimes(validTypes []string)` - Field must be a file of one of the given MIME types
- `NotIn(invalidValues []interface{})` - Field must not be included in the given array
- `NotRegex(pattern string)` - Field must not match the given regular expression
- `Numeric()` - Field must be numeric
- `Present()` - Field must be present
- `Prohibited()` - Field must be absent
- `ProhibitedIf(otherValue interface{})` - Field must be absent when another field has a specific value
- `Regex(pattern string)` - Field must match the given regular expression
- `Required()` - Field must be present and not empty
- `RequiredIf(otherValue interface{})` - Field must be present and not empty when another field has a specific value
- `RequiredUnless(otherValue interface{}, validValues []interface{})` - Field must be present and not empty unless another field has a specific value
- `Same(otherValue interface{})` - Field must match another field
- `StartsWith(values []string)` - Field must start with one of the given values
- `String()` - Field must be a string
- `Timezone()` - Field must be a valid timezone
- `Unique(existingValues []interface{})` - Field must be unique among the given values
- `URL()` - Field must be a valid URL
- `UUID()` - Field must be a valid UUID

### Password Validation Rules

- `PasswordLetters()` - Password must contain at least one letter
- `PasswordMixed()` - Password must contain at least one uppercase and one lowercase letter
- `PasswordNumbers()` - Password must contain at least one number
- `PasswordSymbols()` - Password must contain at least one symbol
- `PasswordUncompromised(leakedValues []string)` - Password must not have appeared in a data leak

### String Length Validation

- `MaxString(max int)` - String must not be longer than the given length
- `MinString(min int)` - String must not be shorter than the given length
- `MaxNumeric(max int)` - Numeric value must not be greater than the given value

## Language Support

Golidation supports both English and Persian languages. You can set the language using:

```go
// For English
v.Lang("en")
// or
v.EnMsg()

// For Persian
v.Lang("fa")
// or
v.FaMsg()
```

## Error Handling

You can handle validation errors in two ways:

```go
// Get the first error
if err := v.Error(); err != nil {
    // Handle single error
}

// Get all errors
if errors := v.Errors(); len(errors) > 0 {
    // Handle multiple errors
}
```

## Example

```go
package main

import (
    "fmt"
    "github.com/KarimpourMehrab/golidation/package/validator"
)

func main() {
    // Create a validator for email
    v := validator.Attribute("email").Is("test@example.com")
    
    // Set language to English
    v.Lang("en")
    
    // Add validation rules
    v.Required().Email()
    
    // Check for errors
    if err := v.Error(); err != nil {
        fmt.Println(err)
    }
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details. 