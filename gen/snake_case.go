package gen

import (
	"regexp"
	"strings"
)

//小写下划线
func LowerName(name string) string {
	if len(name) == 0 {
		return ""
	}
	name = regexp.MustCompile(`([\W]+)`).ReplaceAllString(name, "_")
	name = regexp.MustCompile(`([A-Z]+)`).ReplaceAllString(name, "_$1")
	name = strings.ToLower(name)
	name = strings.Trim(name, "_")
	if name[0] >= '0' && name[0] <= '9' {
		name = "_" + name
	}
	return name
}

//大写驼峰
func UpperName(name string) string {
	if len(name) == 0 {
		return ""
	}
	name = regexp.MustCompile(`([\W]+)`).ReplaceAllString(name, "_")
	name = regexp.MustCompile(`(\d+[a-z])`).ReplaceAllStringFunc(name, strings.ToUpper)
	name = regexp.MustCompile(`(_[A-Za-z0-9])`).ReplaceAllStringFunc(name, func(s string) string {
		return strings.ToUpper(s[1:2]) + s[2:]
	})
	if name[0] >= '0' && name[0] <= '9' {
		name = "_" + name
	} else {
		name = strings.ToUpper(name[:1]) + name[1:]
	}
	return name
}
