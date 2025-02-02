package main

import (
	"fmt"
	"github.com/KarimpourMehrab/golidation/validator"
	"time"
)

func main() {

	errors := validator.
		Attribute("pan").
		Is(nil).
		Required().
		String().
		Numeric().
		Accepted().
		AcceptedIf("name", "mehrab").
		ActiveURL().
		After(time.Now()).
		AfterOrEqual(time.Now()).
		Alpha().
		AlphaDash().
		AlphaNum().
		Array().
		Before(time.Now()).
		BeforeOrEqual(time.Now()).
		Boolean().
		Confirmed("name").
		Date().
		DateEquals(time.Now()).
		DateFormat("YYY-mm-dd HH:mm:ss").
		Different("test").
		DeclinedIf("name", "test").
		Digits(10).
		DigitsBetween(1, 20).
		Dimensions(1, 1, 1, 1).
		Distinct().
		Email().
		EndsWith([]string{"tcp"}).
		Exists([]interface{}{"ip"}).
		ExistsInString([]string{"test"}).
		Filled().
		Image().
		In([]interface{}{"mehrab"}).
		Integer().
		InArray([]interface{}{"test"}).
		IP().
		IPv4().
		IPv6().
		JSON().
		Mimes([]string{"test"}).
		NotIn([]interface{}{"test"}).
		NotRegex(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).
		Numeric().
		Present().
		Prohibited().
		ProhibitedIf("name", "mehrab").
		Regex(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).
		RequiredIf("name", "mehrab").
		RequiredUnless("name", []interface{}{"mehrab"}).
		Same("").
		StartsWith([]string{"https://"}).
		Timezone().
		Unique([]interface{}{"mehrab"}).
		URL().
		UUID().
		PasswordLetters().
		PasswordMixed().
		PasswordNumbers().
		PasswordSymbols().
		PasswordUncompromised([]string{"123456"}).
		MaxNumeric(120).
		MaxString(255).
		MinString(1).
		Errors()

	for attribute, err := range errors {
		fmt.Println(attribute, err)
	}
}
