package keybinding

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"io/ioutil"
	"os/exec"
	"pkg.deepin.io/lib/log"
	"sort"
	"sync"
)

var (
	xu     *xgbutil.XUtil
	logger *log.Logger
	locker sync.Mutex

	keycodeMap = make(map[string]*shortcutInfo)
)

const (
	keybindingFile = "/etc/simple-session/keybinding.json"
)

type configInfo struct {
	List shortcutInfos `json:"Shortcuts"`
}

type shortcutInfo struct {
	Shortcut string `json:"Shortcut"`
	Action   string `json:"Action"`
}
type shortcutInfos []*shortcutInfo

func Load(l *log.Logger) error {
	l.Info("Load keybinding...")
	_xu, err := xgbutil.NewConn()
	if err != nil {
		return err
	}

	config, err := loadConfigInfo(keybindingFile)
	if err != nil {
		return err
	}

	xu = _xu
	logger = l
	keybind.Initialize(xu)
	locker.Lock()
	for _, info := range config.List {
		id, err := grabAccel(info.Shortcut)
		if err != nil {
			logger.Warningf("Failed to grab '%s', reason: %v",
				info.Shortcut, err)
			continue
		}
		keycodeMap[id] = info
	}
	locker.Unlock()
	return nil
}

func StartLoop() {
	xevent.KeyPressFun(func(_xu *xgbutil.XUtil, ev xevent.KeyPressEvent) {
		handleKeyPress(ev.State, ev.Detail)
	}).Connect(xu, xu.RootWin())
}

func Terminate() {
	xevent.Quit(xu)
	// xu = nil
	locker.Lock()
	keycodeMap = nil
	locker.Unlock()
}

func handleKeyPress(state uint16, detail xproto.Keycode) {
	modStr := keybind.ModifierString(state)
	keyStr := keybind.LookupString(xu, state, detail)
	if detail == 65 {
		keyStr = "space"
	}
	if modStr != "" {
		keyStr = modStr + "-" + keyStr
	}

	logger.Info("[handleKeyPress] event key:", keyStr)
	accel, err := formatAccel(keyStr)
	if err != nil {
		logger.Warning("[handleKeyPress] Failed to format:", keyStr, err)
		return
	}

	id, err := getAccelId(accel)
	if err != nil {
		logger.Warning("[handleKeyPress] Failed to get id:", accel, err)
		return
	}
	locker.Lock()
	info, ok := keycodeMap[id]
	locker.Unlock()
	if !ok || info == nil {
		logger.Debug("[handleKeyPress] No shortcut info found for:", id)
		return
	}
	logger.Debug("[handleKeyPress] Will exec:", info.Action)
	doAction(info.Action)
}

func grabAccel(accel string) (string, error) {
	tmp, err := formatAccel(accel)
	if err != nil {
		return "", err
	}

	mod, codes, err := keybind.ParseString(xu, tmp)
	if err != nil {
		return "", err
	}

	var hasGrab []int
	for _, code := range codes {
		err = keybind.GrabChecked(xu, xu.RootWin(), mod, code)
		if err != nil {
			break
		}
		hasGrab = append(hasGrab, int(code))
	}

	if len(hasGrab) != len(codes) {
		// failed
		for _, code := range hasGrab {
			keybind.Ungrab(xu, xu.RootWin(), mod, xproto.Keycode(code))
		}
		return "", err
	}
	return doGetAccelId(mod, hasGrab), nil
}

func getAccelId(accel string) (string, error) {
	mod, codes, err := keybind.ParseString(xu, accel)
	if err != nil {
		return "", err
	}
	var tmp []int
	for _, code := range codes {
		tmp = append(tmp, int(code))
	}
	return doGetAccelId(mod, tmp), nil
}

func doGetAccelId(mod uint16, codes []int) string {
	var id = fmt.Sprintf("%v", mod)
	sort.Ints(codes)
	for _, code := range codes {
		id += fmt.Sprintf("-%v", code)
	}
	return id
}

func doAction(cmd string) {
	app := exec.Command("/bin/sh", "-c", cmd)
	err := app.Start()
	if err != nil {
		logger.Warning("Failed to start cmd:", cmd, err)
		return
	}
	go func(c string) {
		err := app.Wait()
		if err != nil {
			logger.Warning("Failed to wait cmd:", c, err)
			return
		}
	}(cmd)
}

func loadConfigInfo(file string) (*configInfo, error) {
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if len(contents) == 0 {
		return nil, fmt.Errorf("file(%s) is empty", file)
	}

	var info configInfo
	err = json.Unmarshal(contents, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}
