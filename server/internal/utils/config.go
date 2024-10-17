package utils

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DBConn struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Dbname   string `yaml:"dbname"`
	} `yaml:"db_conn"`

	DBTables []string `yaml:"db_tables"`
}

func LoadConfig(path string) Config {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	var config Config
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}

	return config
}
