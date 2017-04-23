package display

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/xgb"
	"io/ioutil"
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
)

const (
	displayFile = "/etc/simple-session/display.json"
)

func Load(l *log.Logger) error {
	_conn, err := xgb.NewConn()
	if err != nil {
		return err
	}

	config, err := loadConfigInfo(displayFile)
	if err != nil {
		_conn.Close()
		return err
	}

	screenInfo, err := drandr.GetScreenInfo(_conn)
	if err != nil {
		_conn.Close()
		return err
	}

	conn = _conn
	logger = l
	sinfo = screenInfo

	outputs := screenInfo.Outputs.ListConnectionOutputs()
	outputs = filterByBlacklist(outputs, config.Blacklist)
	outputs = sortByPriority(outputs, config.Priority)

	return nil
}

func switchToExtend(outputs drandr.OutputInfos) {
}

func switchToMirror() {}

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
