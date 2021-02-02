// Package generate_model Create at 2019-06-18 9:35
package generate_model

import (
	"fmt"
	"strings"
)

// camel string, xx_yy to XxYy
func camelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

// snake string, XxYy to xx_yy , XxYY to xx_yy
func snakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

func sqlTypeToGoType(t string) string {
	switch t {
	case "VARCHAR":
		return "*string"
	case "BIGINT":
		return "*int64"
	case "TINYINT", "INT":
		return "*int"
	case "DOUBLE":
		return "*float64"
	default:
		return ""
	}
}

func sqlTypeToHandlerType(t string) string {
	switch t {
	case "VARCHAR":
		return "string"
	case "BIGINT":
		return "*int64"
	case "TINYINT", "INT":
		return "*int"
	case "DOUBLE":
		return "*float64"
	default:
		return ""

	}
}

func sqlParamToGoParam(t string) string {
	if "id" == t {
		return "ID"
	} else if strings.Contains(t, "_") {
		return camelString(t)
	} else {
		return fmt.Sprintf("%s%s", strings.ToUpper(t[:1]), t[1:])
	}
}
