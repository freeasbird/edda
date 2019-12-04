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

func FindOneSerial(coll string, id string) (instances []*pb.SerialNum, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", oid}}
	instance := new(pb.SerialNum)
	err = db.FindOne(coll, filter, instance)
	return []*pb.SerialNum{instance}, err
}

func FindAllSerial(coll string, skip, limit int64) (instances []*pb.SerialNum, err error) {
	instances = make([]*pb.SerialNum, 0)
	fu := func(cursor *mongo.Cursor) (err error) {
		// 遍历结果集
		ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
		for cursor.Next(ctx) {
			instance := new(pb.SerialNum)
			if err = cursor.Decode(instance); err == nil { // 反序列化bson到对象
				instances = append(instances, instance)
			}
		}
		return
	}
	err = db.Find(coll, make(map[string]string), fu, skip, limit, 1)
	return
}

func InsertSerial(coll string, body io.Reader) (id string, err error) {
	byt, err := ioutil.ReadAll(body)
	if err != nil {
		return
	}
	instance := new(pb.SerialNum)
	err = json.Unmarshal(byt, instance)
	if err != nil {
		return
	}
	// instance.ID = primitive.NewObjectID()
	return db.Insert(coll, instance)
}

//
// func ResolveSerial( cipher string ) (msg []string, err error) {
//	var (
//		filter   bson.D
//		ins      []*pb.Node
//		byt      []byte
//		licenses []*pb.License
//		sn       *pb.SerialNum
//	)
//
//	if byt,err=base64.StdEncoding.DecodeString(cipher);err!=nil{
//		return
//	}
//	if byt,err=endecrypt.Decrypt(endecrypt.Pri1AesRsa2048,byt);err!=nil{
//		return
//	}
//	sn=new(pb.SerialNum)
//	fmt.Println(string(byt))
//	err = json.Unmarshal(byt, &sn)
//	//byt, err = json.Marshal(sn)
//	for _, n := range sn.Nodes {
//		// 硬件md5
//		filter = bson.D{{"attr.md5", n.Attrs.Hwmd5}}
//		nodes, _ := FindNode("nodes", filter, 0, 0)
//		if len(nodes) == 0 {
//			msg = append(msg, "未在数据库中找到该设备:"+n.Attrs.Name+"。")
//			// 设备id
//			filter = bson.D{{"hardware.host.machineid", n.Hardware.Host.Machineid}}
//			machine, _ := FindNode("nodes", filter, 0, 0)
//			if len(machine) > 0 {
//				for _, mh := range machine {
//					msg = append(msg, "在数据库中找到与该设备一致的设备id:"+mh.Hardware.Host.Machineid+"。")
//				}
//			}
//			nodes = append(nodes, machine...)
//			// 产品序列号
//			filter = bson.D{{"hardware.product.serial", n.Hardware.Product.Serial}}
//			serial, _ := FindNode("nodes", filter, 0, 0)
//			if len(serial) > 0 {
//				for _, sr := range serial {
//					msg = append(msg, "在数据库中找到与该设备一致的产品序列号:"+sr.Hardware.Product.Serial+"。")
//				}
//			}
//			nodes = append(nodes, serial...)
//			for _, nw := range n.Hardware.Networks {
//				//  mac
//				filter = bson.D{{"hardware.networks.macaddress", nw.Macaddress}}
//				mac, _ := FindNode("nodes", filter, 0, 0)
//				if len(mac) > 0 {
//					for _, m := range mac {
//						for _, ma := range m.Hardware.Networks {
//							msg = append(msg, "在数据库中找到与该设备一致的MAC地址:"+ma.Macaddress+"。")
//						}
//					}
//				}
//				nodes = append(nodes, mac...)
//			}
//		}
//		ins = append(ins, nodes...)
//	}
//	devices := make(map[primitive.ObjectID]*pb.Node)
//	for _, in := range ins {
//		devices[in.ID] = in
//	}
//	for _, dev := range devices {
//		filter = bson.D{{"devices.md5", dev.Attrs.Hwmd5}}
//		licenses, err = FindLicense("licenses", filter, 0, 0)
//		licenses = append(licenses, licenses...)
//	}
//
//	licenseM := make(map[primitive.ObjectID]*pb.License)
//	for _, lic := range licenses {
//		licenseM[lic.ID] = lic
//	}
//	msg = append(msg, "在数据库中找到相似的设备:"+fmt.Sprintf("%d个。", len(devices)))
//	return
// }
