package configs

import (
	"encoding/json"
	"io"
	"os"
)

type GlobalConfigs struct {
	Addr     string   `json:"addr"`
	Postgres pg_db    `json:"pg_db"`
	Redis    redis_db `json:"redis_db"`
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

func readConfigFile() string {
	file, err := os.Open("configs.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	return string(content)
}

func GetConfigs() GlobalConfigs {
	str := readConfigFile()

	var cfg GlobalConfigs
	json.Unmarshal([]byte(str), &cfg)

	return cfg
}
