package logic

import (
	"context"
	"encoding/json"
	"github.com/offer365/edda/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"io/ioutil"
	"time"
)

func FindOneCustomer(coll string, id string) (instances []*model.Customer, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", oid}}
	instance := new(model.Customer)
	err = device.FindOne(coll, filter, instance)
	return []*model.Customer{instance}, err
}

func FindAllCustomer(coll string, skip, limit int64) (instances []*model.Customer, err error) {
	instances = make([]*model.Customer, 0)
	fu := func(cursor *mongo.Cursor) (err error) {
		// 遍历结果集
		ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
		for cursor.Next(ctx) {
			instance := new(model.Customer)
			if err = cursor.Decode(instance); err == nil { // 反序列化bson到对象
				instances = append(instances, instance)
			}
		}
		return
	}
	err = device.Find(coll, make(map[string]string), fu, skip, limit, -1)
	return
}

func InsertCustomer(coll string, body io.Reader) (id string, err error) {
	byt, err := ioutil.ReadAll(body)
	if err != nil {
		return
	}
	instance := new(model.Customer)
	err = json.Unmarshal(byt, instance)
	if err != nil {
		return
	}
	//instance.ID = primitive.NewObjectID()
	return device.Insert(coll, instance)
}
