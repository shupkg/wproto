package gen

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type Param map[string]string

func (p Param) Encode() string {
	var b bytes.Buffer
	for key, val := range p {
		if b.Len() > 0 {
			b.WriteByte(',')
		}
		fmt.Fprintf(&b, "%s=%s", key, val)
	}
	return b.String()
}

func (p Param) Set(key, val string) {
	p[key] = val
}

func (p Param) GetBool(key string) bool {
	yes, _ := strconv.ParseBool(p[key])
	return yes
}

func (p Param) Get(key string) string {
	return p[key]
}

func ParseParam(s string) Param {
	p := Param{}
	items := strings.Split(s, ",")
	for _, item := range items {
		kv := strings.SplitN(item, "=", 2)
		if len(kv) == 2 {
			p[kv[0]] = kv[1]
		}
	}
	return p
}
