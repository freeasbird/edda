package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Customer struct {
	Name    string `bson:"name" json:"name"`       // 客户姓名
	Info    string `bson:"info" json:"info"`       // 客户信息
	Project string `bson:"project" json:"project"` // 相关项目
}

type APPs struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	Name         string             `bson:"name" json:"name"`                 // 客户姓名
	Key          string             `json:"key" json:"key"`                   // 客户信息
	Introduction string             `bson:"introduction" json:"introduction"` // 相关项目
	Attrs        []*Att             `bson:"attrs" json:"attrs"`
}

type Att struct {
	Name string `bson:"name" json:"name"`
	Key  string `bson:"key" json:"key"`
}
