package model

import (
	"encoding/json"
	"errors"
	"github.com/offer365/endecrypt"
	"github.com/offer365/endecrypt/endeaesrsa"
	"github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"sync"
	"time"
)

// 授权码
type License struct {
	ID primitive.ObjectID `bson:"_id" json:"id"`
	Auth
	Cipher    string `bson:"cipher" json:"cipher"`
	SerialNum string `bson:"serial_num" json:"serial_num"`
	Customer  string `bson:"customer" json:"customer"`
	Project   string `bson:"project" json:"project"`
}

type Auth struct {
	Lid            string            `bson:"lid" json:"lid"`                              // 授权码唯一uuid,用来甄别是否重复授权。
	Sid            string            `bson:"sid" json:"sid"`                              // 机器码的id, lid与sid 一一对应
	Devices        map[string]string `bson:"devices" json:"devices"`                      // 设备与 硬件信息md5
	GenerationTime time.Time         `bson:"generation_time" json:"generation_time"`      // 授权生成时间
	UpdateTime     time.Time         `bson:"update_time" title:"更新时间" json:"update_time"` //当前时间 最后一次授权更新时间
	APPs           map[string]*APP   `bson:"apps"  title:"产品" json:"apps"`                //key:app英文名请求中url标识字段
}

//type Device struct {
//	Name string `bson:"name" json:"name"`
//	Md5  string `bson:"md5" json:"md5"`
//}

// 应用
type APP struct {
	Key          string         `bson:"key" json:"key"`
	Name         string         `bson:"name" title:"服务" json:"name"`
	Attr         map[string]int `bson:"attr" json:"attr"`                                    // 属性
	Introduction string         `bson:"introduction" json:"introduction"`                    // 简介
	Instance     int            `bson:"instance" title:"最大实例" json:"instance"`               // 实例
	ExpireTime   time.Time      `bson:"expire_time" title:"到期时间" json:"expire_time"`         // 授权到期的时间戳
	MaxLifeCycle int64          `bson:"max_life_cycle" title:"最大生存周期" json:"max_life_cycle"` // 最大生存周期 (授权到期时间-生成授权时间)/周期时间60s

	rv reflect.Value
	rt reflect.Type
	mu sync.RWMutex
}

// 生成密文
func (l *License) CipherText() {
	var (
		byt []byte
		err error
	)
	if byt, err = json.Marshal(l.Auth); err != nil {
		return
	}
	// 私钥加密
	if l.Cipher, err = endeaesrsa.PriEncrypt(byt, endecrypt.PirkeyServer2048, endecrypt.AesKeyServer2); err != nil {
		return
	}
}

// 创建 license
func NewLicense(sid string, devices map[string]string, apps ...*APP) (lic *License, err error) {
	// 如果没有服务器硬件信息，序列号错误
	if len(devices) == 0 {
		err = errors.New("Serial number hardware error.")
		return
	}
	lic = new(License)
	lic.Lid = uuid.Must(uuid.NewV4()).String()
	lic.Sid = sid
	lic.Devices = devices
	lic.APPs = make(map[string]*APP, 0)
	lic.UpdateTime = time.Now()
	lic.GenerationTime = time.Now()
	for _, app := range apps {
		lic.APPs[app.Key] = app
	}
	return
}

func NewAPP(key string, name string, expired string, instances int64, attr map[string]int) (app *APP) {
	app = new(APP)
	app.Key = key
	app.Name = name
	exp, err := time.Parse("2006-01-02", expired)
	if err != nil {
		return
	}
	app.ExpireTime = exp
	app.Instance = int(instances)
	app.MaxLifeCycle = int64(exp.Unix()-time.Now().Unix()) / 60
	app.Attr = attr
	return
}
