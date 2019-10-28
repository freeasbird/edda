package logic

import (
	"context"
	"encoding/json"
	pb "github.com/offer365/edda/eddacore/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"io/ioutil"
	"time"
)

type APP struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	Name         string             `bson:"name" json:"name"`
	Key          string             `bson:"key" json:"key"`
	Attrs        []*pb.Attr         `bson:"attrs" json:"attrs"`
	Expire       int64              `bson:"expire" json:"expire"`
	Instance     int64              `bson:"instance" json:"instance"`
	MaxLifeCycle int64              `bson:"maxLifeCycle" json:"maxLifeCycle"`
}

func FindOneApp(coll string, id string) (instance *APP, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", oid}}
	instance = new(APP)
	err = db.FindOne(coll, filter, instance)
	return
}

func FindAllApp(coll string, skip, limit int64) (instances []*APP, err error) {
	instances = make([]*APP, 0)
	fu := func(cursor *mongo.Cursor) (err error) {
		// 遍历结果集
		ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
		for cursor.Next(ctx) {
			instance :=new(APP)
			if err = cursor.Decode(instance); err == nil { // 反序列化bson到对象
				instances = append(instances, instance)
			}
		}
		return
	}
	err = db.Find(coll, make(map[string]string), fu, skip, limit, 1)
	return
}

func InsertApp(coll string, body io.Reader) (id string, err error) {
	byt, err := ioutil.ReadAll(body)
	if err != nil {
		return
	}
	instance := new(pb.App)
	err = json.Unmarshal(byt, instance)
	if err != nil {
		return
	}
	//instance.ID = primitive.NewObjectID()
	return db.Insert(coll, instance)
}
