package json_parser

import (
	"encoding/json"
	"os"
)

type Config struct {
	Db_url            string `json: "db_url"`
	Current_user_name string `json: "current_user_name"`
}

const configFileName = `.gatorconfig.json`

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	filePath := homeDir + configFileName
	return filePath, nil
}

func Read() (*Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
