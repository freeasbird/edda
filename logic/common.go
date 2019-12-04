package logic

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"time"

	"github.com/offer365/example/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var db mongodb.DB

func Init(addr, user, pwd, database string, timeout time.Duration, ci map[string]string) (err error) {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	db = mongodb.NewDB("mongo")
	return db.Init(ctx,
		mongodb.WithAddr(addr),
		mongodb.WithUsername(user),
		mongodb.WithPwd(pwd),
		mongodb.WithDB(database),
		mongodb.WithTimeout(timeout),
		mongodb.WithCollIndex(ci),
	)
}

func Update(coll string, id string, body io.Reader) (err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", oid}}
	byt, err := ioutil.ReadAll(body)
	if err != nil {
		return
	}
	data := make(map[string]interface{}, 0)
	err = json.Unmarshal(byt, &data)
	if err != nil {
		return
	}
	update := bson.D{}
	for k, v := range data {
		update = append(update, bson.E{"$set", bson.D{{k, v}}})
	}
	return db.Update(coll, filter, update)
}

func Delete(coll string, id string) (err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", oid}}
	return db.Delete(coll, filter)
}

func Count(coll string) (num int64, err error) {
	return db.Count(coll, make(map[string]string))
}

type Show struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Customer string             `bson:"customer" json:"customer"` // 客户
	APPs     []string           `bson:"apps" json:"apps"`         // 被授权了哪些产品
	Attr     string             `bson:"attr" json:"attr"`         // 授权属性  测试或正式
	Duration time.Duration      `bson:"duration" json:"duration"` // 授权时长
	Expire   time.Time          `bson:"expire" json:"expire"`     // 到期时间
}
