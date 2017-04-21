package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
)

const (
	configFile = "/etc/simple-session/config.json"
)

type Manager struct {
	config *configInfo
}

func NewManager() (*Manager, error) {
	config, err := loadConfigInfo(configFile)
	if err != nil {
		logger.Error("Load config file failed:", err)
	}
	return &Manager{
		config: config,
	}, nil
}

func (m *Manager) startSession() {
	// start pulseaudio
	m.Launch("pulseaudio")

	if m.config == nil {
		return
	}
	m.Launch(m.config.WM)
	m.launchScripts()
}

func (m *Manager) launchScripts() {
	scripts, err := getAutoScriptList(m.config.AutoScripts)
	if err != nil {
		logger.Error("[launchScripts] Failed to get script list:", err)
		return
	}
	for _, script := range scripts {
		m.Launch(script)
	}
}

func getAutoScriptList(dir string) ([]string, error) {
	finfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var scripts []string
	for _, finfo := range finfos {
		if finfo.IsDir() {
			continue
		}
		scripts = append(scripts, path.Join(dir, finfo.Name()))
	}
	return scripts, nil
}

type configInfo struct {
	WM          string `json:"WM"`
	AutoScripts string `json:"AutoScripts"`
	Background  string `json:"Background"`
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
