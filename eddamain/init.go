package eddamain

import (
	"flag"
	"github.com/offer365/edda/asset"
	"github.com/offer365/edda/config"
	"github.com/offer365/edda/logic"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	ConfFilePath string

	debug     bool
	AssetPath string
)

func args() {
	flag.StringVar(&ConfFilePath, "f", "config.json", "CFG file path.")
	flag.Parse()
}

// 释放静态资源
func RestoreAsset() {
	// 解压 静态文件的位置
	if runtime.GOOS == "linux" {
		AssetPath = "/usr/share/.asset/.temp/"
	} else {
		AssetPath = "./"
	}
	// go get -u github.com/jteeuwen/go-bindata/...
	// 重新生成静态资源在项目的根目录下 go-bindata -o=asset/asset.go -pkg=asset views/... static/...
	dirs := []string{"views", "static"}
	for _, dir := range dirs {
		if err := asset.RestoreAssets(AssetPath, dir); err != nil {
			log.Error("RestoreAssets error. ", err.Error())
			_ = os.RemoveAll(filepath.Join(AssetPath, dir))
			continue
		}
	}
}

func init() {
	args()
	RestoreAsset()
	log.SetFormatter(&log.JSONFormatter{})
	debug = true
	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}
	config.CFG.LoadConfig(ConfFilePath)
	cfg := config.CFG
	debug = true
	collIndex := make(map[string]string, 0)
	collIndex["customers"] = "name"
	collIndex["apps"] = "name"
	collIndex["nodes"] = "attr.md5"
	collIndex["licenses"] = "auth.lid"
	collIndex["serial"] = "sid"

	if err := logic.Init(cfg.MongoDB.Host, cfg.MongoDB.Port, cfg.MongoDB.User, cfg.MongoDB.Pwd, cfg.MongoDB.Database, 2*time.Second, collIndex); err != nil {
		log.Fatal("init error.")
	}
}
