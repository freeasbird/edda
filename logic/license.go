package logic

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"time"

	pb "github.com/offer365/eddacore/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Result struct {
	SerialNum string             `json:"serial_num"`
	Apps      map[string]*pb.App `json:"apps"`
}

func FindOneLicense(coll string, id string) (instances []*pb.License, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", oid}}
	instance := new(pb.License)
	err = db.FindOne(coll, filter, instance)
	return []*pb.License{instance}, err
}

func FindLicense(coll string, filter interface{}, skip, limit int64) (instances []*pb.License, err error) {
	instances = make([]*pb.License, 0)
	callback := func(cursor *mongo.Cursor) (err error) {
		// 遍历结果集
		ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
		for cursor.Next(ctx) {
			instance := new(pb.License)
			if err = cursor.Decode(instance); err == nil { // 反序列化bson到对象
				instances = append(instances, instance)
			}
		}
		return
	}
	err = db.Find(coll, make(map[string]string), callback, skip, limit, -1)
	return
}

func InsertLicense(coll string, body io.Reader) (cipher, id string, err error) {
	byt, err := ioutil.ReadAll(body)
	if err != nil {
		return
	}
	res := new(Result)
	res.Apps = make(map[string]*pb.App, 0)
	err = json.Unmarshal(byt, res)
	if err != nil {
		return
	}
	// if byt,err=base64.StdEncoding.DecodeString(res.SerialNum);err!=nil{
	//	return
	// }
	// if byt,err=endecrypt.Decrypt(endecrypt.Pri1AesRsa2048,byt);err!=nil{
	//	return
	// }
	//
	// sn:=new(pb.SerialNum)
	// if err=json.Unmarshal(byt,sn);err!=nil{
	//	return
	// }
	// dev := make(map[string]string)
	// for k, v := range sn.Nodes {
	//	dev[k] = v.Attrs.Hwmd5
	//	fmt.Println(v.Attrs.Hwmd5)
	// }

	apps := make(map[string]*pb.App, 0)
	for key, app := range res.Apps {
		app.Key = key
		app.MaxLifeCycle = (app.Expire - time.Now().Unix()) / 60
		apps[key] = app
	}
	req := pb.AuthReq{
		Cipher: &pb.Cipher{Code: res.SerialNum},
		Apps:   apps,
	}
	resp, err := pb.Auth.Authorized(context.TODO(), &req)
	if err != nil {
		return
	}
	// lic.ID = primitive.NewObjectID()
	id, err = db.Insert(coll, resp.Lic)
	cipher = resp.Cipher.Code
	return
}

// func Aggregation(coll string, id string, skip, limit int64) (data interface{}, err error) {
//	var pipe mongo.Pipeline
//	//{"$limit",limit},
//	//{"$skip",skip},
//	//{"$sort",-1},
//	if id != "" {
//		if oid, err := primitive.ObjectIDFromHex(id); err == nil {
//			o1 := bson.D{
//				{"$match", bson.M{"_id": oid}},
//			}
//			pipe = append(pipe, o1)
//		}
//	}
//
//	// db.products.aggregate([{$match:{"_id":ObjectId("5d5d0a3a306ba203ca7447a1")}},{$lookup:{from:"projects",localField:"projects",foreignField:"_id",as:"projects"}},{$lookup:{from:"users",localField:"authors",foreignField:"_id",as:"authors"}}]).pretty()
//	o2 := bson.D{{
//		"$lookup", bson.M{
//			"from":         "projects",
//			"localField":   "projects",
//			"foreignField": "_id", "as": "projects",
//		},
//	}}
//	pipe = append(pipe, o2)
//	o3 := bson.D{{
//		"$lookup", bson.M{
//			"from":         "principals",
//			"localField":   "principal",
//			"foreignField": "_id", "as": "principal",
//		},
//	}}
//	pipe = append(pipe, o3)
//	o4 := bson.D{{
//		"$lookup", bson.M{
//			"from":         "fs.files",
//			"localField":   "files",
//			"foreignField": "_id", "as": "files",
//		},
//	}}
//	pipe = append(pipe, o4)
//	o5 := bson.D{{
//		"$lookup", bson.M{
//			"from":         "copyrights",
//			"localField":   "copyright",
//			"foreignField": "_id", "as": "copyright",
//		},
//	}}
//	pipe = append(pipe, o5)
//	pipe = append(pipe, bson.D{{"$limit", limit}})
//	pipe = append(pipe, bson.D{{"$skip", skip}})
//	//pipe=append(pipe,bson.D{{"$sort",-1}})
//
//	instances := make([]*model.Show, 0)
//	fu := func(cursor *mongo.Cursor) (err error) {
//		// 遍历结果集
//		ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
//		for cursor.Next(ctx) {
//			instance := new(model.Show)
//			if err = cursor.Decode(instance); err == nil { // 反序列化bson到对象
//				instances = append(instances, instance)
//			}
//		}
//		return
//	}
//
//	err = db.Aggregation(coll, pipe, fu)
//	return instances, err
// }
