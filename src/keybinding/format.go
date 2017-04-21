package keybinding

import (
	"fmt"
	"pkg.deepin.io/lib/strv"
	"strings"
)

const (
	accelDelim = "-"
)

var keysymModMap = map[string]string{
	"alt":       "mod1",
	"meta":      "mod1",
	"num_lock":  "mod2",
	"super":     "mod4",
	"hyper":     "mod4",
	"caps_lock": "lock",
}

var modKeysymMap = map[string]string{
	"mod1": "alt",
	"mod2": "num_lock",
	"mod4": "super",
	"lock": "caps_lock",
}

var keysymWiredMap = map[string]string{
	"-":  "minus",
	"=":  "equal",
	"\\": "backslash",
	"?":  "question",
	"!":  "exclam",
	"#":  "numbersign",
	";":  "semicolon",
	"'":  "apostrophe",
	"<":  "less",
	".":  "period",
	"/":  "slash",
	"(":  "parenleft",
	"[":  "bracketleft",
	")":  "parenright",
	"]":  "bracketright",
	"\"": "quotedbl",
	" ":  "space",
	"$":  "dollar",
	"+":  "plus",
	"*":  "asterisk",
	"_":  "underscore",
	"|":  "bar",
	"`":  "grave",
	"@":  "at",
	"%":  "percent",
	">":  "greater",
	"^":  "asciicircum",
	"{":  "braceleft",
	":":  "colon",
	",":  "comma",
	"~":  "asciitilde",
	"&":  "ampersand",
	"}":  "braceright",
}

var invalidKeysymList = strv.Strv{"caps_lock", "num_lock"}

// accel format, input: Control-Alt-T, output: control-mod1-t
func formatAccel(accel string) (string, error) {
	accel = strings.ToLower(accel)
	list := strings.Split(accel, accelDelim)
	if len(list) < 2 {
		return "", fmt.Errorf("Invalid accel: %s", accel)
	}

	var ret []string
	for i, v := range list {
		if i == len(list)-1 {
			wired, ok := keysymWiredMap[v]
			if ok {
				v = wired
			}
		} else {
			if invalidKeysymList.Contains(v) {
				continue
			}

			mod, ok := keysymModMap[v]
			if ok {
				v = mod
			}
		}
		ret = append(ret, v)
	}
	return strings.Join(ret, accelDelim), nil
}
