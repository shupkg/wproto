package gen

import (
	"strings"
	"unicode"
)

// func Leading(prefix string, leading protogen.CommentSet, breakLine bool) string {
// 	w := &bytes.Buffer{}
// 	if len(leading.LeadingDetached) > 0 {
// 		for _, comment := range leading.LeadingDetached {
// 			c := TrimComment(comment, "")
// 			if c != "" {
// 				w.WriteString(c)
// 				w.WriteByte('\n')
// 				w.WriteByte('\n')
// 			}
// 		}
// 	}

// 	if leading.Leading != "" {
// 		w.WriteString(TrimComment(leading.Leading, prefix))
// 		if breakLine {
// 			w.WriteByte('\n')
// 		}
// 	}
// 	return w.String()
// }

// func Trailing(raw protogen.Comments) (tagMap map[string]string, comment string) {
// 	s := TrimComment(raw, "")
// 	if s == "" {
// 		return
// 	}
// 	ss := strings.SplitN(s[2:], "//", 2)
// 	for _, c := range ss {
// 		c = strings.TrimSpace(c)
// 		if IsWrap(c) {
// 			tagMap = ParseTags(c)
// 		} else {
// 			comment = "//" + c
// 		}
// 	}
// 	return
// }
//
// func TrimComment(raw protogen.Comments, prefix string) string {
// 	if raw == "" {
// 		return ""
// 	}
// 	s := Trim(string(raw), '/')
// 	if s != "" {
// 		if prefix != "" && !strings.HasPrefix(s, prefix) {
// 			return "//" + prefix + " " + s
// 		}
// 		return "//" + s
// 	}
// 	return s
// }

//name:phone;required
func ParseTags(tags ...string) map[string]string {
	tagMap := make(map[string]string)
	for _, tag := range tags {
		items := strings.Split(Trim(tag, '`', ';'), ";")
		for _, item := range items {
			splitItem := strings.SplitN(strings.TrimSpace(item), ":", 2)
			var (
				name string
				val  string
			)
			if len(splitItem) > 0 {
				name = strings.TrimSpace(splitItem[0])
			}
			if len(splitItem) > 1 {
				val = strings.TrimSpace(splitItem[1])
			}
			if name != "" {
				tagMap[name] = val
			}
		}
	}
	return tagMap
}

func IsWrap(s string, quotes ...string) bool {
	var (
		start = "`"
		end   = "`"
	)
	if len(quotes) > 0 {
		start = quotes[0]
		end = quotes[1]
	}
	if len(quotes) > 1 {
		end = quotes[1]
	}
	return strings.HasPrefix(s, start) && strings.HasSuffix(s, end)
}

func Wrap(s string, quotes ...string) string {
	var (
		start = "`"
		end   = "`"
	)
	if len(quotes) > 0 {
		start = quotes[0]
		end = quotes[1]
	}
	if len(quotes) > 1 {
		end = quotes[1]
	}
	return start + s + end
}

func Trim(s string, cuts ...rune) string {
	s = strings.TrimFunc(s, func(r rune) bool {
		if unicode.IsSpace(r) {
			return true
		}
		for _, cut := range cuts {
			if r == cut {
				return true
			}
		}
		return false
	})

	return s
}
