package display

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/xgb"
	"io/ioutil"
	"os/exec"
	"pkg.deepin.io/dde/api/drandr"
	"pkg.deepin.io/lib/log"
	"pkg.deepin.io/lib/strv"
)

const (
	ModeExtern = "extern"
	ModeMirror = "mirror"
)

var (
	conn   *xgb.Conn
	logger *log.Logger
	sinfo  *drandr.ScreenInfo

	config *configInfo

	_failed = false
)

const (
	displayFile = "/etc/simple-session/display.json"
)

func Load(l *log.Logger) error {
	l.Info("Load display...")
	_conn, err := xgb.NewConn()
	if err != nil {
		_failed = true
		return err
	}

	screenInfo, err := drandr.GetScreenInfo(_conn)
	if err != nil {
		_conn.Close()
		_failed = true
		return err
	}

	conn = _conn
	logger = l
	sinfo = screenInfo

	config, err = loadConfigInfo(displayFile)
	if err != nil {
		logger.Warning("[Load] Failed:", err)
		switchToExtend(sinfo.Outputs.ListConnectionOutputs())
		return nil
	}

	handleOutputs()
	if config.AutoAdaptation == 1 {
		logger.Debug("Start listen output changed")
		go listenEvent()
	}
	return nil
}

func GetOutputInfos() string {
	if _failed {
		return ""
	}

	connected := sinfo.Outputs.ListConnectionOutputs()
	if config == nil {
		return toJson(&connected)
	}

	connected = filterByBlacklist(connected, config.Blacklist)
	connected = sortByPriority(connected, config.Priority)
	return toJson(&connected)
}

func handleOutputs() {
	outputs := sinfo.Outputs.ListConnectionOutputs()
	if len(outputs) == 0 {
		logger.Warning("[handleOutputs] No output connected")
		return
	}

	if config == nil {
		return
	}
	outputs = filterByBlacklist(outputs, config.Blacklist)
	outputs = sortByPriority(outputs, config.Priority)
	if config.Mode == ModeMirror {
		switchToMirror(outputs)
	} else {
		switchToExtend(outputs)
	}
}

func switchToExtend(outputs drandr.OutputInfos) {
	if len(outputs) == 0 {
		return
	}

	var cmd = "xrandr "
	x := uint16(0)
	for i, output := range outputs {
		cmd += fmt.Sprintf(" --output %s ", output.Name)
		if i == 0 {
			cmd += " --primary "
		}
		mode := getPreferredMode(output)
		cmd += fmt.Sprintf(" --pos %dx0 --mode %dx%d ", x, mode.Width, mode.Height)
		if mode.Rate > 0 {
			cmd += fmt.Sprintf(" --rate %f ", mode.Rate)
		}
		x += mode.Width
	}
	logger.Debug("[switchToExtend] cmd:", cmd)
	doAction(cmd)
}

func switchToMirror(outputs drandr.OutputInfos) {
	if len(outputs) == 0 {
		return
	}

	common := findCommonModes(outputs)
	if len(common) == 0 {
		return
	}

	var cmd = "xrandr "
	for i, output := range outputs {
		cmd += " --output " + output.Name
		if i == 0 {
			cmd += " --primary "
		}
		cmd += fmt.Sprintf(" --pos 0x0 --mode %dx%d ", common[0].Width, common[0].Height)
		if common[0].Rate > 0 {
			cmd += fmt.Sprintf(" --rate %f ", common[0].Rate)
		}
	}
	logger.Debug("[switchToExtend] cmd:", cmd)
	doAction(cmd)
}

func findCommonModes(outputs drandr.OutputInfos) drandr.ModeInfos {
	var mGroup []drandr.ModeInfos
	for _, output := range outputs {
		modes := getOutputModes(output)
		if len(modes) == 0 {
			continue
		}
		mGroup = append(mGroup, modes)
	}
	return drandr.FindCommonModes(mGroup...)
}

func getOutputModes(output drandr.OutputInfo) drandr.ModeInfos {
	var modes drandr.ModeInfos
	for _, id := range output.Modes {
		info := sinfo.Modes.Query(id)
		if info.Id != id {
			continue
		}
		modes = append(modes, info)
	}
	return modes
}

func getPreferredMode(output drandr.OutputInfo) drandr.ModeInfo {
	if output.Crtc.Id != 0 && output.Crtc.Width != 0 &&
		output.Crtc.Height != 0 {
		return sinfo.Modes.QueryBySize(output.Crtc.Width, output.Crtc.Height)
	}

	return sinfo.Modes.Query(output.Modes[0])
}

func filterByBlacklist(outputs drandr.OutputInfos, list []string) drandr.OutputInfos {
	if len(list) == 0 {
		return outputs
	}

	var infos drandr.OutputInfos
	for _, info := range outputs {
		if strv.Strv(list).Contains(info.Name) {
			continue
		}
		infos = append(infos, info)
	}
	return infos
}

func sortByPriority(outputs drandr.OutputInfos, list []string) drandr.OutputInfos {
	if len(list) == 0 {
		return outputs
	}

	var infos drandr.OutputInfos
	for _, v := range list {
		if info := outputs.QueryByName(v); info.Name == v {
			infos = append(infos, info)
		}
	}
	if len(infos) == 0 {
		return outputs
	}

	for _, info := range outputs {
		if tmp := infos.Query(info.Id); tmp.Id == info.Id {
			continue
		}
		infos = append(infos, info)
	}
	return infos
}

type configInfo struct {
	Mode string `json:"Mode"`
	// whether auto changed when output number changed
	AutoAdaptation int `json:"AutoAdaptation"`

	Blacklist []string `json:"Blacklist"`
	Priority  []string `json:"Priority"`
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

func toJson(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}
