package logic

import (
	"../model"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"io/ioutil"
	"time"
)

func FindOneNode(coll string, id string) (instances []*model.Node, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", oid}}
	instance := new(model.Node)
	err = device.FindOne(coll, filter, instance)
	return []*model.Node{instance}, err
}

func FindNode(coll string, filter interface{}, skip, limit int64) (instances []*model.Node, err error) {
	instances = make([]*model.Node, 0)
	fu := func(cursor *mongo.Cursor) (err error) {
		// 遍历结果集
		ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
		for cursor.Next(ctx) {
			instance := new(model.Node)
			if err = cursor.Decode(instance); err == nil { // 反序列化bson到对象
				instances = append(instances, instance)
			}
		}
		return
	}
	err = device.Find(coll, filter, fu, skip, limit, 1)
	return
}

func InsertNode(coll string, body io.Reader) (id string, err error) {
	byt, err := ioutil.ReadAll(body)
	if err != nil {
		return
	}
	instance := new(model.Node)
	err = json.Unmarshal(byt, instance)
	if err != nil {
		return
	}
	instance.ID = primitive.NewObjectID()
	return device.Insert(coll, instance)
}
