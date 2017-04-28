package main

import (
	"os"
	"os/exec"
	"pkg.deepin.io/lib/dbus"
)

const (
	dbusDest = "org.jouyouyun.SimpleSession"
	dbusPath = "/org/jouyouyun/SimpleSession"
	dbusIFC  = dbusDest
)

func (*Manager) ToggleDebug() {
	doToggleDebug()
}

func (*Manager) Launch(cmd string) error {
	var handler = exec.Command("/bin/sh", "-c", cmd)
	err := handler.Start()
	if err != nil {
		logger.Error("[Launch] Failed to start:", err)
	}
	go func(cmd string) {
		err := handler.Wait()
		if err != nil {
			logger.Errorf("[Launch] Failed to wait(%s): %v", cmd, err)
		}
	}(cmd)
	return nil
}

func (m *Manager) Logout() {
	os.Exit(0)
}

func (m *Manager) ListOutput() string {
	return getOutputInfos()
}

func (*Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}
