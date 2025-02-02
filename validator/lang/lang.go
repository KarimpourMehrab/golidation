package lang

import (
	"github.com/KarimpourMehrab/golidation/internal/validator/lang/en"
	"github.com/KarimpourMehrab/golidation/internal/validator/lang/fa"
)

const (
	En = "en"
	Fa = "fa"
)

func TransMsg(lan interface{}, variable string) string {
	if lan == "fa" {
		return fa.Messages[variable]
	}
	return en.Messages[variable]
}
func TransErr(lan interface{}, variable string) string {
	if lan == "fa" {
		return fa.Errors[variable]
	}
	return en.Errors[variable]
}
func SysErr(variable string) string {
	return en.Errors[variable]
}
func SysInfo(variable string) string {
	return en.Info[variable]
}
