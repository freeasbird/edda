package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

var (
	cfg string
	Cfg *Configuration
)

func args() {
	flag.StringVar(&cfg, "f", "edda.json", "Cfg file path.")
	flag.Parse()
}

func init() {
	args()
	Cfg = new(Configuration)
	Cfg.Users = make(map[string]string, 0)
	Cfg.LoadConfig(cfg)
}

type Configuration struct {
	Port    string            `json:"port"`
	Core    string            `json:"core"`
	Users   map[string]string `json:"users"`
}



func (cfg *Configuration) LoadConfig(filename string) {
	var (
		content []byte
		err     error
	)
	// 读取配置文件
	if content, err = ioutil.ReadFile(filename); err != nil {
		goto ERR
	}
	// json反序列化
	if err = json.Unmarshal(content, cfg); err != nil {
		goto ERR
	}

	return
ERR:
	cfg.Port = "1999"
	log.Error("failed to load configuration file. error:", err.Error())
	return
}
