package json_parser

import (
	"encoding/json"
	"os"
	"sync"
)

type Config struct {
	Db_url            string `json: "db_url"`
	Current_user_name string `json: "current_user_name"`
}

var configLock sync.RWMutex

const configFileName = `.gatorconfig.json`

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	filePath := homeDir + "/" + configFileName
	return filePath, nil
}

func Read() (Config, error) {
	path, err := getConfigFilePath()
	configLock.RLock()
	defer configLock.RUnlock()
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, err
	}
	return config, nil
}

func SetUser(config Config, user string) (Config, error) {
	config.Current_user_name = user
	writeCfg(config)
	return config, nil
}

func writeCfg(cfg Config) error {
	configLock.Lock()
	defer configLock.Unlock()
	jsonBytes, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	err = os.WriteFile(path, jsonBytes, 0600)
	if err != nil {
		return err
	}
	return nil
}
