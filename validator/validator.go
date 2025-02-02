package validator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"image"
	"net"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	Fa = "fa"
	En = "en"
)

type Validator struct {
	attribute      string
	value          interface{}
	optional       bool
	langMessages   map[string]string
	attributeTrans map[string]string
	messages       []string
}

func Attribute(attribute string) *Validator {
	var v Validator
	v.attribute = attribute
	return &v
}

func (v *Validator) Is(val interface{}) *Validator {
	v.value = val
	return v
}

func (v *Validator) EnMsg() *Validator {
	v.langMessages = EnValidationMsg
	return v
}

func (v *Validator) FaMsg() *Validator {
	v.langMessages = FaValidationMsg
	v.attributeTrans = FaAttributes
	return v
}

func (v *Validator) Lang(lan interface{}) *Validator {
	if lan == Fa {
		v.langMessages = FaValidationMsg
		v.attributeTrans = FaAttributes
		return v
	}

	v.langMessages = EnValidationMsg
	return v
}

func (v *Validator) messageMaker(rule string) *Validator {
	if len(v.langMessages) == 0 {
		v.langMessages = EnValidationMsg
	}
	val, ok := v.langMessages[rule]
	attributeTrans, okTrans := v.attributeTrans[v.attribute]
	if okTrans {
		v.attribute = attributeTrans
	}
	if ok {
		//fmt.Println("append new message...", v.attribute)
		v.messages = append(v.messages, strings.Replace(val, ":attr", v.attribute, 1))
	}
	return v
}

func (v *Validator) messageParamMaker(rule string, replacements map[string]string) *Validator {
	if len(v.langMessages) == 0 {
		v.langMessages = EnValidationMsg
	}
	val, ok := v.langMessages[rule]
	attributeTrans, okTrans := v.attributeTrans[v.attribute]
	if okTrans {
		v.attribute = attributeTrans
	}
	if ok {
		text := strings.Replace(val, ":attr", v.attribute, 1)
		for key, value := range replacements {
			placeholder := fmt.Sprintf("%s", key)
			text = strings.ReplaceAll(text, placeholder, value)
		}
		fmt.Println(text)

		v.messages = append(v.messages, text)
	}
	return v
}

func (v *Validator) Error() map[string]error {
	if len(v.messages) != 0 {
		fErrMsg := v.messages[0]
		return map[string]error{v.attribute: fmt.Errorf(fErrMsg)}
	}
	return nil
}

func (v *Validator) Errors() (errors map[string][]error) {
	errors = map[string][]error{}
	for _, message := range v.messages {
		_, ok := errors[v.attribute]
		if ok {
			errors[v.attribute] = append(errors[v.attribute], fmt.Errorf(message))
		} else {
			errors[v.attribute] = []error{fmt.Errorf(message)}
		}

	}
	return errors
}

type File interface {
	MimeType() string
}

func (v *Validator) Accepted() *Validator {
	var status bool
	switch val := v.value.(type) {
	case string:
		val = strings.TrimSpace(strings.ToLower(val))
		status = val == "yes" || val == "on" || val == "1" || val == "true"
	case bool:
		status = val
	case int:
		status = val == 1
	default:
		status = false
	}
	if !status {
		v.messageMaker("accepted")
	}
	return v
}

func (v *Validator) AcceptedIf(other string, otherValue interface{}) *Validator {
	var status bool

	// Check if other field's value matches the specified otherValue
	otherMatches := false
	switch val := otherValue.(type) {
	case string:
		otherMatches = strings.TrimSpace(strings.ToLower(other)) == strings.TrimSpace(strings.ToLower(val))
	case bool:
		otherMatches = other == "true" && val == true
	case int:
		otherMatches = other == "1" && val == 1
	}

	// If other field matches, then check if the current value is accepted
	if otherMatches {
		switch val := v.value.(type) {
		case string:
			val = strings.TrimSpace(strings.ToLower(val))
			status = val == "yes" || val == "on" || val == "1" || val == "true"
		case bool:
			status = val
		case int:
			status = val == 1
		default:
			status = false
		}
	} else {
		status = true // No validation needed if the condition is not met
	}
	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("accepted_if")
	}
	return v
}

func (v *Validator) ActiveURL() *Validator {
	var status bool

	// Attempt to parse the value as a URL
	if str, ok := v.value.(string); ok {
		_, err := url.ParseRequestURI(str)
		status = err == nil
	} else {
		status = false // If the value is not a string, it's not a valid URL
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("active_url")
	}
	return v
}

func (v *Validator) After(compareDate time.Time) *Validator {
	var status bool

	// Check if the value can be parsed as a date
	if dateStr, ok := v.value.(string); ok {
		// Parse the date string
		date, err := time.Parse("2006-01-02", dateStr)
		if err == nil {
			// Check if the date is after the compareDate
			status = date.After(compareDate)
		} else {
			status = false
		}
	} else {
		status = false // If the value is not a string, it's not a valid date
	}

	// If validation fails, call messageMaker
	if !status {

		replacements := map[string]string{":date": compareDate.String()}
		v.messageParamMaker("after", replacements)
	}
	return v
}

func (v *Validator) AfterOrEqual(compareDate time.Time) *Validator {
	var status bool

	// Check if the value can be parsed as a date
	if dateStr, ok := v.value.(string); ok {
		// Parse the date string
		date, err := time.Parse("2006-01-02", dateStr)
		if err == nil {
			// Check if the date is after or equal to the compareDate
			status = date.Equal(compareDate) || date.After(compareDate)
		} else {
			status = false
		}
	} else {
		status = false // If the value is not a string, it's not a valid date
	}

	// If validation fails, call messageMaker
	if !status {
		replacements := map[string]string{":date": compareDate.String()}
		v.messageParamMaker("after_or_equal", replacements)
	}
	return v
}

func (v *Validator) Alpha() *Validator {
	var status bool

	// Check if the value is a string and contains only alphabetic characters
	if str, ok := v.value.(string); ok {
		for _, ch := range str {
			if !('A' <= ch && ch <= 'Z' || 'a' <= ch && ch <= 'z') {
				status = false
				break
			}
		}
		status = true
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("alpha")
	}
	return v
}

func (v *Validator) AlphaDash() *Validator {
	var status bool

	// Check if the value is a string and contains only letters, numbers, and dashes
	if str, ok := v.value.(string); ok {
		for _, ch := range str {
			if !('A' <= ch && ch <= 'Z' || 'a' <= ch && ch <= 'z' || '0' <= ch && ch <= '9' || ch == '-') {
				status = false
				break
			}
		}
		status = true
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("alpha_dash")
	}
	return v
}

func (v *Validator) AlphaNum() *Validator {
	var status bool

	// Check if the value is a string and contains only letters and numbers
	if str, ok := v.value.(string); ok {
		for _, ch := range str {
			if !('A' <= ch && ch <= 'Z' || 'a' <= ch && ch <= 'z' || '0' <= ch && ch <= '9') {
				status = false
				break
			}
		}
		status = true
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("alpha_num")
	}
	return v
}

func (v *Validator) Array() *Validator {
	var status bool

	// Check if the value is of type slice (array)
	if _, ok := v.value.([]interface{}); ok {
		status = true
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("array")
	}
	return v
}

func (v *Validator) Before(compareDate time.Time) *Validator {
	var status bool

	// Check if the value can be parsed as a date
	if dateStr, ok := v.value.(string); ok {
		// Parse the date string
		date, err := time.Parse("2006-01-02", dateStr)
		if err == nil {
			// Check if the date is before the compareDate
			status = date.Before(compareDate)
		} else {
			status = false
		}
	} else {
		status = false // If the value is not a string, it's not a valid date
	}

	// If validation fails, call messageMaker
	if !status {
		replacements := map[string]string{":date": compareDate.String()}
		v.messageParamMaker("before", replacements)
	}
	return v
}

func (v *Validator) BeforeOrEqual(compareDate time.Time) *Validator {
	var status bool

	// Check if the value can be parsed as a date
	if dateStr, ok := v.value.(string); ok {
		// Parse the date string
		date, err := time.Parse("2006-01-02", dateStr)
		if err == nil {
			// Check if the date is before or equal to the compareDate
			status = date.Before(compareDate) || date.Equal(compareDate)
		} else {
			status = false
		}
	} else {
		status = false // If the value is not a string, it's not a valid date
	}

	// If validation fails, call messageMaker
	if !status {
		replacements := map[string]string{":date": compareDate.String()}
		v.messageParamMaker("before_or_equal", replacements)
	}
	return v
}

func (v *Validator) Boolean() *Validator {
	var status bool

	// Check if the value is a boolean
	if _, ok := v.value.(bool); ok {
		status = true
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("boolean")
	}
	return v
}

func (v *Validator) Confirmed(confirmValue interface{}) *Validator {
	var status bool

	// Check if the confirmation value matches the original value
	if confirmStr, ok := confirmValue.(string); ok {
		if originalStr, ok := v.value.(string); ok {
			status = originalStr == confirmStr
		} else {
			status = false
		}
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("confirmed")
	}
	return v
}

func (v *Validator) Date() *Validator {
	var status bool

	// Check if the value can be parsed as a date
	if dateStr, ok := v.value.(string); ok {
		_, err := time.Parse("2006-01-02", dateStr)
		status = err == nil
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("date")
	}
	return v
}

func (v *Validator) DateEquals(compareDate time.Time) *Validator {
	var status bool

	// Check if the value can be parsed as a date
	if dateStr, ok := v.value.(string); ok {
		// Parse the date string
		date, err := time.Parse("2006-01-02", dateStr)
		if err == nil {
			// Check if the date is equal to the compareDate
			status = date.Equal(compareDate)
		} else {
			status = false
		}
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		replacements := map[string]string{":date": compareDate.String()}
		v.messageParamMaker("date_equals", replacements)
	}
	return v
}

func (v *Validator) DateFormat(format string) *Validator {
	var status bool

	// Check if the value can be parsed as a date with the specified format
	if dateStr, ok := v.value.(string); ok {
		_, err := time.Parse(format, dateStr)
		status = err == nil
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		replacements := map[string]string{":format": format}
		v.messageParamMaker("date_format", replacements)
	}
	return v
}

func (v *Validator) Declined() *Validator {
	var status bool

	// Check if the value is declined (typically "no" or similar negative response)
	if str, ok := v.value.(string); ok {
		val := strings.TrimSpace(strings.ToLower(str))
		status = val == "no" || val == "declined" || val == "0" || val == "false"
	} else if boolVal, ok := v.value.(bool); ok {
		status = !boolVal
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("declined")
	}
	return v
}

func (v *Validator) DeclinedIf(otherValue interface{}, value interface{}) *Validator {
	var status bool

	// Check if the otherValue matches the specified value
	otherMatches := false
	switch val := otherValue.(type) {
	case string:
		otherMatches = strings.TrimSpace(strings.ToLower(val)) == strings.TrimSpace(strings.ToLower(value.(string)))
	case bool:
		otherMatches = val == value.(bool)
	case int:
		otherMatches = val == value.(int)
	}

	// If other field matches, then check if the current value is declined
	if otherMatches {
		if str, ok := v.value.(string); ok {
			val := strings.TrimSpace(strings.ToLower(str))
			status = val == "no" || val == "declined" || val == "0" || val == "false"
		} else if boolVal, ok := v.value.(bool); ok {
			status = !boolVal
		} else {
			status = false
		}
	} else {
		status = true // No validation needed if the condition is not met
	}

	// If validation fails, call messageMaker
	if !status {
		replacements := map[string]string{}
		otherValueStr, _ := otherValue.(string)
		replacements[":other"] = otherValueStr
		valueStr, _ := value.(string)
		replacements[":value"] = valueStr

		v.messageParamMaker("declined_if", replacements)
	}
	return v
}

func (v *Validator) Different(otherValue interface{}) *Validator {
	var status bool

	// Check if the value and otherValue are different
	switch val := v.value.(type) {
	case string:
		if otherStr, ok := otherValue.(string); ok {
			status = val != otherStr
		} else {
			status = false
		}
	case bool:
		if otherBool, ok := otherValue.(bool); ok {
			status = val != otherBool
		} else {
			status = false
		}
	case int:
		if otherInt, ok := otherValue.(int); ok {
			status = val != otherInt
		} else {
			status = false
		}
	default:
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		replacements := map[string]string{":other": otherValue.(string)}
		v.messageParamMaker("different", replacements)
	}
	return v
}

func (v *Validator) Digits(digits int) *Validator {
	var status bool

	// Check if the value is a string and has exactly the specified number of digits
	if str, ok := v.value.(string); ok {
		// Remove any non-digit characters
		digitCount := 0
		for _, ch := range str {
			if ch >= '0' && ch <= '9' {
				digitCount++
			}
		}
		status = digitCount == digits
	} else if num, ok := v.value.(int); ok {
		// For integer values, check the number of digits
		digitCount := len(fmt.Sprintf("%d", num))
		status = digitCount == digits
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		digitsStr := fmt.Sprintf("%d", digits)
		replacements := map[string]string{":digits": digitsStr}
		v.messageParamMaker("digits", replacements)
	}
	return v
}

func (v *Validator) DigitsBetween(min int, max int) *Validator {
	var status bool

	// Check if the value is a string and has a number of digits between min and max
	if str, ok := v.value.(string); ok {
		// Count the number of digit characters
		digitCount := 0
		for _, ch := range str {
			if ch >= '0' && ch <= '9' {
				digitCount++
			}
		}
		status = digitCount >= min && digitCount <= max
	} else if num, ok := v.value.(int); ok {
		// For integer values, check the number of digits
		digitCount := len(fmt.Sprintf("%d", num))
		status = digitCount >= min && digitCount <= max
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		minStr := fmt.Sprintf("%d", min)
		maxStr := fmt.Sprintf("%d", max)
		replacements := map[string]string{
			":min": minStr,
			":max": maxStr,
		}
		v.messageParamMaker("digits_between", replacements)
	}
	return v
}

func (v *Validator) Dimensions(minWidth, minHeight, maxWidth, maxHeight int) *Validator {
	var status bool

	// Check if the value is a byte slice representing an image
	if imgData, ok := v.value.([]byte); ok {
		// Decode the image
		_, _, err := image.DecodeConfig(bytes.NewReader(imgData))
		if err != nil {
			status = false
		} else {
			// Get the image dimensions
			img, _, err := image.Decode(bytes.NewReader(imgData))
			if err == nil {
				bounds := img.Bounds()
				width, height := bounds.Dx(), bounds.Dy()
				status = width >= minWidth && width <= maxWidth && height >= minHeight && height <= maxHeight
			} else {
				status = false
			}
		}
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("dimensions")
	}
	return v
}

func (v *Validator) Distinct() *Validator {
	var status bool

	// Check if the value is a slice and has distinct elements
	if slice, ok := v.value.([]interface{}); ok {
		seen := make(map[interface{}]bool)
		for _, item := range slice {
			if seen[item] {
				status = false
				break
			}
			seen[item] = true
		}
		status = true
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("distinct")
	}
	return v
}

func (v *Validator) Email() *Validator {
	var status bool

	// Define a regular expression for validating email addresses
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)

	// Check if the value is a valid email address
	if email, ok := v.value.(string); ok {
		status = re.MatchString(email)
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("email")
	}
	return v
}

func (v *Validator) EndsWith(values []string) *Validator {
	var status bool

	// Check if the value is a string and ends with one of the specified values
	if str, ok := v.value.(string); ok {
		for _, suffix := range values {
			if strings.HasSuffix(str, suffix) {
				status = true
				break
			}
		}
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		valuesStr := ""
		for _, value := range values {
			valuesStr += value
		}
		replacements := map[string]string{
			":values": valuesStr,
		}
		v.messageParamMaker("ends_with", replacements)
	}
	return v
}

func (v *Validator) Exists(validValues []interface{}) *Validator {
	var status bool

	// Check if the value exists in the list of valid values
	for _, validValue := range validValues {
		if v.value == validValue {
			status = true
			break
		}
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("exists")
	}
	return v
}

func (v *Validator) ExistsInString(validValues []string) *Validator {
	var status bool

	// Check if the value exists in the list of valid values
	for _, validValue := range validValues {
		if v.value == validValue {
			status = true
			break
		}
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("exists")
	}
	return v
}

func (v *Validator) Filled() *Validator {
	var status bool

	// Check if the value is not empty
	switch v.value.(type) {
	case string:
		status = strings.TrimSpace(v.value.(string)) != ""
	case []interface{}:
		status = len(v.value.([]interface{})) > 0
	case map[string]interface{}:
		status = len(v.value.(map[string]interface{})) > 0
	default:
		status = v.value != nil
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("filled")
	}
	return v
}

func (v *Validator) Image() *Validator {
	var status bool

	// Check if the value is a byte slice representing an image
	if imgData, ok := v.value.([]byte); ok {
		_, _, err := image.DecodeConfig(bytes.NewReader(imgData))
		status = err == nil
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("image")
	}
	return v
}

func (v *Validator) In(validValues []interface{}) *Validator {
	var status bool

	// Check if the value is in the list of valid values
	for _, validValue := range validValues {
		if reflect.DeepEqual(v.value, validValue) {
			status = true
			break
		}
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("in")
	}
	return v
}

func (v *Validator) Integer() *Validator {
	var status bool

	// Check if the value is an integer
	switch v.value.(type) {
	case int:
		status = true
	case string:
		_, err := strconv.Atoi(v.value.(string))
		status = err == nil
	default:
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("integer")
	}
	return v
}

func (v *Validator) InArray(arr []interface{}) *Validator {
	var status bool

	// Check if the value is in the list of otherValues
	for _, item := range arr {
		if reflect.DeepEqual(v.value, item) {
			status = true
			break
		}
	}

	// If validation fails, call messageMaker
	if !status {
		valuesStr := ""
		for _, value := range arr {
			valuesStr += value.(string)
		}
		replacements := map[string]string{
			":other": valuesStr,
		}
		v.messageParamMaker("in_array", replacements)
	}
	return v
}

func (v *Validator) IP() *Validator {
	var status bool

	// Check if the value is a valid IP address
	if ipStr, ok := v.value.(string); ok {
		status = net.ParseIP(ipStr) != nil
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("ip")
	}
	return v
}

func (v *Validator) IPv4() *Validator {
	var status bool

	// Check if the value is a valid IPv4 address
	if ipStr, ok := v.value.(string); ok {
		ip := net.ParseIP(ipStr)
		status = ip != nil && ip.To4() != nil
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("ipv4")
	}
	return v
}

func (v *Validator) IPv6() *Validator {
	var status bool

	// Check if the value is a valid IPv6 address
	if ipStr, ok := v.value.(string); ok {
		ip := net.ParseIP(ipStr)
		status = ip != nil && ip.To16() != nil && ip.To4() == nil
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("ipv6")
	}
	return v
}

func (v *Validator) JSON() *Validator {
	var status bool

	// Check if the value is a valid JSON string
	if jsonString, ok := v.value.(string); ok {
		var js json.RawMessage
		status = json.Unmarshal([]byte(jsonString), &js) == nil
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("json")
	}
	return v
}

func (v *Validator) Mimes(validTypes []string) *Validator {
	var status bool

	// Check if the value is a file with a valid MIME type
	if file, ok := v.value.(File); ok {
		// Extract the MIME type of the file
		fileType := file.MimeType()
		for _, validType := range validTypes {
			if strings.EqualFold(fileType, validType) {
				status = true
				break
			}
		}
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		valuesStr := ""
		for _, value := range validTypes {
			valuesStr += value
		}
		replacements := map[string]string{
			":values": valuesStr,
		}
		v.messageParamMaker("mimes", replacements)
	}
	return v
}

func (v *Validator) NotIn(invalidValues []interface{}) *Validator {
	var status bool

	// Check if the value is not in the list of invalidValues
	for _, invalidValue := range invalidValues {
		if reflect.DeepEqual(v.value, invalidValue) {
			status = false
			break
		}
		status = true
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("not_in")
	}
	return v
}

func (v *Validator) NotRegex(pattern string) *Validator {
	var status bool

	// Check if the value is a string and does not match the given regex pattern
	if str, ok := v.value.(string); ok {
		matched, err := regexp.MatchString(pattern, str)
		status = err == nil && !matched
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("not_regex")
	}
	return v
}

func (v *Validator) Numeric() *Validator {
	var status bool

	// Check if the value is a number
	switch v.value.(type) {
	case int, float32, float64, int32, int64:
		if v.value == 0 && !v.optional {
			status = false
		} else {
			status = true
		}

	case string:
		_, err := strconv.ParseFloat(v.value.(string), 64)
		status = err == nil
	default:
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("numeric")
	}
	return v
}

func (v *Validator) Present() *Validator {
	var status bool

	// Check if the value is present (not zero value)
	val := reflect.ValueOf(v.value)
	status = val.IsValid() && !reflect.DeepEqual(val.Interface(), reflect.Zero(val.Type()).Interface())

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("present")
	}
	return v
}

func (v *Validator) Prohibited() *Validator {
	var status bool

	// Check if the value is absent (zero value or nil)
	val := reflect.ValueOf(v.value)
	status = !val.IsValid() || reflect.DeepEqual(val.Interface(), reflect.Zero(val.Type()).Interface())

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("prohibited")
	}
	return v
}

func (v *Validator) ProhibitedIf(otherField string, otherValue interface{}) *Validator {
	var status bool

	// Check if the value is prohibited when the otherValue matches the specified condition
	val := reflect.ValueOf(v.value)
	otherVal := reflect.ValueOf(otherValue)

	// Check if the condition is met and the value is present (not zero value)
	if reflect.DeepEqual(v.value, otherVal) {
		status = !val.IsValid() || reflect.DeepEqual(val.Interface(), reflect.Zero(val.Type()).Interface())
	} else {
		status = true
	}

	// If validation fails, call messageMaker
	if !status {
		replacements := map[string]string{
			":other": otherField,
			":value": otherValue.(string),
		}
		v.messageParamMaker("prohibited_if", replacements)
	}
	return v
}

func (v *Validator) Regex(pattern string) *Validator {
	var status bool

	// Check if the value is a string and matches the given regex pattern
	if str, ok := v.value.(string); ok {
		matched, err := regexp.MatchString(pattern, str)
		status = err == nil && matched
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("regex")
	}
	return v
}

func (v *Validator) Required() *Validator {
	var status bool

	// Check if the value is present (not zero value)
	val := reflect.ValueOf(v.value)
	status = val.IsValid() && !reflect.DeepEqual(val.Interface(), reflect.Zero(val.Type()).Interface())

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("required")
	}
	return v
}

func (v *Validator) Optional() *Validator {
	v.optional = true
	return v
}

func (v *Validator) RequiredIf(otherField string, otherValue interface{}) *Validator {
	var status bool

	// Check if the value is required when the otherValue matches the specified condition
	val := reflect.ValueOf(v.value)
	otherVal := reflect.ValueOf(otherValue)

	// Check if the other condition is met
	if reflect.DeepEqual(v.value, otherVal) {
		status = val.IsValid() && !reflect.DeepEqual(val.Interface(), reflect.Zero(val.Type()).Interface())
	} else {
		status = true
	}

	// If validation fails, call messageMaker
	if !status {
		replacements := map[string]string{
			":other": otherField,
			":value": otherValue.(string),
		}
		v.messageParamMaker("required_if", replacements)
	}
	return v
}

func (v *Validator) RequiredUnless(otherValue interface{}, validValues []interface{}) *Validator {
	var status bool

	// Check if the value is required unless the otherValue is in the list of validValues
	val := reflect.ValueOf(v.value)
	otherVal := reflect.ValueOf(otherValue)
	valid := false

	for _, validValue := range validValues {
		if reflect.DeepEqual(otherVal.Interface(), validValue) {
			valid = true
			break
		}
	}

	if !valid {
		status = val.IsValid() && !reflect.DeepEqual(val.Interface(), reflect.Zero(val.Type()).Interface())
	} else {
		status = true
	}

	// If validation fails, call messageMaker
	if !status {
		validValuesStr := ""
		for _, value := range validValues {
			validValuesStr += value.(string) + ","
		}
		replacements := map[string]string{
			":other":  otherValue.(string),
			":values": validValuesStr,
		}
		v.messageParamMaker("required_unless", replacements)
	}
	return v
}

func (v *Validator) Same(otherValue interface{}) *Validator {
	var status bool

	// Check if the value and otherValue are the same
	if reflect.DeepEqual(v.value, otherValue) {
		status = true
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		replacements := map[string]string{
			":other": otherValue.(string),
		}
		v.messageParamMaker("same", replacements)
	}
	return v
}

func (v *Validator) StartsWith(values []string) *Validator {
	var status bool

	// Check if the value is a string and starts with any of the specified values
	if str, ok := v.value.(string); ok {
		for _, prefix := range values {
			if strings.HasPrefix(str, prefix) {
				status = true
				break
			}
		}
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		valuesStr := ""
		for _, value := range values {
			valuesStr += value
		}
		replacements := map[string]string{
			":values": valuesStr,
		}
		v.messageParamMaker("starts_with", replacements)
	}
	return v
}

func (v *Validator) String() *Validator {
	var status bool

	// Check if the value is a string
	if _, ok := v.value.(string); ok {
		status = true
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("string")
	}
	return v
}

func (v *Validator) Timezone() *Validator {
	var status bool

	// Check if the value is a string and is a valid timezone
	if tz, ok := v.value.(string); ok {
		_, err := time.LoadLocation(tz)
		status = err == nil
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("timezone")
	}
	return v

}

func (v *Validator) Unique(existingValues []interface{}) *Validator {
	var status bool

	// Check if the value is unique among the existing values
	for _, existingValue := range existingValues {
		if reflect.DeepEqual(v.value, existingValue) {
			status = false
			break
		}
		status = true
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("unique")
	}
	return v
}

func (v *Validator) URL() *Validator {
	var status bool

	// Check if the value is a string and is a valid URL
	if str, ok := v.value.(string); ok {
		_, err := url.ParseRequestURI(str)
		status = err == nil
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("url")
	}
	return v
}

func (v *Validator) UUID() *Validator {
	var status bool

	// Check if the value is a string and is a valid UUID
	if str, ok := v.value.(string); ok {
		_, err := uuid.Parse(str)
		status = err == nil
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("uuid")
	}
	return v
}

func (v *Validator) PasswordLetters() *Validator {
	var status bool

	// Check if the Value contains at least one letter
	if str, ok := v.value.(string); ok {
		for _, char := range str {
			if unicode.IsLetter(char) {
				status = true
				break
			}
		}
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("password_letters")
	}
	return v
}

func (v *Validator) PasswordMixed() *Validator {
	var status bool
	var hasUpper, hasLower bool

	// Check if the Value contains at least one uppercase and one lowercase letter
	if str, ok := v.value.(string); ok {
		for _, char := range str {
			if unicode.IsUpper(char) {
				hasUpper = true
			}
			if unicode.IsLower(char) {
				hasLower = true
			}
			if hasUpper && hasLower {
				status = true
				break
			}
		}
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("password_mixed")
	}
	return v
}

func (v *Validator) PasswordNumbers() *Validator {
	var status bool

	// Check if the Value contains at least one number
	if str, ok := v.value.(string); ok {
		for _, char := range str {
			if unicode.IsDigit(char) {
				status = true
				break
			}
		}
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("password_numbers")
	}
	return v
}

func (v *Validator) PasswordSymbols() *Validator {
	var status bool

	// Check if the Value contains at least one symbol
	if str, ok := v.value.(string); ok {
		for _, char := range str {
			if unicode.IsPunct(char) || unicode.IsSymbol(char) {
				status = true
				break
			}
		}
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("password_symbols")
	}
	return v
}

func (v *Validator) PasswordUncompromised(leakedValues []string) *Validator {
	var status bool

	// Check if the Value has not appeared in a data leak (in leakedValues)
	if str, ok := v.value.(string); ok {
		for _, leaked := range leakedValues {
			if str == leaked {
				status = false
				break
			}
		}
		if status != false {
			status = true
		}
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		v.messageMaker("password_uncompromised")
	}
	return v
}

func (v *Validator) MaxNumeric(max int) *Validator {
	var status bool

	// Get the value and ensure it's a numeric type
	val := reflect.ValueOf(v.value)

	if val.Kind() == reflect.Int || val.Kind() == reflect.Int8 || val.Kind() == reflect.Int16 ||
		val.Kind() == reflect.Int32 || val.Kind() == reflect.Int64 || val.Kind() == reflect.Uint ||
		val.Kind() == reflect.Uint8 || val.Kind() == reflect.Uint16 || val.Kind() == reflect.Uint32 ||
		val.Kind() == reflect.Uint64 || val.Kind() == reflect.Float32 || val.Kind() == reflect.Float64 {

		// Compare value against max
		var numericValue float64
		switch val.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			numericValue = float64(val.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			numericValue = float64(val.Uint())
		case reflect.Float32, reflect.Float64:
			numericValue = val.Float()
		}

		// Validate that numeric value does not exceed max
		status = numericValue <= float64(max)
	} else {
		// If the value is not numeric, validation fails
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		strVal := strconv.Itoa(max)
		replacements := map[string]string{":max": strVal}
		v.messageParamMaker("max_numeric", replacements)
	}

	return v
}

func (v *Validator) MaxString(max int) *Validator {
	var status bool

	// Get the value and ensure it's a string
	val := reflect.ValueOf(v.value)
	if val.Kind() == reflect.String {
		str := val.String()
		// Check if the length of the string exceeds max
		status = len(str) <= max
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		strVal := strconv.Itoa(max)
		replacements := map[string]string{":max": strVal}
		v.messageParamMaker("max_string", replacements)
	}

	return v
}

func (v *Validator) MinString(min int) *Validator {
	var status bool

	// Get the value and ensure it's a string
	val := reflect.ValueOf(v.value)
	if val.Kind() == reflect.String {
		str := val.String()
		// Check if the length of the string exceeds max
		status = len(str) >= min
	} else {
		status = false
	}

	// If validation fails, call messageMaker
	if !status {
		strVal := strconv.Itoa(min)
		replacements := map[string]string{":min": strVal}
		v.messageParamMaker("min_string", replacements)
	}

	return v
}
