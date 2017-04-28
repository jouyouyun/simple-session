package display

import (
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
	"pkg.deepin.io/dde/api/drandr"
	"sync"
)

var (
	evLocker sync.Mutex
)

func listenEvent() {
	randr.SelectInputChecked(conn, xproto.Setup(conn).DefaultScreen(conn).Root,
		randr.NotifyMaskOutputChange|randr.NotifyMaskOutputProperty|
			randr.NotifyMaskCrtcChange|randr.NotifyMaskScreenChange)
	for {
		e, err := conn.WaitForEvent()
		if err != nil {
			continue
		}

		evLocker.Lock()
		switch e.(type) {
		// case randr.NotifyCrtcChange:
		// case randr.NotifyOutputChange:
		// case randr.NotifyOutputProperty:
		case randr.ScreenChangeNotifyEvent:
			handleScreenChanged()
		}
		evLocker.Unlock()
	}
}

func handleScreenChanged() {
	screenInfo, err := drandr.GetScreenInfo(conn)
	if err != nil {
		logger.Error("[listenEvent] Failed to get screen info:", err)
		return
	}

	oldLen := len(sinfo.Outputs.ListConnectionOutputs())
	sinfo = screenInfo
	if oldLen != len(sinfo.Outputs.ListConnectionOutputs()) {
		handleOutputs()
	}
}
