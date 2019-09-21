package config

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

var CFG *Configuration

func init() {
	CFG = new(Configuration)
	CFG.Users = make(map[string]string, 0)
}

type Configuration struct {
	Port    string            `json:"port"`
	MongoDB MongoDB           `json:"mongodb"`
	Users   map[string]string `json:"users"`
}

type MongoDB struct {
	Host        string   `json:"host"`
	Port        string   `json:"port"`
	User        string   `json:"user"`
	Pwd         string   `json:"pwd"`
	Database    string   `json:"database"`
	Collections []string `json:"collections"`
}

func (cfg *Configuration) LoadConfig(filename string) {
	var (
		content []byte
		err     error
	)
	//读取配置文件
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
	cfg.MongoDB.Host = "127.0.0.1"
	cfg.MongoDB.Port = "27017"
	cfg.MongoDB.User = "admin"
	cfg.MongoDB.Pwd = "666666"
	cfg.MongoDB.Database = "products"
	cfg.MongoDB.Collections = []string{"principals", "projects", "copyrights", "products"}
	log.Error("failed to load configuration file. error:", err.Error())
	return
}
