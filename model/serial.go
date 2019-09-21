package model

import (
	"encoding/json"
	"github.com/offer365/endecrypt"
	"github.com/offer365/endecrypt/endeaesrsa"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SerialNum struct {
	ID    primitive.ObjectID `bson:"_id" json:"-"`
	Sid   string             `bson:"sid" json:"sid"`     // 序列号唯一uuid，用来标识序列号，并与 授权码相互校验，一一对应。
	Nodes map[string]*Node   `bson:"nodes" json:"nodes"` // 节点的具体硬件信息。这里不使用map的原因是map是无序的。无法保证每次生成的hws是一致的。
	Time  int64              `bson:"time" json:"time"`   // 生成 序列号的时间。
}

// 解密序列号
func Decrypt(src string) (sn *SerialNum, err error) {
	sn = new(SerialNum)
	sn.Nodes = make(map[string]*Node, 0)

	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	// 私钥解密
	if src, err = endeaesrsa.PirDecrypt(src, endecrypt.PirkeyServer2048, endecrypt.AesKeyServer2); err != nil {
		return
	}
	if src == "" {
		err = errors.Errorf("解密失败。")
		return
	}
	if err = json.Unmarshal([]byte(src), sn); err != nil {
		return
	}
	return
}
