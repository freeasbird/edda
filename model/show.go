package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Show struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Customer string             `bson:"customer" json:"customer"` //客户
	APPs     []string           `bson:"apps" json:"apps"`         // 被授权了哪些产品
	Attr     string             `bson:"attr" json:"attr"`         // 授权属性  测试或正式
	Duration time.Duration      `bson:"duration" json:"duration"` // 授权时长
	Expire   time.Time          `bson:"expire" json:"expire"`     // 到期时间
}
