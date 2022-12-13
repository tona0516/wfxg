package repo

import (
	"encoding/json"
	"os"
	"path/filepath"
	"wfxg/vo"
)

type Local struct{}

func (l *Local) IsClientInstalled(installPath string) bool {
	replaysPath := filepath.Join(installPath, "replays")
	if f, err := os.Stat(replaysPath); os.IsNotExist(err) || !f.IsDir() {
		return false
	} else {
		return true
	}
}

func (l *Local) GetTempArenaInfo(installPath string) (vo.TempArenaInfo, error) {
	var tempArenaInfo vo.TempArenaInfo
	data, err := os.ReadFile(filepath.Join(installPath, "replays", "tempArenaInfo.json"))
	if err != nil {
		return tempArenaInfo, err
	}

	err = json.Unmarshal(data, &tempArenaInfo)
	if err != nil {
		return tempArenaInfo, err
	}

	return tempArenaInfo, nil
}
