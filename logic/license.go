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

type Result struct {
	Customer  string         `json:"customer"`
	Project   string         `json:"project"`
	SerialNum string         `json:"serial_num"`
	Apps      map[string]App `json:"apps"`
}

type App struct {
	Name       string         `json:"name"`
	ExpireTime string         `json:"expire_time"`
	Instance   int64          `json:"instance"`
	Attr       map[string]int `json:"attr"`
}

func FindOneLicense(coll string, id string) (instances []*model.License, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", oid}}
	instance := new(model.License)
	err = device.FindOne(coll, filter, instance)
	return []*model.License{instance}, err
}

func FindLicense(coll string, filter interface{}, skip, limit int64) (instances []*model.License, err error) {
	instances = make([]*model.License, 0)
	fu := func(cursor *mongo.Cursor) (err error) {
		// 遍历结果集
		ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
		for cursor.Next(ctx) {
			instance := new(model.License)
			if err = cursor.Decode(instance); err == nil { // 反序列化bson到对象
				instances = append(instances, instance)
			}
		}
		return
	}
	err = device.Find(coll, make(map[string]string), fu, skip, limit, -1)
	return
}

func InsertLicense(coll string, body io.Reader) (id string, err error) {
	byt, err := ioutil.ReadAll(body)
	if err != nil {
		return
	}
	res := new(Result)
	res.Apps = make(map[string]App, 0)
	err = json.Unmarshal(byt, res)
	if err != nil {
		return
	}
	//instance.ID = primitive.NewObjectID()
	sn, _ := model.Decrypt(res.SerialNum)
	dev := make(map[string]string)
	for k, v := range sn.Nodes {
		dev[k] = v.HwMd5
	}
	apps := make([]*model.APP, 0)
	for key, app := range res.Apps {
		ap := model.NewAPP(key, app.Name, app.ExpireTime, app.Instance, app.Attr)
		apps = append(apps, ap)
	}
	lic, err := model.NewLicense(sn.Sid, dev, apps...)
	lic.ID = primitive.NewObjectID()
	lic.CipherText()
	lic.SerialNum = res.SerialNum
	lic.Customer = res.Customer
	lic.Project = res.Project
	byt, err = json.Marshal(lic)
	id, err = device.Insert(coll, lic)
	for _, node := range sn.Nodes {
		if err == nil {
			_, err = device.Insert("nodes", node)
		}
	}
	return
}

func Aggregation(coll string, id string, skip, limit int64) (data interface{}, err error) {
	var pipe mongo.Pipeline
	//{"$limit",limit},
	//{"$skip",skip},
	//{"$sort",-1},
	if id != "" {
		if oid, err := primitive.ObjectIDFromHex(id); err == nil {
			o1 := bson.D{
				{"$match", bson.M{"_id": oid}},
			}
			pipe = append(pipe, o1)
		}
	}

	// db.products.aggregate([{$match:{"_id":ObjectId("5d5d0a3a306ba203ca7447a1")}},{$lookup:{from:"projects",localField:"projects",foreignField:"_id",as:"projects"}},{$lookup:{from:"users",localField:"authors",foreignField:"_id",as:"authors"}}]).pretty()
	o2 := bson.D{{
		"$lookup", bson.M{
			"from":         "projects",
			"localField":   "projects",
			"foreignField": "_id", "as": "projects",
		},
	}}
	pipe = append(pipe, o2)
	o3 := bson.D{{
		"$lookup", bson.M{
			"from":         "principals",
			"localField":   "principal",
			"foreignField": "_id", "as": "principal",
		},
	}}
	pipe = append(pipe, o3)
	o4 := bson.D{{
		"$lookup", bson.M{
			"from":         "fs.files",
			"localField":   "files",
			"foreignField": "_id", "as": "files",
		},
	}}
	pipe = append(pipe, o4)
	o5 := bson.D{{
		"$lookup", bson.M{
			"from":         "copyrights",
			"localField":   "copyright",
			"foreignField": "_id", "as": "copyright",
		},
	}}
	pipe = append(pipe, o5)
	pipe = append(pipe, bson.D{{"$limit", limit}})
	pipe = append(pipe, bson.D{{"$skip", skip}})
	//pipe=append(pipe,bson.D{{"$sort",-1}})

	instances := make([]*model.Show, 0)
	fu := func(cursor *mongo.Cursor) (err error) {
		// 遍历结果集
		ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
		for cursor.Next(ctx) {
			instance := new(model.Show)
			if err = cursor.Decode(instance); err == nil { // 反序列化bson到对象
				instances = append(instances, instance)
			}
		}
		return
	}

	err = device.Aggregation(coll, pipe, fu)
	return instances, err
}

//func UpdateProduct(coll string, id string, body io.Reader) (err error) {
//	oid, err := primitive.ObjectIDFromHex(id)
//	filter := bson.D{{"_id", oid}}
//	byt, err := ioutil.ReadAll(body)
//	if err != nil {
//		return
//	}
//
//	pd := new(model.Show)
//	err = json.Unmarshal(byt, pd)
//	if err != nil {
//		return
//	}
//	data := make(map[string]interface{}, 0)
//	principal := make([]primitive.ObjectID, 0)
//	copyrights := make([]primitive.ObjectID, 0)
//	projects := make([]primitive.ObjectID, 0)
//	files := make([]primitive.ObjectID, 0)
//
//
//	if oid, err := primitive.ObjectIDFromHex(pd.Principal); err == nil {
//		principal = append(principal, oid)
//	}
//
//	if oid, err := primitive.ObjectIDFromHex(pd.Copyright); err == nil {
//		copyrights = append(copyrights, oid)
//	}
//
//	for _, id := range pd.Projects {
//		if oid, err := primitive.ObjectIDFromHex(id); err == nil {
//			projects = append(projects, oid)
//		}
//	}
//	for _, id := range pd.Files {
//		if oid, err := primitive.ObjectIDFromHex(id); err == nil {
//			files = append(files, oid)
//		}
//	}
//
//	data["name"] = pd.Name
//	data["version"] = pd.Version
//	data["number"] = pd.Number
//	data["introduction"] = pd.Introduction
//	data["files"] = files
//	data["projects"] = projects
//	data["principal"] = principal
//	data["copyright"] = copyrights
//	update := bson.D{}
//	for k, v := range data {
//		update = append(update, bson.E{"$set", bson.D{{k, v}}})
//	}
//	return device.Update(coll, filter, update)
//}
