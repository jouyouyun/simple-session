package keybinding

import (
	"fmt"
	"testing"
)

func TestFormatAccel(t *testing.T) {
	var infos = []struct {
		accel string
		ret   string
	}{
		{
			accel: "Super-T",
			ret:   "mod4-t",
		},
		{
			accel: "Control-Super-T",
			ret:   "control-mod4-t",
		},
		{
			accel: "Super-Caps_Lock-T",
			ret:   "mod4-t",
		},
		{
			accel: "Super-Num_Lock-T",
			ret:   "mod4-t",
		},
	}

	for _, info := range infos {
		accel, _ := formatAccel(info.accel)
		if accel != info.ret {
			msg := fmt.Sprintf("Expected: %s; but got: %s", info.ret, accel)
			panic(msg)
		}
	}
}

func TestDoGetAccelId(t *testing.T) {
	var infos = []struct {
		mod   uint16
		codes []int
		ret   string
	}{
		{
			mod:   10,
			codes: []int{5, 3, 1, 7},
			ret:   "10-1-3-5-7",
		},
		{
			codes: []int{5, 3, 1, 7},
			ret:   "0-1-3-5-7",
		},
		{
			mod: 10,
			ret: "10",
		},
	}
	for _, info := range infos {
		id := doGetAccelId(info.mod, info.codes)
		if id != info.ret {
			msg := fmt.Sprintf("Expected: %s; but got: %s", info.ret, id)
			panic(msg)
		}
	}
}

func TestLoadFile(t *testing.T) {
	var ret = shortcutInfos{
		{
			Shortcut: "Super-T",
			Action:   "xterm",
		},
		{
			Shortcut: "Control-Super-Delete",
			Action:   "logout",
		},
	}
	config, _ := loadConfigInfo("testdata/keybinding.json")
	for i, info := range config.List {
		if ret[i].Shortcut != info.Shortcut || ret[i].Action != info.Action {
			msg := fmt.Sprintf("Expected: (%s, %s); but got: (%s, %s)",
				ret[i].Shortcut, ret[i].Action, info.Shortcut, info.Action)
			panic(msg)
		}
	}
}
