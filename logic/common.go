package logic

import (
	"../dao"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"io/ioutil"
	"time"
)

var device dao.DB

func Init(host, port, user, pwd, db string, timeout time.Duration, ci map[string]string) (err error) {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	device = dao.NewDB("mongo")
	return device.Init(ctx,
		dao.WithHost(host),
		dao.WithPort(port),
		dao.WithUsername(user),
		dao.WithPwd(pwd),
		dao.WithDB(db),
		dao.WithTimeout(timeout),
		dao.WithCollIndex(ci),
	)
}

func Update(coll string, id string, body io.Reader) (err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", oid}}
	byt, err := ioutil.ReadAll(body)
	if err != nil {
		return
	}
	fmt.Println(string(byt))
	data := make(map[string]interface{}, 0)
	err = json.Unmarshal(byt, &data)
	if err != nil {
		return
	}
	update := bson.D{}
	for k, v := range data {
		update = append(update, bson.E{"$set", bson.D{{k, v}}})
	}
	return device.Update(coll, filter, update)
}

func Delete(coll string, id string) (err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", oid}}
	return device.Delete(coll, filter)
}

func Count(coll string) (num int64, err error) {
	return device.Count(coll, make(map[string]string))
}
