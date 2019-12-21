package eddaX

import (
	"errors"
	"reflect"

	"github.com/offer365/edda/log"
)

var (
	Cfg *Config
)

// 加解密方法
type CryptFunc func(src []byte) ([]byte, error)

// hash 方法
type HashFunc func(byt []byte) string

type Config struct {
	GRpcServerCrt  string
	GRpcServerKey  string
	GRpcClientCrt  string
	GRpcClientKey  string
	GRpcCaCrt      string
	GRpcUser       string
	GRpcPwd        string
	GRpcServerName string
	GRpcListen     string

	// encrypt decrypt func
	// odin & edda
	LicenseEncrypt CryptFunc // license 加解密
	LicenseDecrypt CryptFunc
	SerialEncrypt  CryptFunc // 序列号 加解密
	SerialDecrypt  CryptFunc
	UntiedEncrypt  CryptFunc // 解绑码 加解密
	UntiedDecrypt  CryptFunc
	TokenHash      HashFunc
}

func NewConfig() *Config {
	return &Config{
		GRpcServerCrt:  "",
		GRpcServerKey:  "",
		GRpcClientCrt:  "",
		GRpcClientKey:  "",
		GRpcCaCrt:      "",
		GRpcUser:       "",
		GRpcPwd:        "",
		GRpcServerName: "",
		GRpcListen:     "",
		LicenseEncrypt: nil,
		LicenseDecrypt: nil,
		SerialEncrypt:  nil,
		SerialDecrypt:  nil,
		UntiedEncrypt:  nil,
		UntiedDecrypt:  nil,
		TokenHash:      nil,
	}
}

func (cfg Config) CheckValue() (err error) {
	v := reflect.ValueOf(cfg)
	vt := v.Type()

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		switch vt.Field(i).Type.Name() {
		case "string":
			if f.String() == "" {
				return errors.New(vt.Field(i).Name + " field cannot be empty")
			}
		default:
			if f.IsNil() {
				return errors.New(vt.Field(i).Name + " field cannot be nil")
			}
		}
	}
	return
}

func Start(cfg *Config) {
	if err := cfg.CheckValue(); err != nil {
		log.Sugar.Fatal(err)
		return
	}
	Cfg = cfg
	Main()
	return
}
