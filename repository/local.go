package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Local struct{}

type TempArenaInfo struct {
	Vehicles []struct {
		ShipID   int    `json:"shipId"`
		Relation int    `json:"relation"`
		ID       int    `json:"id"`
		Name     string `json:"name"`
	} `json:"vehicles"`
}

func (l *Local) IsClientInstalled(installPath string) bool {
	replaysPath := filepath.Join(installPath, "replays")
	if f, err := os.Stat(replaysPath); os.IsNotExist(err) || !f.IsDir() {
		return false
	} else {
		return true
	}
}

func (l *Local) GetTempArenaInfo(installPath string) (TempArenaInfo, error) {
	var tempArenaInfo TempArenaInfo
	data, err := os.ReadFile(filepath.Join(installPath, "replays", "tempArenaInfo.json"))
	if err != nil {
		return tempArenaInfo, err
	}

	fmt.Println(string(data))

	err = json.Unmarshal(data, &tempArenaInfo)
	if err != nil {
		return tempArenaInfo, err
	}

	return tempArenaInfo, nil
}
