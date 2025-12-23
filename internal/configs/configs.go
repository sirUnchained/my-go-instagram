package configs

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

type GlobalConfigs struct {
	Addr          string        `json:"addr"`
	IsDevelopment bool          `json:"is_development"`
	Postgres      pg_db         `json:"pg_db"`
	Redis         redis_db      `json:"redis_db"`
	Auth          authenticator `json:"auth"`
}

type pg_db struct {
	Addr         string `json:"addr"`
	MaxOpenConns int    `json:"max_open_conns"`
	MaxIdleConns int    `json:"max_idle_conns"`
	MaxIdleTime  string `json:"max_idle_time"`
}

type redis_db struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DBNumber int    `json:"db_number"`
	Enabled  bool   `json:"enabled"`
}

type authenticator struct {
	SecretKey string        `json:"secret_key"`
	Aud       string        `json:"aud"`
	Iss       string        `json:"iss"`
	Exp       time.Duration `json:"exp"`
}

func readConfigFile() (string, error) {
	file, err := os.Open("configs.json")
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func GetConfigs() (*GlobalConfigs, error) {
	str, err := readConfigFile()
	if err != nil {
		return nil, err
	}

	var cfg GlobalConfigs
	json.Unmarshal([]byte(str), &cfg)

	return &cfg, nil
}
