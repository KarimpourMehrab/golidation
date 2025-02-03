
# Golang Validation Package

This package provides a set of validation methods for data validation in Golang. Below are the descriptions and examples for using each method.

## Installation

To install the validation package, run the following command:

```bash
go get github.com/KarimpourMehrab/golidation
```

## Usage


### Group Validation

### Group()
This function groups multiple validation error checks together. It likely returns a slice of errors or a combined validation result.

```bash
errors := validator.Group(
		validator.Attribute("firstName").Is(nil).Required().String().Errors(),
		validator.Attribute("lastName").Is(nil).Required().String().Errors(),
	)
```



### Basic Methods

### Required()  
Ensures the attribute is present and not empty.

```bash
validator.Attribute("name").Required()
```

### String()  
Ensures the attribute is a valid string.

```bash
validator.Attribute("name").String()
```

### Numeric()  
Ensures the attribute contains only numbers.

```bash
validator.Attribute("age").Numeric()
```

### Email()  
Ensures the attribute is a valid email address.

```bash
validator.Attribute("email").Email()
```

### Alpha()  
Ensures the attribute contains only alphabetic characters.

```bash
validator.Attribute("name").Alpha()
```

### AlphaNum()  
Ensures the attribute contains only letters and numbers.

```bash
validator.Attribute("username").AlphaNum()
```

### Date and Time Methods

### After(time)  
Ensures the attribute is after the specified time.

```bash
validator.Attribute("event_date").After(time.Now())
```

### Before(time)  
Ensures the attribute is before the specified time.

```bash
validator.Attribute("event_date").Before(time.Now())
```

### DateEquals(time)  
Ensures the attribute equals the specified date.

```bash
validator.Attribute("birth_date").DateEquals(time.Now())
```

### DateFormat(format)  
Ensures the attribute matches the specified date format.

```bash
validator.Attribute("birth_date").DateFormat("YYYY-MM-DD")
```

### Conditional Methods

### AcceptedIf(attribute, value)  
Ensures the attribute is accepted only if another attribute matches the specified value.

```bash
validator.Attribute("agreement").AcceptedIf("name", "mehrab")
```

### DeclinedIf(attribute, value)  
Ensures the attribute is declined only if another attribute matches the specified value.

```bash
validator.Attribute("agreement").DeclinedIf("name", "test")
```

### RequiredIf(attribute, value)  
Ensures the attribute is required if another attribute matches the specified value.

```bash
validator.Attribute("email").RequiredIf("name", "mehrab")
```

### RequiredUnless(attribute, values)  
Ensures the attribute is required unless another attribute matches one of the specified values.

```bash
validator.Attribute("phone").RequiredUnless("name", []interface{}{"mehrab"})
```

### Array and String Methods

### Array()  
Ensures the attribute is an array.

```bash
validator.Attribute("items").Array()
```

### InArray(values)  
Ensures the attribute is one of the values in the array.

```bash
validator.Attribute("status").InArray([]interface{}{"active", "inactive"})
```

### EndsWith(suffixes)  
Ensures the attribute ends with one of the specified suffixes.

```bash
validator.Attribute("url").EndsWith([]string{"com", "org"})
```

### StartsWith(prefixes)  
Ensures the attribute starts with one of the specified prefixes.

```bash
validator.Attribute("website").StartsWith([]string{"https://"})
```

### Regex and Format Methods

### Regex(pattern)  
Ensures the attribute matches the specified pattern.

```bash
validator.Attribute("email").Regex(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
```

### NotRegex(pattern)  
Ensures the attribute does not match the specified pattern.

```bash
validator.Attribute("password").NotRegex(`^[a-zA-Z0-9]*$`)
```

### Special Methods

### UUID()  
Ensures the attribute is a valid UUID.

```bash
validator.Attribute("user_id").UUID()
```

### IP()  
Ensures the attribute is a valid IP address.

```bash
validator.Attribute("ip_address").IP()
```

### IPv4()  
Ensures the attribute is a valid IPv4 address.

```bash
validator.Attribute("ip_address").IPv4()
```

### IPv6()  
Ensures the attribute is a valid IPv6 address.

```bash
validator.Attribute("ip_address").IPv6()
```

### JSON()  
Ensures the attribute is a valid JSON string.

```bash
validator.Attribute("data").JSON()
```

### File and Image Methods

### Image()  
Ensures the attribute is a valid image file.

```bash
validator.Attribute("profile_picture").Image()
```

### Mimes(mimeTypes)  
Ensures the attribute is a file with one of the specified MIME types.

```bash
validator.Attribute("profile_picture").Mimes([]string{"image/jpeg", "image/png"})
```

### Other Methods

### Boolean()  
Ensures the attribute is a boolean value.

```bash
validator.Attribute("is_active").Boolean()
```

### Unique(fields)  
Ensures the attribute is unique among the specified fields.

```bash
validator.Attribute("username").Unique([]interface{}{"username"})
```

### Exists(fields)  
Ensures the attribute exists among the specified fields.

```bash
validator.Attribute("user_id").Exists([]interface{}{"user"})
```

### Prohibited()  
Ensures the attribute is prohibited.

```bash
validator.Attribute("name").Prohibited()
```

### Filled()  
Ensures the attribute is filled (not empty).

```bash
validator.Attribute("email").Filled()
```

## Examples

Here is an example of how to combine multiple validation methods:

```bash
validator.Attribute("email").
    Required().
    Email().
    Unique([]interface{}{"email"}).
    NotRegex(`^[a-zA-Z0-9]*$`)
```

## License

MIT License. See the [LICENSE](LICENSE) file for details.